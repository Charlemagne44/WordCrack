package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"wordcrack/trie"

	"golang.org/x/exp/slices"
)

type Game struct {
	Tiles       [4][4]Tile
	Valid_Words []string
	Best_Words  [][]string
	trie        trie.Trie
	Visited     [][]bool
}

type Tile struct {
	Value   string
	Visited bool
}

type KeyVal struct {
	Key   string
	Value []Tile
}

func main() {
	// init game and board
	game := Game{}
	game.Visited = make([][]bool, 4)
	game.Best_Words = make([][]string, 4)
	for i := 0; i < 4; i++ {
		game.Visited[i] = make([]bool, 4)
		game.Best_Words[i] = make([]string, 4)
	}

	// load in the english dictionary json
	jsonFile, err := os.Open("resources/10000.json")
	if err != nil {
		fmt.Printf("Open dict json: %v\n", err)
	}
	defer jsonFile.Close()

	var dictionary []string
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("Read json: %v\n", err)
	}
	err = json.Unmarshal(byteValue, &dictionary)
	if err != nil {
		fmt.Printf("Unmarshal %v\n", err)
	}

	// load the dictionary into a trie
	trie := trie.InitTrie()
	for _, word := range dictionary {
		// disallow hyphens
		if strings.Contains(word, "-") {
			continue
		}
		trie.Insert(word)
	}
	game.trie = *trie

	// load in tiles
	game.LoadTestWords()
	// game.LoadWords()

	// for each cell, explore each path appending valid words-loc as you find them with the trie
	for row := 0; row < 1; row++ {
		for col := 0; col < 1; col++ {
			// find all valid words from origin tile
			game.backtrack(row, col, *game.trie.Root, "")

			// sort them by size, largest first
			sort.Slice(game.Valid_Words, func(i, j int) bool {
				return len(game.Valid_Words[i]) < len(game.Valid_Words[j])
			})
			// print the largest solutions first and display their path
			fmt.Print("COORDS: ", row, col)
			fmt.Println(game.Valid_Words)

			// append the highest scoring word to the best word grid
			if len(game.Valid_Words) > 1 {
				game.Best_Words[row][col] = game.Valid_Words[len(game.Valid_Words)-1]
			}

			// flush the found words and visited  words
			for row := range game.Visited {
				for col := range game.Visited {
					game.Visited[row][col] = false
				}
			}
			game.Valid_Words = []string{}
		}
	}

	PrettyPrint(game.Best_Words)
}

func (g *Game) backtrack(x, y int, trie trie.Node, word string) {
	fmt.Println("currently at x,  y:  ", x, y)
	tileChar := string(g.Tiles[x][y].Value)
	fmt.Println("tile char val: ", tileChar)
	_, exists := trie.Chars[tileChar]
	if !exists {
		fmt.Println("char: " + tileChar + " doesn't exist for word: " + word)
		return
	}
	word = word + g.Tiles[x][y].Value
	fmt.Println("Word constructed:", word)
	trie = *trie.Chars[tileChar]
	g.Visited[x][y] = true
	if g.trie.Search(word) {
		if !slices.Contains(g.Valid_Words, word) && len(word) > 2 {
			fmt.Println("Is word:  ", word)
			g.Valid_Words = append(g.Valid_Words, word)
		}
	}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	for i := 0; i < 8; i++ {
		nextx := x + dx[i]
		nexty := y + dy[i]
		if nextx < 4 && nextx >= 0 && nexty < 4 && nexty >= 0 {
			if !g.Visited[nexty][nextx] {
				fmt.Println("Going to: ", nextx, nexty)
				g.backtrack(nextx, nexty, trie, word)
			}
		}
	}
	g.Visited[x][y] = false
	trie = *trie.Parent
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
	return string(s)
}

func (g *Game) LoadTestWords() {
	// load in a predefined tile grid for testing
	g.Tiles = [4][4]Tile{
		{Tile{Value: "e", Visited: false}, Tile{Value: "a", Visited: false}, Tile{Value: "t", Visited: false}, Tile{Value: "c", Visited: false}},
		{Tile{Value: "e", Visited: false}, Tile{Value: "s", Visited: false}, Tile{Value: "i", Visited: false}, Tile{Value: "a", Visited: false}},
		{Tile{Value: "m", Visited: false}, Tile{Value: "h", Visited: false}, Tile{Value: "p", Visited: false}, Tile{Value: "s", Visited: false}},
		{Tile{Value: "k", Visited: false}, Tile{Value: "y", Visited: false}, Tile{Value: "b", Visited: false}, Tile{Value: "i", Visited: false}},
	}
}

func (g *Game) LoadWords() {
	// laod in the tile character values
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var char string
			fmt.Scanln(&char)
			g.Tiles[i][j].Value = char
		}
	}
}
