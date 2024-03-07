package io

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Append(timestamp int, fileName string) {
	WriteToFile([]int{timestamp}, fileName)
}

func Rewrite(timestamps []int, fileName string) {
	truncateFile(fileName)
	WriteToFile(timestamps, fileName)
}

func WriteToFile(timestamps []int, fileName string) {
	file := OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	defer CloseFile(file)

	writer := bufio.NewWriter(file)
	defer flushWriter(writer)

	for _, timestamp := range timestamps {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", timestamp))
		if err != nil {
			log.Fatalf("Could not save timestamps to file \n %s", err)
		}
	}

	err := writer.Flush()
	if err != nil {
		log.Fatalf("Could not flush writer \n %s", err)
	}
}
