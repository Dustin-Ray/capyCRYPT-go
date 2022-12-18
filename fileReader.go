package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// reads pos * i number of bytes from a file and pushes the resulting array into link channel.
func produce(filePath string, pos int64, link chan []byte, wg *sync.WaitGroup) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Create a buffer to read the file into
	buf := make([]byte, 8192)
	// Read the file into the buffer
	_, err = f.ReadAt(buf, pos)
	if err != nil {
		panic(err)
	}
	// Push the buffer into the channel
	link <- buf
	wg.Done()
}

// consumes any available byte array in link channel.
func consume(link <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for ch := range link {
		SpongeSqueeze(SpongeAbsorb(&ch, 256), 1344/8, 256)
	}
}

func run() {

	file, _ := os.Open("/home/dr/Downloads/movie/movie.mp4")
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	chunkSize := 8192
	chunks := fileSize/int64(chunkSize) - 2
	fmt.Println("number of goroutines: ", chunks)
	link := make(chan []byte)
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}
	start := time.Now()
	for i := 1; int64(i) < chunks-1; i++ {
		wp.Add(1)
		go produce("/home/dr/Downloads/movie/movie.mp4", int64(i*8192), link, wp)
	}

	for i := 0; int64(i) < chunks; i++ {
		wc.Add(1)
		go consume(link, wc)
	}
	wp.Wait()
	close(link)
	wc.Wait()
	end := time.Now()
	fmt.Println("Elapsed time: ", end.Sub(start))
}
