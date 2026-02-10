package main

import (
	"flag"
	"log"
	"os"

	"github.com/omegaji/directory-diff/pkg/trie"
)

func main() {

	commit := flag.Bool("commit", false, "to update saved tree")
	status := flag.Bool("status", false, "gives the status of what has changed in the residing directory")

	flag.Parse()
	var t trie.Trie

	_, err := os.Stat(trie.TRIE_BACKUP_PATH)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if err == nil {
		if commit == nil && status == nil {
			log.Fatal("No action provided, returning early. Common actions are either 'status' or 'commit'")
		}

		t = trie.Load()
		if status != nil && *status {
			t.Compare(true)
		} else {
			t.Compare(false)
		}

		if commit != nil && *commit {
			t.Save()
		}
	} else {
		if commit != nil && *commit {
			log.Fatal("you cannot commit a tree if it does not exist, first create a tree")
		}
		t, err = trie.Trie_Init("./", os.DirFS("./"))
		if err != nil {
			log.Fatal(err)
		}

		t.Save()
		log.Println("Saved state")
	}

}
