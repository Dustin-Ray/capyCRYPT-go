package main

import (
	"io"
	"os"
	"sync"
)

func ReadFileInChunks(fileName string, chunkSize int64, threadCount int, output chan []byte) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(threadCount)
	for threadID := 0; threadID < threadCount; threadID++ {
		go func(threadID int) {
			defer wg.Done()
			for {
				chunk := make([]byte, chunkSize)
				n, err := file.Read(chunk)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					break
				}
				output <- chunk
			}
		}(threadID)
	}
	wg.Wait()
	close(output)
}

func main() {
	chunks := make(chan []byte)

	go ReadFileInChunks("video.vid", 136, 80, chunks)

	// open the output file
	output, err := os.Create("output.vid")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	for chunk := range chunks {
		output.Write(chunk)
	}

}
