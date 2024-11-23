package main

import (
	"log"
	"os"
	"strconv"

	"github.com/akrck02/image-compressor/image"
)

func main() {

	args := os.Args[1:]
	if len(args) < 2 {
		end("width is not a number.")
	}

	path := args[0]
	creationPath := args[1]

	width := image.THUMBNAIL_DEFAULT_WIDTH
	if len(args) > 3 {
		w, err := strconv.Atoi(args[2])
		if err != nil {
			end("width is not a number.")
		}
		width = w
	}

	suffix := ""
	if len(args) > 4 {
		suffix = args[3]
	}

	// string to int
	image.Thumbnail(path, creationPath, width, suffix)
}

func end(text string) {
	log.Print(text)
	help()
	log.Fatal("")
}

func help() {
	log.Print("Command usage: (arguments in correct order)")
	log.Print("thumbnail <path> <new-path> [new-width] [new-suffix]")
}
