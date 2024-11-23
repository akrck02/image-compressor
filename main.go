package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const MIME_TYPES_JPEG string = "image/jpeg"
const MIME_TYPES_PNG string = "image/png"
const THUMBNAIL_DEFAULT_WIDTH int = 300
const THUMBNAIL_DEFAULT_SUFFIX string = "min"

func main() {
	thumbnailImages("resources", "thumbnails", THUMBNAIL_DEFAULT_WIDTH, THUMBNAIL_DEFAULT_SUFFIX)
}

// Thumbnail the images in the given path to the specified width and suffix
func thumbnailImages(path string, creationPath string, width int, suffix string) {

	// get the files in the directory
	items, error := os.ReadDir(path)
	if error != nil {
		log.Fatal("Cannot access the given path", path)
	}

	// loop through the files and thumbnail them
	for _, item := range items {
		thumbnailDirEntry(path, creationPath, item, width, suffix)
	}

}

// Thumbnail the files in the given directory thumbnailDirEntry
// and its subdirectories to the specified width
func thumbnailDirEntry(path string, creationPath string, item fs.DirEntry, width int, suffix string) {

	currentPath := path + "/" + item.Name()
	currentCreationPath := creationPath + "/" + item.Name()

	log.Println("Checking route " + currentPath)
	log.Println("Checking creation route " + currentCreationPath)

	// if the item is a file, thumbnail it
	if !item.IsDir() {
		thumbnailFile(currentPath, currentCreationPath, width, suffix)
		return
	}

	// if the item is a directory, thumbnail the files in it
	subItems, subItemError := os.ReadDir(item.Name())
	if subItemError != nil {
		subItems = []fs.DirEntry{}
		log.Print("Cannot open " + item.Name())
	}

	for _, subItem := range subItems {
		currentPath = currentPath + subItem.Name()
		currentCreationPath = currentCreationPath + subItem.Name()
		thumbnailDirEntry(currentPath, currentCreationPath, subItem, width, suffix)
	}
}

// ThumbnailFile creates a resized image from the file and writes it to
// another file with suffix added to the original file name.
func thumbnailFile(path string, creationPath string, width int, suffix string) {

	// Determine the mimetype of the file
	// and create the thumbnail
	// note: the mimetype is determined by the file extension
	var extension, _ = filepath.Ext(path), "."
	var mimetype string = ""

	switch extension {
	case ".jpg", ".jpeg":
		mimetype = MIME_TYPES_JPEG
	case ".png":
		mimetype = MIME_TYPES_PNG
	}

	log.Print("Extension: " + extension)
	log.Print("MimeType: " + mimetype)

	// if the file is not an image, return
	if mimetype == "" {
		log.Print("The file " + path + " has not valid extension (" + extension + ").")
		return
	}

	// If the directory does not exist, create it
	directory := filepath.Dir(creationPath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, os.ModePerm)
	}

	// Create the output file
	var thumbnailPath string = strings.TrimSuffix(creationPath, extension) + "-" + suffix + extension
	output, err := os.Create(thumbnailPath)
	panicIfNeeded(err)

	log.Println("Generating thumbnail " + thumbnailPath)

	// Generate the thumbnail
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		log.Print("Cannot create file for thumbnail " + thumbnailPath)
		return
	}

	err = thumbnail(thumbnailFile, output, mimetype, width)
	panicIfNeeded(err)

	// close the output file
	defer thumbnailFile.Close()
	err = output.Close()
	panicIfNeeded(err)
}

// Thumbnail creates a resized image from the reader and writes it to
// the writer. The mimetype determines how the image will be decoded
// and must be either "image/jpeg" or "image/png". The desired width
// of the thumbnail is specified in pixels, and the resulting height
// will be calculated to preserve the aspect ratio.
// ..................................................................
// The original code of the resize function was taken
// from https://roeber.dev/posts/resize-an-image-in-go/
// thanks to the author (Roeber.dev) for the code.
func thumbnail(r io.Reader, w io.Writer, mimetype string, width int) error {

	var src image.Image
	var err error
	print(mimetype)

	switch mimetype {
	case MIME_TYPES_JPEG:
		src, err = jpeg.Decode(r)
	case MIME_TYPES_PNG:
		src, err = png.Decode(r)
	}

	if err != nil {
		return err
	}

	ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	height := int(math.Round(float64(width) * ratio))
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	err = jpeg.Encode(w, dst, nil)
	if err != nil {
		return err
	}

	return nil
}

// Check function to handle errors
func panicIfNeeded(err error) {
	if err != nil {
		panic(err)
	}
}
