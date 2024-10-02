package main

import (
	"fmt"
	"math/rand"
	"testing"

	"xaavian.com/Golfle/wordLists"
)

var (
	wordsFive  = wordLists.WordsFive
	wordsAll   = append(wordLists.WordsStart, wordLists.WordsRest...)
	x          = &wordsAll //Change this to what you want to benchmark
	size       = 10        // How many test iterations really
	benchStart = makeList(size)
	benchEnd   = makeList(size)
)

func makeList(n int) []string {
	l := make([]string, n)
	q := len(*x)
	for i := 0; i < n; i++ {
		l = append(l, (*x)[rand.Intn(q)])
	}
	return l
}

func BenchmarkFind1(b *testing.B) {
	fmt.Println(*x)
	for n := 0; n < b.N; n++ {
		shortestPath(benchStart[n], benchEnd[n], 2, *x)
	}
}

func BenchmarkFind2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		shortestPath2(benchStart[n], benchEnd[n], 2, *x)
	}
}
