package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const producerCount int = 4
const consumerCount int = 4

func produce(file *os.File, pos int64, link chan []byte, wg *sync.WaitGroup) {
	// Open the file
	f, err := os.Open("bible.txt")
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

func consume(link <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for ch := range link {
		ch = <-link
		SpongeSqueeze(SpongeAbsorb(&ch, 256), 1344/8, 256)
	}
}

func main() {

	file, _ := os.Open("bible.txt")
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	chunkSize := 8192
	chunks := fileSize/int64(chunkSize) - 2
	fmt.Println("number of goroutines: ", chunks)
	link := make(chan []byte)
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}
	start := time.Now()
	for i := 0; int64(i) < chunks; i++ {
		wp.Add(1)
		go produce(file, int64(i*8192), link, wp)
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
