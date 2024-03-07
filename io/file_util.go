package io

import (
	"bufio"
	"log"
	"os"
)

func flushWriter(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Fatalf("Could not flush writer: %s", err)
	}
}

func truncateFile(fileName string) {
	file := OpenFile(fileName, os.O_TRUNC)
	defer CloseFile(file)
}

func OpenFile(fileName string, flags int) *os.File {
	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file = CreateFile(fileName)
		} else {
			log.Fatalf("failed to open file '%s': %s", fileName, err)
		}
	}
	return file
}

func CreateFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create file '%s': %s", fileName, err)
	}
	return file
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatalf("failed to close file: %s", err)
	}
}
