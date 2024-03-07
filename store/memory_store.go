package store

import (
	"bufio"
	"log"
	"os"
	"simplesurance/io"
	"strconv"
	"sync"
)

type TimestampMemoryStore struct {
	timestamps []int
	fileName   string
	input      chan int
	mutex      sync.Mutex
	storeDone  chan struct{}
}

func New(fileName string) *TimestampMemoryStore {
	store := &TimestampMemoryStore{
		timestamps: []int{},
		fileName:   fileName,
		input:      make(chan int),
		storeDone:  make(chan struct{}, 1),
	}
	go store.start()
	return store
}

func (store *TimestampMemoryStore) start() {
	for timestamp := range store.input {
		store.mutex.Lock()
		store.timestamps = append(store.timestamps, timestamp)
		store.mutex.Unlock()
		store.storeDone <- struct{}{} 
	}
}

func (store *TimestampMemoryStore) Store(timestamp int) {
	store.input <- timestamp
	<-store.storeDone 
}

func (store *TimestampMemoryStore) View() []int {
	return append([]int{}, store.timestamps...)
}

func (store *TimestampMemoryStore) Count() int {
	return len(store.timestamps)
}

func (store *TimestampMemoryStore) Load() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	file := io.OpenFile(store.fileName, os.O_RDONLY)
	defer io.CloseFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		timestamp, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("failed parsing timestamp '%s': %s", scanner.Text(), err)
			continue
		}
		store.timestamps = append(store.timestamps, timestamp)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error while scanning file: %s", err)
	}
}

func (store *TimestampMemoryStore) RemoveExpired(current, threshold int) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var timestamps []int
	for _, timestamp := range store.timestamps {
		if current-timestamp < threshold {
			timestamps = append(timestamps, timestamp)
		}
	}
	store.timestamps = timestamps
}

func (store *TimestampMemoryStore) Sync() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	io.Rewrite(store.timestamps, store.fileName)
}
