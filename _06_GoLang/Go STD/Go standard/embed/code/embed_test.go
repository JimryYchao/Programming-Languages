package gostd

import (
	"embed"
	"io/fs"
	"testing"
)

//go:embed img.png
var img []byte

func TestEmbedImg(t *testing.T) {
	log(len(img))
}

//go:embed hello.txt
var hello string

func TestEmbedToString(t *testing.T) {
	log(hello)
}

//go:embed code.go hello.txt
var embedfs embed.FS

//go:embed \\
var emAll embed.FS

func TestEmbedFS(t *testing.T) {
	if dirs, err := fs.ReadDir(embedfs, "."); err == nil {
		for _, d := range dirs {
			log(d.Name())
		}
	}
	log(" ")
	if dirs, err := fs.ReadDir(emAll, "."); err == nil {
		for _, d := range dirs {
			log(d.Name())
		}
	}
}
