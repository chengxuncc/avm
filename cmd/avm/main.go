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

func main() {
	world := make([]byte, wordSize)
	for {
		if len(world) == 0 {
			world = world[:wordSize]
		}
		for i := uint64(0); i < uint64(len(world)); i++ {
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
						fmt.Println(i, "jump", dst)
						i = dst
						break
					}
				}
				world[i] = 0
			case 2: // word step
				fmt.Println(i, "data")
				i += wordSize
			case 3: // add
				fmt.Println(i, "add")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v += getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += wordSize
			case 4: // sub
				fmt.Println(i, "sub")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v -= getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += wordSize
			case 5: // mul
				fmt.Println(i, "mul")
				var v uint64
				if i > wordSize {
					v = getUint64(world[i-wordSize:], wordSize)
				}
				v *= getUint64(world[i:], wordSize)
				putUint64(world[i:], wordSize, v)
				i += wordSize
			case 6: // div
				fmt.Println(i, "div")
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
				i += wordSize
			default:
				world[i] = 0
			}
		}
	}
}
