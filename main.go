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

	firstPage, err := strconv.Atoi(cliArgs[2])
	lastPage, err := strconv.Atoi(cliArgs[3])

	if err != nil {
		fmt.Println("----- Last page param error")
		os.Exit(-1)
	}

	for i := firstPage; i <= lastPage; i++ {
		currentPageURL := fullURL + "?" + pageParam + strconv.Itoa(i)
		fmt.Println(" ------------------- " + strconv.Itoa(i))
		fmt.Println("Getting page " + currentPageURL)
		fmt.Println(" ------------------- " + strconv.Itoa(i) + "\n")

		page, err := goquery.NewDocument(currentPageURL)

		if err != nil {
			fmt.Println("----- Error trying to get " + currentPageURL)
		}

		parsePage(page, cliArgs)
	}
}

func parsePage(page *goquery.Document, cliArgs []string) {
	host := cliArgs[4]
	classToSearch := cliArgs[5]
	fileClass := cliArgs[6]
	fileAttr := cliArgs[7]

	page.Find(classToSearch).Each(func(i int, s *goquery.Selection) {
		playValue, _ := s.Find(fileClass).First().Attr(fileAttr)
		file := strings.Split(playValue, "'")[1]
		fileURL := host
		if file != "" {
			fileURL = fileURL + file
		}

		name := s.Find("a").Text()
		fmt.Println(fileURL)

		downloadFile(fileURL, name, cliArgs)
	})

}

func downloadFile(url string, name string, cliArgs []string) {
	basicFolder := cliArgs[8]
	err := os.Mkdir(filepath.Join("sounds", name), os.ModePerm)

	remoteFile, err := http.Get(url)
	if err != nil {
		log.Println(err)
		fmt.Println("----- Error getting page " + url)
	}

	fileChunks := strings.Split(url, "/")
	fileName := fileChunks[len(fileChunks)-1]

	defer remoteFile.Body.Close()

	fileFullPath := filepath.Join(basicFolder, filepath.Join(name, fileName))

	fmt.Println(fileFullPath + "\n")

	out, err := os.Create(fileFullPath)

	if err != nil {
		log.Println(err)
		fmt.Println("----- Error creating / downloading file " + url)
	}
	defer out.Close()

	io.Copy(out, remoteFile.Body)
}
