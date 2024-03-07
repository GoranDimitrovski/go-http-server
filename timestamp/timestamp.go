package timestamp

import (
	"fmt"
	"net/http"
	"simplesurance/config"
	"simplesurance/io"
	"simplesurance/store"
	"time"
)

var memoryStore *store.TimestampMemoryStore

func Init(fileName string) {
	timestamp := time.Now().Unix()

	memoryStore = store.New(fileName)
	memoryStore.Load()
	memoryStore.RemoveExpired(int(timestamp), config.ServerConfig.Threshold)

	io.Rewrite(memoryStore.View(), fileName)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Unix()

	memoryStore.RemoveExpired(int(timestamp), config.ServerConfig.Threshold)
	memoryStore.Store(int(timestamp))
	memoryStore.Sync()
	fmt.Fprint(w, memoryStore.Count())
}
