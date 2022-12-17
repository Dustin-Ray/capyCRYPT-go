package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func readFileChunk(file *os.File, chunkSize int64, pos int64, wg *sync.WaitGroup, chunks *chan []byte) {
	b := make([]byte, chunkSize)
	file.ReadAt(b, pos*chunkSize)
	// Keccak(b, 1344)
	wg.Done()
	*chunks <- b
}

// func KeccakChunk(chunks *chan []byte)

func run() {

	// start := time.Now() // start time
	fileName := "/home/dr/Downloads/bible.txt"
	// fileName := "/home/dr/Downloads/input.file"

	var wg sync.WaitGroup

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileSize := fileInfo.Size()
	// fmt.Println(fileSize / 64)
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
	// end := time.Since(start)
	// fmt.Println("Done reading file.")
	// fmt.Println("Time elapsed: ", end)
}

func fileToByteArray(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	var bytes []byte
	for scanner.Scan() {
		bytes = append(bytes, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	return bytes
}

func main() {

	b, err := ioutil.ReadFile("bible.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(BytesToHexString(SHA3(&b, 512)))
}
