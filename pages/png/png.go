package png

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-server/helpers"
	"hash/crc32"
	"io"
	"net/http"
	"slices"
)

const maxAllowedPngSize = 10 * 1024 * 1024

// Validate accepts a PNG, checks it according to the PNG specification www.libpng.org/pub/png/spec/1.2/png-1.2.pdf
// and returns parsed chunks as JSON body
func Validate(w helpers.MyWriter, r *http.Request) {
	ok, err := isValidPNGHeader(r.Body)
	if err != nil {
		w.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
		return
	}
	if !ok {
		w.WriteResponse(http.StatusBadRequest, []byte("png is invalid"))
		return
	}

	chunks, err := readChunks(r.Body)
	if err != nil {
		w.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
		return
	}

	err = validateIHDRChunk(chunks)
	if err != nil {
		w.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
	}

	chunksJson, err := json.Marshal(chunks)
	if err != nil {
		w.WriteResponse(http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteResponse(http.StatusOK, chunksJson)
}

const validPNGHeader = "\x89PNG\r\n\x1a\n"

func readBytesToBuffer(body io.ReadCloser, wantBytes int) ([]byte, error) {
	buffer := make([]byte, wantBytes)
	readBytes, err := io.ReadAtLeast(body, buffer, wantBytes)
	if err != nil {
		return nil, err
	}
	if readBytes != wantBytes {
		return nil, fmt.Errorf("expected to read %d bytes, but read %d bytes", wantBytes, readBytes)
	}
	return buffer, nil
}

func isValidPNGHeader(body io.ReadCloser) (bool, error) {
	givenHeader, err := readBytesToBuffer(body, len(validPNGHeader))
	if err != nil {
		return false, err
	}
	return bytes.HasPrefix(givenHeader, []byte(validPNGHeader)), nil
}

type Chunk struct {
	ChunkType string
	Data      []byte
	Length    int
}

var criticalChunkTypes = []string{"IHDR", "PLTE", "IDAT", "IEND"}
var ancillaryChunkTypes = []string{"bKGD", "cHRM", "cICP", "dSIG", "eXIf", "gAMA", "hIST", "iCCP", "iTXt", "pHYs", "sBIT", "sPLT", "sRGB", "sTER", "tEXt", "tIME", "tRNS", "zTXt"}

func (chunk *Chunk) IsCritical() bool {
	return slices.Contains(criticalChunkTypes, chunk.ChunkType)
}

func (chunk *Chunk) IsAncillary() bool {
	return !chunk.IsCritical()
}

func isValidChunkType(chunkType string) bool {
	return slices.Contains(criticalChunkTypes, chunkType) ||
		slices.Contains(ancillaryChunkTypes, chunkType)
}

func readChunks(body io.ReadCloser) ([]*Chunk, error) {
	chunks := make([]*Chunk, 0)
	bytesRead := 0

	for {
		chunk, err := readChunk(body)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
		if chunk.ChunkType == "IEND" {
			return chunks, nil
		}

		bytesRead += chunk.Length
		if bytesRead >= maxAllowedPngSize {
			return nil, fmt.Errorf("PNG too big")
		}
	}
}

func readChunk(body io.ReadCloser) (*Chunk, error) {
	chunkLengthBuffer, err := readBytesToBuffer(body, 4)
	if err != nil {
		return nil, err
	}
	chunkLengthUint := binary.BigEndian.Uint32(chunkLengthBuffer)
	chunkLength := int(chunkLengthUint)

	chunkTypeBuffer, err := readBytesToBuffer(body, 4)
	if err != nil {
		return nil, err
	}
	chunkType := string(chunkTypeBuffer)
	if !isValidChunkType(chunkType) {
		return nil, fmt.Errorf("invalid chunk type: %s", chunkType)
	}

	contentBuffer, err := readBytesToBuffer(body, chunkLength)
	if err != nil {
		return nil, err
	}

	crcBuffer, err := readBytesToBuffer(body, 4)
	if err != nil {
		return nil, err
	}

	crc := binary.BigEndian.Uint32(crcBuffer)
	computedCRC := crc32.ChecksumIEEE(append(chunkTypeBuffer, contentBuffer...))

	if computedCRC != crc {
		return nil, fmt.Errorf("CRC mismatch for chunk type: %s", chunkType)
	}

	return &Chunk{ChunkType: chunkType, Data: contentBuffer, Length: chunkLength}, nil
}

func validateIHDRChunk(chunks []*Chunk) error {
	firstChunk := chunks[0]
	if firstChunk.ChunkType != "IHDR" {
		return fmt.Errorf("IHDR chunk must appear first, found %s instead", firstChunk.ChunkType)
	}

	widthBuffer := firstChunk.Data[0:4]
	width := binary.BigEndian.Uint32(widthBuffer)
	if width == 0 {
		return fmt.Errorf("IHDR's width cannot be zero")
	}

	heightBuffer := firstChunk.Data[4:8]
	height := binary.BigEndian.Uint32(heightBuffer)
	if height == 0 {
		return fmt.Errorf("IHDR's height cannot be zero")
	}

	bitDepth := int(firstChunk.Data[8])
	colorType := int(firstChunk.Data[9])

	bitDepthsByColorType := map[int][]int{
		0: {1, 2, 4, 8, 16},
		2: {8, 16},
		3: {1, 2, 4, 8},
		4: {8, 16},
		6: {8, 16},
	}

	possibleBitDepths, found := bitDepthsByColorType[colorType]
	if !found {
		return fmt.Errorf("IHDR's color type is invalid, found %v", colorType)
	}

	if !slices.Contains(possibleBitDepths, bitDepth) {
		return fmt.Errorf("IHDR's bit depth is invalid, found %v", bitDepth)
	}

	compression := int(firstChunk.Data[10])
	if compression != 0 {
		return fmt.Errorf("IHDR's compression is invalid, found %v, only '0' is valid", compression)
	}

	filterMethod := int(firstChunk.Data[11])
	if filterMethod != 0 {
		return fmt.Errorf("IHDR's filter method is invalid, found %v, only '0' is valid", filterMethod)
	}

	interlace := int(firstChunk.Data[12])
	if interlace != 0 && interlace != 1 {
		return fmt.Errorf("IHDR's interlace is invalid, found %v, only '0' or '1' are valid", interlace)
	}

	if len(firstChunk.Data) > 12 {
		return fmt.Errorf("IHDR's data is invalid, found more data after 13 required bytes")
	}

	return nil
}
