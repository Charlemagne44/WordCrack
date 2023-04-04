package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"wordcrack/trie"

	"github.com/eiannone/keyboard"
	"golang.org/x/exp/slices"
)

const TOPWORDS = 5
const GAMESIZE = 4
const TABSPACING = 4

//go:embed resources/scrabble.json
var f embed.FS

type Game struct {
	Tiles       [GAMESIZE][GAMESIZE]Tile
	Valid_Words []string
	Best_Words  [][]string
	Trie        trie.Trie
	Visited     [GAMESIZE][GAMESIZE]bool
}

type Tile struct {
	Value   string
	Visited bool
}

func main() {
	// init game and board
	game := Game{}
	game.Best_Words = make([][]string, GAMESIZE)
	for i := 0; i < GAMESIZE; i++ {
		game.Best_Words[i] = make([]string, GAMESIZE)
	}

	// load in the english dictionary json
	data, _ := f.ReadFile("resources/scrabble.json")

	// unmarshal the dictionary into a list
	var dictionary []string
	err := json.Unmarshal(data, &dictionary)
	if err != nil {
		fmt.Printf("Unmarshal %v\n", err)
	}

	// load the dictionary into a trie
	trie := trie.InitTrie()
	for _, word := range dictionary {
		if strings.Contains(word, "-") {
			continue
		}
		trie.Insert(strings.ToLower(word))
	}
	game.Trie = *trie

	// load in tiles
	fmt.Println("Type game grid:")
	game.LoadWords()

	// for each cell, explore each path appending valid words-loc as you find them via backtracing
	for row := 0; row < GAMESIZE; row++ {
		for col := 0; col < GAMESIZE; col++ {
			// find all valid words from origin tile
			game.backtrack(row, col, *game.Trie.Root, "")

			// sort them by size, largest first
			sortListBySize(game.Valid_Words)

			// print the top 5 solutions first and display their path
			game.printTopOptions(row, col)

			// append the highest scoring Non repeating word IF it is not already in the matrix
			game.insertBestWord(row, col)

			// flush the found words and visited words
			for row := range game.Visited {
				for col := range game.Visited {
					game.Visited[row][col] = false
				}
			}
			game.Valid_Words = []string{}
		}
	}
}

func (g *Game) printTopOptions(row, col int) {
	if len(g.Valid_Words) > TOPWORDS {
		PrettyPrint(g.Valid_Words[len(g.Valid_Words)-TOPWORDS-1:len(g.Valid_Words)-1], row, col)
	} else {
		PrettyPrint(g.Valid_Words, row, col)
	}
	fmt.Print("\n")
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
	// get the current char from the tile at row, col
	tileChar := string(g.Tiles[row][col].Value)
	// check if char appears at current node of trie
	_, exists := trie.Chars[tileChar]
	if !exists {
		// not a path to a valid word, return
		return
	}
	// append to current recursively constructed word and append char from current tile
	word = word + g.Tiles[row][col].Value
	// set trie node to the existing chile not at that specific character
	trie = *trie.Chars[tileChar]
	// mark the visited matrix true to avoid backtracking repeats
	g.Visited[row][col] = true
	// check for the existence of the constructed word within the scrabble dict trie
	if g.Trie.Search(word) {
		if !slices.Contains(g.Valid_Words, word) && len(word) > 2 {
			g.Valid_Words = append(g.Valid_Words, word)
		}
	}
	// create delta arrays to explore surrounding cells
	drow := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dcol := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	// for every adjacent cell within the matrix and not visited, recurse this backtrack function
	for i := 0; i < 8; i++ {
		nextrow := row + drow[i]
		nextcol := col + dcol[i]
		if nextrow < GAMESIZE && nextrow >= 0 && nextcol < GAMESIZE && nextcol >= 0 {
			if !g.Visited[nextrow][nextcol] {
				g.backtrack(nextrow, nextcol, trie, word)
			}
		}
	}
	// cleanup the visited array at the current cell post tracking and go up a level in Trie node
	g.Visited[row][col] = false
	trie = *trie.Parent
}

func PrettyPrint(row []string, i, j int) {
	writer := tabwriter.NewWriter(os.Stdout, 0, TABSPACING, 1, '\t', 0)
	defer writer.Flush()
	fmt.Fprintf(writer, "Coords: %d %d  -  ", i, j)
	for _, element := range row {
		fmt.Fprintf(writer, "%s\t", element)
	}
	fmt.Printf("\n")
}

func (g *Game) LoadTestWords() {
	// load in a predefined tile grid for testing
	g.Tiles = [GAMESIZE][GAMESIZE]Tile{
		{Tile{Value: "e", Visited: false}, Tile{Value: "a", Visited: false}, Tile{Value: "t", Visited: false}, Tile{Value: "c", Visited: false}},
		{Tile{Value: "e", Visited: false}, Tile{Value: "s", Visited: false}, Tile{Value: "i", Visited: false}, Tile{Value: "a", Visited: false}},
		{Tile{Value: "m", Visited: false}, Tile{Value: "h", Visited: false}, Tile{Value: "p", Visited: false}, Tile{Value: "s", Visited: false}},
		{Tile{Value: "k", Visited: false}, Tile{Value: "y", Visited: false}, Tile{Value: "b", Visited: false}, Tile{Value: "i", Visited: false}},
	}
}

func (g *Game) LoadWords() {
	// laod in the tile character values
	for i := 0; i < GAMESIZE; i++ {
		for j := 0; j < GAMESIZE; j++ {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			fmt.Printf("You pressed: %q\r\n", char)
			g.Tiles[i][j].Value = string(char)
		}
	}
}
