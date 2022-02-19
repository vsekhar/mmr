package main

import (
	"fmt"
	"math/bits"

	"github.com/vsekhar/mmr/internal/bruteforce"
)

func log2(x int) int { return bits.Len(uint(x)) - 1 }

func alloc(x int) int { return log2(x) - 1 }

func main() {
	for i := 1000; i < 100000; i += 449 {
		m := bruteforce.New(i)
		fmt.Printf("%d, %d, %d, %v\n", i, alloc(i), len(m.Peaks), m.Peaks)
		if len(m.Peaks) > alloc(i) {
			fmt.Println("stop")
			break
		}
	}
}
