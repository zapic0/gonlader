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
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	cliArgs := os.Args[1:]

	fullURL := cliArgs[0]

	start := time.Now()
	parseWeb(fullURL, cliArgs)
	end := time.Now()

	elapsed := end.Sub(start)
	fmt.Println("Time to parse full website:" + elapsed.String() + "\n")
}

func parseWeb(fullURL string, cliArgs []string) {
	pageParam := cliArgs[1]

	firstPage, err := strconv.Atoi(cliArgs[2])
	lastPage, err := strconv.Atoi(cliArgs[3])

	chPage := make(chan string)

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

		go parsePage(page, cliArgs, chPage)
	}

	for j := firstPage; j < lastPage; j++ {
		fmt.Println(<-chPage)
	}

}

func parsePage(page *goquery.Document, cliArgs []string, chPage chan<- string) {
	start := time.Now()

	host := cliArgs[4]
	classToSearch := cliArgs[5]
	fileClass := cliArgs[6]
	fileAttr := cliArgs[7]

	ch := make(chan string)
	n := 0

	page.Find(classToSearch).Each(func(i int, s *goquery.Selection) {
		playValue, _ := s.Find(fileClass).First().Attr(fileAttr)
		file := strings.Split(playValue, "'")[1]
		fileURL := host
		if file != "" {
			fileURL = fileURL + file
		}

		name := s.Find("a").Text()
		fmt.Println(fileURL)

		n++
		go downloadFile(fileURL, name, cliArgs, ch)
	})

	for i := 1; i < n; i++ {
		fmt.Println(<-ch)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	finishedMessage := "Time to parse page: " + elapsed.String() + "\n"

	chPage <- finishedMessage
}

func downloadFile(url string, name string, cliArgs []string, ch chan<- string) {
	start := time.Now()

	basicFolder := cliArgs[8]
	err := os.Mkdir(filepath.Join(basicFolder, name), os.ModePerm)
	message := ""

	remoteFile, err := http.Get(url)
	if err != nil {
		log.Println(err)
		message = "----- Error getting page " + url
		fmt.Println("----- Error getting page " + url)
	}

	fileChunks := strings.Split(url, "/")
	fileName := fileChunks[len(fileChunks)-1]

	defer remoteFile.Body.Close()

	fileFullPath := filepath.Join(basicFolder, filepath.Join(name, fileName))

	message = fileFullPath + "\n"

	out, err := os.Create(fileFullPath)

	if err != nil {
		log.Println(err)
		message = "----- Error creating / downloading file " + fileFullPath
	}
	defer out.Close()

	io.Copy(out, remoteFile.Body)

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Time to download file " + fileFullPath + " :" + elapsed.String())

	ch <- message
}
