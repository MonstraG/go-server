package helpers

import (
	"encoding/json"
	"log"
	"os"
)

func ReadData[T any](path string) *[]T {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read database file:\n", err)
	}

	var data []T
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal("Failed to read from database file:\n", err)
	}
	return &data
}

func WriteData(path string, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Failed to marshall data:\n", err)
	}
	err = os.WriteFile(path, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file:\n", err)
	}
}
