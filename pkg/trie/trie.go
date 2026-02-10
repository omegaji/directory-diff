package trie

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const PATH_SPLITTER = string(os.PathSeparator)
const TRIE_BACKUP_PATH = ".triebackup.json"

type Trie struct {
	Fs                  fs.FS  `json:"-"`
	RootPath            string `json:"rootPath"`
	Root                *TrieNode
	LastCommitTimestamp uint64 `json:"lastCommitTimestamp"`
}

func (t *Trie) Add(path string, isDir bool) error {
	if t.Root == nil {
		return fmt.Errorf("trie is not initialized")
	}

	if path == TRIE_BACKUP_PATH {
		return nil
	}

	currNode := t.Root
	splitPath := strings.Split(path, PATH_SPLITTER)

	for i, n := range splitPath {
		child, ok := currNode.Children[n]
		if !ok {
			tNode := TrieNode{
				Name:     n,
				Dir:      true,
				Children: make(map[string]*TrieNode),
			}

			if i == len(splitPath)-1 {
				tNode.Dir = isDir
			}

			child = &tNode
			currNode.Children[n] = child
		}

		child.OnDisk = true
		currNode = child
	}

	return nil
}

func (t *Trie) Has(name string) bool {
	currNode := t.Root
	splitNames := strings.Split(name, PATH_SPLITTER)

	for _, n := range splitNames {
		child, ok := currNode.Children[n]
		if !ok {
			return false
		}
		currNode = child
	}
	return true
}

func (t *Trie) Compare(logChanges bool) {
	if t.Fs == nil {
		log.Fatal("FS not initialized")
	}

	err := fs.WalkDir(t.Fs, ".", func(path string, d fs.DirEntry, err error) error {
		exists := t.Has(path)
		t.Add(path, d.IsDir())

		if exists {
			fInfo, err := d.Info()
			if err != nil {
				log.Fatal(err)
			}

			if logChanges {
				if fInfo.ModTime().After(time.Unix(int64(t.LastCommitTimestamp), 0)) {
					log.Printf("Modified: %s", path)
				}
			}
		} else if logChanges && path != TRIE_BACKUP_PATH {
			log.Printf("Added: %s", path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if logChanges {
		err = t.Walk(func(node *TrieNode, path string) error {
			if !node.OnDisk {
				log.Printf("Deleted: %s", path)
			}
			return nil
		})
	}

	if err != nil {
		log.Fatal(err)
	}
}

func (t *Trie) prune() {
	stack := []*TrieNode{t.Root}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for k, v := range node.Children {
			if !v.OnDisk {
				delete(node.Children, k)
			} else {
				stack = append(stack, v)
			}
		}
	}
}

func (t *Trie) Save() {
	t.prune()
	t.LastCommitTimestamp = uint64(time.Now().Unix())

	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to save tree: %s", err.Error()))
	}

	err = os.WriteFile(TRIE_BACKUP_PATH, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// walk is not deterministic in terms of order as it is not sorted
func (t *Trie) Walk(fn func(n *TrieNode, p string) error) error {
	return t.walk(t.Root, t.RootPath, fn)
}

func (t *Trie) walk(node *TrieNode, path string, fn func(n *TrieNode, p string) error) error {
	for _, v := range node.Children {
		err := fn(v, filepath.Join(path, v.Name))
		if err != nil {
			return err
		}

		err = t.walk(v, filepath.Join(path, v.Name), fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func Load() Trie {
	data, err := os.ReadFile(TRIE_BACKUP_PATH)
	if err != nil {
		log.Fatal(err)
	}

	var root Trie
	err = json.Unmarshal(data, &root)
	if err != nil {
		log.Fatal(err)
	}

	root.Fs = os.DirFS(root.RootPath)
	return root
}

func Trie_Init(rootPath string, f fs.FS) (Trie, error) {
	t := Trie{
		RootPath: rootPath,
		Root: &TrieNode{
			Name:     "*",
			Children: make(map[string]*TrieNode),
			Dir:      true,
		},
		Fs: f,
	}

	err := fs.WalkDir(t.Fs, ".", func(path string, d fs.DirEntry, err error) error {
		return t.Add(path, d.IsDir())
	})
	if err != nil {
		return Trie{}, err
	}
	t.LastCommitTimestamp = uint64(time.Now().Unix())
	return t, nil
}
