package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	// The Keccak permutation
	keccakf func([]uint64)
	// The output of the permutation
	state []uint64
	// The rate of the sponge
	rate int
	// The output length
	outputLength = 32
	// number of rounds
	nrRounds = 24
)
var (
	bc, i, j int
)

// The round constants
var rc = [24]uint64{
	0x0000000000000001, 0x0000000000008082, 0x800000000000808a,
	0x8000000080008000, 0x000000000000808b, 0x0000000080000001,
	0x8000000080008081, 0x8000000000008009, 0x000000000000008a,
	0x0000000000000088, 0x0000000080008009, 0x000000008000000a,
	0x000000008000808b, 0x800000000000008b, 0x8000000000008089,
	0x8000000000008003, 0x8000000000008002, 0x8000000000000080,
	0x000000000000800a, 0x800000008000000a, 0x8000000080008081,
	0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
}

var rotc = [24]uint{1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 2, 14,
	27, 41, 56, 8, 25, 43, 62, 18, 39, 61, 20, 44}

var piln = [24]uint{10, 7, 11, 17, 18, 3, 5, 16, 8, 21, 24, 4,
	15, 23, 19, 13, 12, 2, 20, 14, 22, 9, 6, 1}

// Keccak hashing function
// Accepts a message in bytes and hashes in place.
func Keccak(message []byte, inRate int) {
	// Initialize vars
	rate = inRate
	// Initialize the state
	state = make([]uint64, 25)

	for i := 0; i < 25; i++ {
		state[i] = 0
	}

	// The permutation
	keccakf = func(state []uint64) {
		// The round function
		for bc := 0; bc < nrRounds; bc++ {
			// Theta step
			C := make([]uint64, 5)
			for i := 0; i < 5; i++ {
				C[i] = state[i] ^ state[i+5] ^ state[i+10] ^ state[i+15] ^ state[i+20]
			}
			for i := 0; i < 5; i++ {
				D := C[(i+4)%5] ^ rotl_64(C[(i+1)%5], 1)
				for j := 0; j < 25; j += 5 {
					state[j+i] ^= D
				}
			}

			// Rho and pi steps
			t := state[1]
			for i := 0; i < 24; i++ {
				j := piln[i]
				C[0] = state[j]
				state[j] = rotl_64(t, rotc[i])
				t = C[0]
			}

			// Chi step
			for j := 0; j < 25; j++ {
				C := make([]uint64, 5)
				copy(C, state[5+j-5:])

				for i := 0; i < 5; i += 5 {
					state[j+i] = state[j+i] ^ ((^C[(i+1)%5]) & C[(i+2)%5])
				}
			}
			// Iota step
			state[0] ^= rc[bc]
		}
	}

	// Absorb the message
	for i := 0; i < len(message); i += rate / 8 {
		P := message
		if len(message)%rate/8 != 0 {
			P = padTenOne(message)
		}
		block := P[i : i+(rate/8)]
		for j := 0; j < len(block); j++ {
			state[j/8] ^= uint64(block[j]) << (uint(j) * 8)
		}
		keccakf(state)
	}

	// // Squeeze the output
	// output := make([]byte, outputLength)
	// for i := 0; i < len(output); i += 8 {
	// 	block := state[i/8 : i/8+8]
	// 	for j := 0; j < len(block); j++ {
	// 		output[i+j] = byte(block[j] >> (uint(j) * 8))
	// 	}
	// 	if i+8 >= outputLength {
	// 		break
	// 	}
	// 	keccakf(state)
	// }
	// return output
}

// Rotate left
func rotl_64(x uint64, n uint) uint64 {
	return (x << n) | (x >> (64 - n))
}

func padTenOne(X []byte) []byte {

	q := rate - len(X)%rate
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	b := byte(0x80)
	padded[len(X)+q-1] = b
	return padded
}

func lrEncode(X uint64, dir bool) []byte {

	emptyX := make([]byte, 2)

	if X == 0 && dir {
		emptyX[0] = 0
		emptyX[1] = 1
		return emptyX
	} else if X == 0 && !dir {
		emptyX[0] = 1
		emptyX[1] = 0
		return emptyX
	}

	temp := make([]byte, 255)
	length := X
	count := 0

	for length > 0 {
		b := byte(length & 0xff)
		length = length >> 8
		temp[245-count] = b
		count++
	}

	result := make([]byte, count+1)
	copy(result, temp[255-count:])
	if dir {
		result[len(result)-1] = byte(count)
	} else {
		result[0] = byte(count)
	}
	return result
}

func encode_string(S []byte) []byte { return append(lrEncode(uint64(len(S)*8), false), S...) }

func bytepad(X []byte, w uint64) []byte {

	enc_w := lrEncode(w, false)
	// w * ((enc_w.length + X.length + w - 1) / w) = smallest multiple of w and z.length
	z := make([]byte, w*((uint64(len(enc_w)+len(X))+w-1)/w))
	copy(z, enc_w)
	copy(z[len(enc_w):], X)
	return z

}

func readFileChunk(file *os.File, chunkSize int64, pos int64, wg *sync.WaitGroup, chunks *chan []byte) {
	b := make([]byte, chunkSize)
	file.ReadAt(b, pos*chunkSize)
	Keccak(b, 1088)
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

func main() {

	run()

}

func generateRandomBytes() []byte {
	b := make([]byte, 1000)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return b
}
