package trie

type Trie struct {
	Root *Node
}

type Node struct {
	Chars  map[string]*Node
	IsEnd  bool
	Parent *Node
}

func InitTrie() *Trie {
	return &Trie{Root: &Node{Chars: make(map[string]*Node, 26)}}
}

func (t *Trie) Insert(word string) {
	current := t.Root
	for i := 0; i < len(word); i++ {
		_, exists := current.Chars[string(word[i])]
		if !exists {
			newNode := &Node{Chars: make(map[string]*Node, 26)}
			newNode.Parent = current
			current.Chars[string(word[i])] = newNode
		}
		current = current.Chars[string(word[i])]
	}
	current.IsEnd = true
}

func (t *Trie) Search(word string) bool {
	current := t.Root
	for i := 0; i < len(word); i++ {
		_, exists := current.Chars[string(word[i])]
		if !exists {
			return false
		}
		current = current.Chars[string(word[i])]
	}
	return current.IsEnd
}
