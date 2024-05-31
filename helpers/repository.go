package helpers

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

func ReadData[T any](dbFilePath string) *[]T {
	_, err := os.Stat(dbFilePath)
	if os.IsNotExist(err) {
		return &[]T{}
	}

	bytes, err := os.ReadFile(dbFilePath)
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

func WriteData(dbFilePath string, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Failed to marshall data:\n", err)
	}
	err = os.MkdirAll(path.Dir(dbFilePath), os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create database dbFilePath:\n", err)
	}
	err = os.WriteFile(dbFilePath, bytes, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to write database file:\n", err)
	}
}
