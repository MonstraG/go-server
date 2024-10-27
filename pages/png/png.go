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

// Validate accepts a PNG, checks it according to the PNG specification
// http://www.libpng.org/pub/png/spec/1.2/png-1.2.pdf
// and returns parsed chunks as JSON body
func Validate(w helpers.MyWriter, r *helpers.MyRequest) {
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

	err = validateChunks(chunks)
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
var ancillaryChunkTypes = []string{
	"bKGD", "cHRM", "cICP", "dSIG", "eXIf",
	"gAMA", "hIST", "iCCP", "iTXt", "pHYs",
	"sBIT", "sPLT", "sRGB", "sTER", "tEXt",
	"tIME", "tRNS", "zTXt"}

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

func validateChunks(chunks []*Chunk) error {
	var firstChunk = chunks[0]
	if firstChunk.ChunkType != "IHDR" {
		return fmt.Errorf(
			"IHDR chunk must appear first, PLTEFound %s instead",
			firstChunk.ChunkType)
	}

	plteRequirements, err := validateIHDRChunk(firstChunk)
	if err != nil {
		return err
	}

	for chunkIndex, chunk := range chunks {
		if chunkIndex == 0 {
			// not just chunks[1:] to keep chunk index
			continue
		}

		err = validateChunk(chunk, plteRequirements)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateChunk(chunk *Chunk, plteRequirements *ValidatorState) error {
	switch chunk.ChunkType {
	case "PLTE":
		err := validatePLTEChunk(chunk, plteRequirements)
		if err != nil {
			return err
		}
	case "IDAT":
		plteRequirements.IDATFound = true
	}
	return nil
}

type ValidatorState struct {
	PLTERequired    bool
	PLTEDissallowed bool
	PLTEFound       bool
	IDATFound       bool
}

func validateIHDRChunk(chunk *Chunk) (*ValidatorState, error) {
	widthBuffer := chunk.Data[0:4]
	width := binary.BigEndian.Uint32(widthBuffer)
	if width == 0 {
		return nil, fmt.Errorf("IHDR's width cannot be zero")
	}

	heightBuffer := chunk.Data[4:8]
	height := binary.BigEndian.Uint32(heightBuffer)
	if height == 0 {
		return nil, fmt.Errorf("IHDR's height cannot be zero")
	}

	bitDepth := int(chunk.Data[8])
	colorType := int(chunk.Data[9])

	bitDepthsByColorType := map[int][]int{
		0: {1, 2, 4, 8, 16},
		2: {8, 16},
		3: {1, 2, 4, 8},
		4: {8, 16},
		6: {8, 16},
	}

	possibleBitDepths, found := bitDepthsByColorType[colorType]
	if !found {
		return nil, fmt.Errorf("IHDR's color type is invalid, PLTEFound %v", colorType)
	}

	if !slices.Contains(possibleBitDepths, bitDepth) {
		return nil, fmt.Errorf("IHDR's bit depth is invalid, PLTEFound %v", bitDepth)
	}

	compression := int(chunk.Data[10])
	if compression != 0 {
		return nil, fmt.Errorf("IHDR's compression is invalid, PLTEFound %v, only '0' is valid", compression)
	}

	filterMethod := int(chunk.Data[11])
	if filterMethod != 0 {
		return nil, fmt.Errorf("IHDR's filter method is invalid, PLTEFound %v, only '0' is valid", filterMethod)
	}

	interlace := int(chunk.Data[12])
	if interlace != 0 && interlace != 1 {
		return nil, fmt.Errorf("IHDR's interlace is invalid, PLTEFound %v, only '0' or '1' are valid", interlace)
	}

	if len(chunk.Data) > 12 {
		return nil, fmt.Errorf("IHDR's data is invalid, PLTEFound more data after 13 PLTERequired bytes")
	}

	PLTERequired := false
	if colorType == 3 {
		PLTERequired = true
	}

	PLTEDisallowed := false
	if colorType == 0 || colorType == 4 {
		PLTEDisallowed = true
	}

	return &ValidatorState{
		PLTERequired:    PLTERequired,
		PLTEDissallowed: PLTEDisallowed,
	}, nil
}

func validatePLTEChunk(chunk *Chunk, requirements *ValidatorState) error {
	if requirements.PLTEFound {
		return fmt.Errorf("found second PLTE chunk")
	}
	if requirements.IDATFound {
		return fmt.Errorf("found IDAT chunk earlier than PLTE chunk")
	}
	if requirements.PLTEDissallowed {
		return fmt.Errorf("PLTE must not appear due to specified color mode")
	}
	if chunk.Length%3 != 0 {
		return fmt.Errorf("PLTE chunk's length must be a multiple of 3 (pallete colors have 3 channels, RGB)")
	}

	requirements.PLTEFound = true
	return nil
}
