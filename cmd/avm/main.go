package main

import (
	"fmt"
	"math/rand"
)

func getUint64(p []byte, size int) (v uint64) {
	for i := 1; i <= size && i < len(p); i++ {
		v <<= 8
		v |= uint64(p[i])
	}
	return v
}

func putUint64(p []byte, offset, v uint64) {
	for i := offset; i <= wordSize && i < uint64(len(p)); i++ {
		p[i] = byte(v)
		v >>= 8
	}
}

func randBytes(size int) []byte {
	p := make([]byte, size)
	rand.Read(p)
	return p
}

const wordSize = 8

var world = make([]byte, wordSize)

func main() {
	threadCount := uint64(1)
	var loop func(uint64, uint64)
	loop = func(thread, start uint64) {
		for i := start; i < uint64(len(world)); i++ {
			switch world[i] {
			case 0: // random
				rand8 := randBytes(wordSize)
				world = append(world, make([]byte, wordSize-1)...)
				copy(world[i:], rand8)
			case 1: // jump
				var dst uint64
				if i > wordSize {
					dst = getUint64(world[i-wordSize:], wordSize)
					if dst < uint64(len(world)) && i != dst+1 {
						fmt.Println(thread, i, "jump", dst)
						i = dst
						break
					}
				}
				world[i] = 0
			case 2: // word step
				fmt.Println(thread, i, "data")
				i += wordSize
			case 3: // add
				fmt.Println(thread, i, "add")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v += getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += 2 * wordSize
			case 4: // sub
				fmt.Println(thread, i, "sub")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v -= getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += 2 * wordSize
			case 5: // mul
				fmt.Println(thread, i, "mul")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v *= getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += 2 * wordSize
			case 6: // div
				fmt.Println(thread, i, "div")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				p := getUint64(world[i:], wordSize)
				if p == 0 {
					v = 0
				} else {
					v /= p
				}
				putUint64(world[i:], wordSize, v)
				i += 2 * wordSize
			case 7: // copy
				fmt.Println(thread, i, "copy")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				putUint64(world[i:], 0, v)
				i += wordSize
			case 8: // go
				if i > wordSize {
					dst := getUint64(world[i-wordSize:], wordSize)
					if dst < uint64(len(world)) && i != dst+1 {
						fmt.Println(thread, i, "go", dst)
						subthread := threadCount
						threadCount++
						go loop(subthread, dst)
					}
				}
			default:
				world[i] = 0
			}
		}
	}
	for {
		loop(0, 0)
	}
}
