package helpers

import (
	"log"
	"os"
)

func ReadData(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read database file", err)
	}

	return data
}

func WriteData(path string, bytes []byte) {
	err := os.WriteFile(path, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
