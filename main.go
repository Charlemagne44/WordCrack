package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"wordcrack/trie"

	"github.com/eiannone/keyboard"
	"golang.org/x/exp/slices"
)

type Game struct {
	Tiles       [4][4]Tile
	Valid_Words []string
	Best_Words  [][]string
	Trie        trie.Trie
	Visited     [4][4]bool
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
	game.Best_Words = make([][]string, 4)
	for i := 0; i < 4; i++ {
		game.Best_Words[i] = make([]string, 4)
	}

	// load in the english dictionary json
	jsonFile, err := os.Open("resources/dictionary.json")
	if err != nil {
		fmt.Printf("Open dict json: %v\n", err)
	}
	defer jsonFile.Close()

	// unmarshal the dictionary into a list
	var dictionary map[string]string
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
	for word, _ := range dictionary {
		// disallow hyphens
		if strings.Contains(word, "-") {
			continue
		}
		trie.Insert(word)
	}
	game.Trie = *trie

	// load in tiles
	// game.LoadTestWords()
	fmt.Println("Type game grid:")
	game.LoadWords()

	var allwords []string
	// for each cell, explore each path appending valid words-loc as you find them via backtracing
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			// find all valid words from origin tile
			game.backtrack(row, col, *game.Trie.Root, "")

			// sort them by size, largest first
			sortListBySize(game.Valid_Words)

			// print the largest solutions first and display their path
			fmt.Print("COORDS: ", row, col)
			fmt.Println(game.Valid_Words)

			// append the highest scoring Non repeating word IF it is not already in the matrix
			game.insertBestWord(row, col)

			// flush the found words and visited words
			for row := range game.Visited {
				for col := range game.Visited {
					game.Visited[row][col] = false
				}
			}
			allwords = append(allwords, game.Valid_Words...)
			game.Valid_Words = []string{}
		}
	}

	// sortListBySize(allwords)
	// fmt.Println("ALL BEST:")
	// for _, element := range allwords {
	// 	fmt.Println(element)
	// }
	// PrettyPrint(game.Best_Words)
}

func sortListBySize(list []string) {
	sort.Slice(list, func(i, j int) bool {
		return len(list[i]) < len(list[j])
	})
}

func (g *Game) insertBestWord(row, col int) {
	for i := len(g.Valid_Words) - 1; i >= 0; i-- {
		curr_best := g.Valid_Words[i]
		best := true
		for _, words := range g.Best_Words {
			if slices.Contains(words, curr_best) {
				best = false
			}
		}
		if best {
			g.Best_Words[row][col] = curr_best
			return
		} else {
			continue
		}
	}
}

func (g *Game) backtrack(row, col int, trie trie.Node, word string) {
	// fmt.Println("currently at row,  col:  ", row, col)
	tileChar := string(g.Tiles[row][col].Value)
	// fmt.Println("tile char val: ", tileChar)
	_, exists := trie.Chars[tileChar]
	if !exists {
		// fmt.Println("char: " + tileChar + " doesn't exist for word: " + word)
		return
	}
	word = word + g.Tiles[row][col].Value
	// fmt.Println("Word constructed:", word)
	trie = *trie.Chars[tileChar]
	g.Visited[row][col] = true
	if g.Trie.Search(word) {
		if !slices.Contains(g.Valid_Words, word) && len(word) > 2 {
			// fmt.Println("Is word:  ", word)
			g.Valid_Words = append(g.Valid_Words, word)
		}
	}
	drow := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dcol := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	for i := 0; i < 8; i++ {
		nextrow := row + drow[i]
		nextcol := col + dcol[i]
		if nextrow < 4 && nextrow >= 0 && nextcol < 4 && nextcol >= 0 {
			if !g.Visited[nextrow][nextcol] {
				// fmt.Println("Going to row, col: ", nextrow, nextcol)
				g.backtrack(nextrow, nextcol, trie, word)
			}
		}
	}
	g.Visited[row][col] = false
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
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			fmt.Printf("You pressed: %q\r\n", char)
			g.Tiles[i][j].Value = string(char)
		}
	}
}
