package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"

	"xaavian.com/Golfle/wordLists"
)

var (
	hardMode   = flag.Bool("hard", false, "Sets game difficulty to hard mode")
	normalMode = flag.Bool("normal", false, "Sets game difficulty to normal mode (Default)")
	diffMax    = 2
	words      = make([]string, 0)
	dictionary = make(map[string]bool)
	start      = ""
	end        = ""
)

// Finds the shortest path you can take between start and end, only making MaxDiff changes
func shortestPath2(start string, end string, MaxDiff int, f []string) []string {

	type node struct {
		word     string
		path     []string
		distance int
	}

	g := make([]node, len(f))
	for h := 0; h < len(f); h++ {
		g[h].word = f[h]
		g[h].distance = -1
	}
	var startint int
	for i := 0; ; i++ {
		if g[i].word == start {
			startint = i
			break
		}
	}
	tmp := g[0]
	g[0] = g[startint]
	g[startint] = tmp

	currentDistance := 0
	g[0].distance = currentDistance

	g[0].path = append(g[0].path, g[0].word)

	//fmt.Printf("%s => %s\n", start, end)

	// O(n^3) is attrocious, just do DFS like an adult
	for z := 0; z < 10; z++ {
		for i := 0; i < len(g); i++ {
			if g[i].distance == currentDistance {
				for j := 0; j < len(g); j++ {
					if i != j && g[j].distance == -1 && check(g[i].word, g[j].word, MaxDiff) {
						g[j].path = append(g[j].path, g[i].path...)
						g[j].path = append(g[j].path, g[j].word)
						g[j].distance = g[i].distance + 1
						if g[j].word == end {
							return g[j].path
						}
					}
				}

			}
		}
		currentDistance++
	}
	fmt.Println("No Path Found.")
	return make([]string, 0)
}

// Finds the shortest path you can take between start and end, only making MaxDiff changes
func shortestPath(start string, end string, MaxDiff int, f []string) []string {

	type node struct {
		word    string
		path    []string // Storing each path to get to each string *cannot* be memory efficient
		visited bool
	}

	searchQueue := make([]*node, 0) // Size 0 to start is fine honestly I doubt resizing the array hampers performance that badly

	// Does the compile optimize this? Like should I call len(f) once and store the value to avoid calling it again or would the compiler just fix it for me? len(f) doesnt change and is the same as len(g)
	listSize := len(f)

	g := make([]node, listSize)
	for h := 0; h < listSize; h++ {
		g[h].word = f[h]
	}

	// Idk why im doing this but I want to
	var startint int
	for i := 0; ; i++ {
		if g[i].word == start {
			startint = i
			break
		}
	}
	tmp := g[0]
	g[0] = g[startint]
	g[startint] = tmp

	g[0].path = append(g[0].path, g[0].word)
	g[0].visited = true
	searchQueue = append(searchQueue, &g[0])

	// I fucked up and now j is the outer loop variable
	for j := 0; j < len(searchQueue); j++ {
		for i := 0; i < listSize; i++ {
			if !g[i].visited && check(g[i].word, searchQueue[j].word, diffMax) {
				searchQueue = append(searchQueue, &g[i])
				g[i].visited = true
				g[i].path = append(searchQueue[j].path, g[i].word)
				if g[i].word == end {
					return g[i].path
				}
			}
		}
	}

	fmt.Println("No Path Found.")
	return make([]string, 0)
}

// Returns if the two given strings have less or equal to the maximum amount of allowed changes between them
// A change is defined as either adding a letter, removing a letter, or replacing a letter for another
// Kinda surprised this works, there may be an edge case for which it fails but I can't mathematically prove anything
func check(x string, y string, maxDiff int) bool {
	var rec func(x string, y string) int
	rec = func(x string, y string) int {
		//fmt.Printf("Comparing \"%s\" to \"%s\"\n", x, y)
		if len(x) == 0 {
			//fmt.Printf("Comparing \"%s\" to empty string, returning %d differences\n", y, len(y))
			return len(y)
		}
		if len(y) == 0 {
			//fmt.Printf("Comparing \"%s\" to empty string, returning %d differences\n", x, len(x))
			return len(x)
		}
		if len(x) > len(y) {
			tmp := x
			x = y
			y = tmp
		}
		biggestSubstring := func(x string, y string, start int) string {
			for i := 1; i <= start; i++ {
				for j := 0; j < i; j++ {
					//fmt.Printf("Looking for \"%s\"\n", string(x[j:j+start-i+1]))
					if strings.Index(y, string(x[j:j+start-i+1])) >= 0 {
						//fmt.Println("Found!")
						return string(x[j : j+start-i+1])
					}
				}
			}

			return ""
		}
		sub := biggestSubstring(x, y, len(x))
		if sub == "" {
			//fmt.Printf("No similarities found between \"%s\" and \"%s\", returning %d differences\n", x, y, len(y))
			return len(y)
		}
		xindex := strings.Index(x, sub)
		yindex := strings.Index(y, sub)
		//fmt.Printf("First portion: \"%s\" and \"%s\"\n", string(x[:xindex]), string(y[:yindex]))
		//fmt.Printf("Second portion: \"%s\" and \"%s\"\n", string(x[xindex+len(sub):]), string(y[yindex+len(sub):]))
		return rec(string(x[:xindex]), string(y[:yindex])) + rec(string(x[xindex+len(sub):]), string(y[yindex+len(sub):]))
	}
	foundDiff := rec(x, y)
	//fmt.Printf("Differences Found: %d\n", foundDiff)
	return foundDiff <= maxDiff
}

// Returns golf-themed string about your score
func findScore(par int, count int) string {
	diff := count - par
	if diff == 0 {
		return fmt.Sprintf("You shot par!")
	} else if diff == 1 {
		return fmt.Sprintf("You shot a bogey")
	} else if diff == 2 {
		return fmt.Sprintf("You shot a double bogey")
	} else if diff == 3 {
		return fmt.Sprintf("You shot a triple bogey")
	} else {
		return fmt.Sprintf("You shot +%d", diff)
	}
}

// Formats the path slice into a nicer string
func formatPath(x []string) string {
	a := x[0]
	for i := 1; i < len(x); i++ {
		a = fmt.Sprintf("%s => %s", a, x[i])
	}
	return a
}

func main() {

	flag.Parse()

	if (!*normalMode && !*hardMode) || *normalMode {
		words = wordLists.WordsFive
		size := len(words)
		start = words[rand.Intn(size)]
		end = words[rand.Intn(size)]
	} else {
		wordsStart := wordLists.WordsStart
		wordsRest := wordLists.WordsRest
		wordsStartSize := len(wordsStart)
		start = wordsStart[rand.Intn(wordsStartSize)]
		end = wordsStart[rand.Intn(wordsStartSize)]
		words = append(wordsRest, wordsStart...)
	}

	for i := 0; i < len(words); i++ {
		dictionary[words[i]] = true
	}

	shortestPath := shortestPath(start, end, diffMax, words)
	if len(shortestPath) == 1 {
		fmt.Printf("ERROR: Path not found between \"%s\" and \"%s\". Exiting game...\n", start, end)
		return
	}
	par := len(shortestPath) - 1
	x := ""
	count := 0

	fmt.Println("Enter 0 to give up.")
	fmt.Printf("Max letter change: %d\n", diffMax)
	fmt.Printf("Goal: %s => %s\n", start, end)
	fmt.Printf("Par %d\n", par)

	for {
		fmt.Scan(&x)
		if x == "0" { // Inputting 0 will end your game in case you get stuck
			fmt.Printf("Loser :( You took %d moves\n", count)
			break
		}
		if x == start {
			fmt.Println("Same word.")
			continue
		}
		if !dictionary[x] {
			fmt.Println("Word not in dictionary.")
			continue
		}
		if !check(x, start, 2) {
			fmt.Printf("Invalid Guess.\n")
			continue
		}
		count++
		fmt.Printf("%d: %s => %s\n", count, start, x)
		start = x
		if start == end {
			fmt.Printf("Winner! You've solved this in %d moves\n", count)
			fmt.Printf("%s\n", findScore(par, count))
			break
		}
	}
	fmt.Printf("Potential shortest path: %s\n", formatPath(shortestPath))
}
