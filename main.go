package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/HasinduLanka/console"
)

var phpfilesDir string
var origin string
var timeStamp string
var outputFilePath string

var sitemapBuilder strings.Builder

func main() {
	console.GlobalWriter = console.NewWriterToStandardOutput()
	console.GlobalReader = console.NewReaderFromFile("cli")

	console.Print("SitemapGenForPHP")

	sitemapBuilder = strings.Builder{}

	phpfilesDir = console.ReadLine() + "/"
	origin = console.ReadLine() + "/"
	outputFilePath = console.ReadLine()

	timeStamp = time.Now().Format("2006-01-02")

	sitemapBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
`)

	// Make sure folder exists
	if _, err := os.Stat(phpfilesDir); os.IsNotExist(err) {
		console.Log("Folder not found")
		os.Mkdir(phpfilesDir, 0777)
	}

	// Read all files in the folder recursively

	e := filepath.Walk(phpfilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			console.Log(err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".php") {
			route := strings.TrimSuffix(strings.TrimPrefix(path, phpfilesDir), ".php") + "/"
			processFile(path, route)
		}

		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	sitemapBuilder.WriteString(`</urlset>`)

	outputFile := console.NewWriterToFile(outputFilePath)
	outputFile.PrintInline(sitemapBuilder.String())

	outputFile.Buff.Flush()
	outputFile.CloseWritterFile()
}

func processFile(path string, route string) {
	console.Print(path)

	// Read the file
	fileBytes, fileErr := os.ReadFile(path)

	if fileErr != nil {
		console.Log(fileErr)
		return
	}

	file := string(fileBytes)

	funcRgx := regexp.MustCompile(`public function ([a-zA-Z0-9]*?)\(`)
	funcMatches := funcRgx.FindAllStringSubmatch(file, len(file))

	if len(funcMatches) == 0 {
		return
	}

	for _, funcMatch := range funcMatches {
		funcName := funcMatch[1]
		console.Print(route + funcName)

		sitemapBuilder.WriteString(`
	<url>
		<loc>` + origin + route + funcName + `</loc>
		<lastmod>` + timeStamp + `</lastmod>
	</url>
`)
	}

}
