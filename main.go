package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = "<!DOCTYPE html>\n<html>\n\u0020\u0020<head>\n\u0020\u0020\u0020\u0020<meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\">\n\u0020\u0020\u0020\u0020<title>Markdown Preview Tool</title>\n\u0020\u0020</head>\n\u0020\u0020<body>\n"
	footer = "\n\u0020\u0020</body>\n</html>\n"
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHtml(outName, htmlData)
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Write HTML to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHtml(outName string, data []byte) error {
	//Write the bytes to the file
	return os.WriteFile(outName, data, 0644)
}
