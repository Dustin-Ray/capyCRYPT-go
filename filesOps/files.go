package CryptoTool

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func readFileChunk(file *os.File, chunkSize int64, pos int64, wg *sync.WaitGroup, chunks *chan []byte) {
	b := make([]byte, chunkSize)
	file.ReadAt(b, pos*chunkSize)
	KeccakP1600(b, 1344)
	wg.Done()
	*chunks <- b
}

// func KeccakChunk(chunks *chan []byte)

func run() {

	start := time.Now() // start time
	fileName := "/home/dr/Downloads/bible.txt"
	// fileName := "/home/dr/Downloads/input.file"

	var wg sync.WaitGroup

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileSize := fileInfo.Size()
	threads := int64(fileSize / 64)
	chunkSize := int64(fileSize / threads)

	chunks := make(chan []byte)

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := 0; int64(i) < threads; i++ {
		wg.Add(1)
		go readFileChunk(f, chunkSize, int64(i), &wg, &chunks)
	}

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Done reading file.")
	fmt.Println("Time elapsed: ", end)
}
