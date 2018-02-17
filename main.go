package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	cliArgs := os.Args[1:]

	fullURL := cliArgs[0]

	parseWeb(fullURL, cliArgs)
}

func parseWeb(fullURL string, cliArgs []string) {
	pageParam := cliArgs[1]

	lastPage, err := strconv.Atoi(cliArgs[2])

	if err != nil {
		fmt.Println("Last page param error")
		os.Exit(-1)
	}

	for i := 1; i <= lastPage; i++ {
		fmt.Println("NEW PAGE ----------------")

		currentPageURL := fullURL + "?" + pageParam + strconv.Itoa(i)

		page, err := goquery.NewDocument(currentPageURL)

		if err != nil {
			fmt.Println("Error trying to get " + currentPageURL)
		}

		parsePage(page, cliArgs)
	}
}

func parsePage(page *goquery.Document, cliArgs []string) {
	host := cliArgs[3]
	classToSearch := cliArgs[4]
	fileClass := cliArgs[5]
	fileAttr := cliArgs[6]

	page.Find(classToSearch).Each(func(i int, s *goquery.Selection) {
		playValue, _ := s.Find(fileClass).First().Attr(fileAttr)
		file := strings.Split(playValue, "'")[1]
		fileURL := host
		if file != "" {
			fileURL = fileURL + file
		}

		name := s.Find("a").Text()
		fmt.Println(fileURL + " " + name)

		downloadFile(file, name)
	})

}

func downloadFile(url string, name string) {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		fmt.Println("Couldn't create folder for " + name)
	} else {
		remoteFile, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}

		fileChunks := strings.Split(name, "/")
		fileName := fileChunks[len(fileChunks)-1]

		defer remoteFile.Body.Close()

		out, err := os.Create(filepath.Join(name, fileName))

		if err != nil {
			log.Println(err)
		}
		defer out.Close()

		io.Copy(out, remoteFile.Body)
	}
}
