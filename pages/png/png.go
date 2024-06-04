package png

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-server/helpers"
	"hash/crc32"
	"io"
	"net/http"
	"slices"
)

func Validate(w helpers.MyWriter, r *http.Request) {
	ok, err := isValidPNGHeader(r.Body)
	if err != nil {
		w.WriteBadRequest(err.Error())
		return
	}
	if !ok {
		w.WriteBadRequest("png is invalid")
		return
	}

	chunk, err := readChunk(r.Body)
	if err != nil {
		w.WriteBadRequest(err.Error())
		return
	}

	w.WriteOk(fmt.Sprintf("png is valid, first chunk is of type %v", chunk.ChunkType))
}

const validPNGHeader = "\x89PNG\r\n\x1a\n"

func readBytesToBuffer(body io.ReadCloser, wantBytes int) ([]byte, error) {
	buffer := make([]byte, wantBytes)
	readBytes, err := body.Read(buffer)
	if err != nil {
		return nil, err
	}
	if readBytes != wantBytes {
		return nil, fmt.Errorf("expected to read %d bytes for header, but read %d bytes", wantBytes, readBytes)
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

	crcBuffer, err := readBytesToBuffer(body, 4)
	if err != nil {
		return nil, err
	}

	crc := binary.BigEndian.Uint32(crcBuffer)
	computedCRC := crc32.ChecksumIEEE(append(chunkTypeBuffer, contentBuffer...))

	if computedCRC != crc {
		return nil, fmt.Errorf("CRC mismatch for chunk type: %s", chunkType)
	}

	return &Chunk{ChunkType: chunkType, Data: crcBuffer, Length: chunkLength}, nil
}
