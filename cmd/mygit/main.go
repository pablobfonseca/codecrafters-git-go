package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file <type> <object>\n")
			os.Exit(1)
		}
		// flag := os.Args[2]
		object := os.Args[3]

		dirName := object[:2]
		blobSha := object[2:]

		buff, err := os.ReadFile(path.Join(".git", "objects", dirName, blobSha))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file content: %s\n", err)
			os.Exit(1)
		}

		b := bytes.NewReader(buff)
		r, err := zlib.NewReader(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "zlib error: %s\n", err)
			os.Exit(1)
		}
		defer r.Close()

		var out bytes.Buffer
		_, err = io.Copy(&out, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading decompressed data: %s\n", err)
			os.Exit(1)
		}
		data := out.String()
		nullIdx := strings.IndexByte(data, '\x00')
		if nullIdx == -1 {
			fmt.Fprintf(os.Stderr, "invalid git object format\n")
			os.Exit(1)
		}

		content := data[nullIdx+1:]
		fmt.Print(content)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
