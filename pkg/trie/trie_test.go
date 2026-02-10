package trie

import (
	"fmt"
	"io/fs"
	"log"
	"testing"
	"testing/fstest"
)

func TestTrieInit(t *testing.T) {
	tests := []struct {
		name string
		fs   fstest.MapFS
	}{

		{"fs-without-dir", fstest.MapFS{
			"hello.txt":             &fstest.MapFile{Data: []byte("hello.txt file")},
			"g,dhfdhfdhjiuruirrei.": &fstest.MapFile{Data: []byte("djdjfd")},
			"dir-hello":             &fstest.MapFile{Mode: fs.ModeDir},
			"dir-hello/dir-go":      &fstest.MapFile{Mode: fs.ModeDir},
		}},

		{"fs-with-dir", fstest.MapFS{
			"hello.txt":                                   &fstest.MapFile{Data: []byte("hello.txt file")},
			"g,dhfdhfdhjiuruirrei.":                       &fstest.MapFile{Data: []byte("djdjfd")},
			"dir-hello":                                   &fstest.MapFile{Mode: fs.ModeDir},
			"dir-hello/dir-howdydo":                       &fstest.MapFile{Mode: fs.ModeDir},
			"dir-hello/dir-howdydo/dir-howdyo2":           &fstest.MapFile{Mode: fs.ModeDir},
			"dir-hello/howdy.tx":                          &fstest.MapFile{Data: []byte("dfdfdfwerewr")},
			"dir-hello/dir-howdydo/dir-howdyo2/hello.txt": &fstest.MapFile{Data: []byte("2222dfdfdfwerewr")},
		}},
	}

	for _, tt := range tests {
		fmt.Printf("starting testing for: %s", tt.name)
		trie, err := Trie_Init(".", tt.fs)
		if err != nil {
			t.Fatalf("error encountered when creating Trie: %s", err.Error())
		}

		fileNames, err := tt.fs.Glob("*")
		if err != nil {
			log.Fatalf("error encousntered when fetching files from mapFs")
		}

		for _, fileName := range fileNames {
			_, ok := trie.Root.Children[fileName]
			if !ok {
				t.Fatalf("Failed to find %s in fs", fileName)
			}
		}
	}
}
