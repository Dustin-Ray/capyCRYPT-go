package main

// func main() {

// 	threadCount := 8

// 	chunks := make(chan []byte)

// 	start := time.Now() // start time
// 	go ReadFileInChunks("/home/dr/Downloads/input.file", 1024*1024, threadCount, chunks)

// 	outFilename := "/home/dr/Downloads/output.file"
// 	// open the output file
// 	output, err := os.Create(outFilename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer output.Close()

// 	WriteFileInChunks(outFilename, 136, threadCount, chunks)
// 	end := time.Since(start) // end time
// 	fmt.Println("Time elapsed: ", end)

// }
