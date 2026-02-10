package trie

type TrieNode struct {
	Name     string `json:"Name"`
	Dir      bool   `json:"Dir"`
	OnDisk   bool   `json:"-"`
	Children map[string]*TrieNode
}
