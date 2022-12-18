package main

import (
	"fmt"
	"os"
	"sync"
)

const producerCount int = 64
const consumerCount int = 64

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
	fmt.Println(pos)
	wg.Done()
}

func consume(link <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	ch := <-link
	SpongeSqueeze(SpongeAbsorb(&ch, 512), 1344/8, 512)

}

func main() {

	file, _ := os.Open("bible.txt")
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	chunkSize := 8192
	chunks := fileSize / int64(chunkSize)

	link := make(chan []byte)
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}
	fmt.Println(chunks)
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
}
