package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"math"
	"os"
	"strings"
	"time"
)

// GetGeneResult takes a disease keyword and the page limit, and gets all the information in cvs file.
func GetGeneResult(keyword string, limit int)  {
	// get work dictionary.
	rootFolderPath, _ := os.Getwd()
	g.Log().Infof("folder %s", rootFolderPath)
	resultFolderPath := fmt.Sprintf("%s/csv", rootFolderPath)
	// If the file don't exist, create a file
	_, err := os.Stat(resultFolderPath)
	if os.IsNotExist(err) {
		_ = os.Mkdir(resultFolderPath, os.ModePerm)
	}
	timeObj := time.Now()
	now := fmt.Sprintf("%d%d%d_%d%d%d", timeObj.Year(), timeObj.Month(), timeObj.Day(), timeObj.Hour(), timeObj.Minute(), timeObj.Second())
	csvFilePath := fmt.Sprintf("%s/%s_%s.csv", resultFolderPath, keyword, now)
	csvFile, err := os.OpenFile(csvFilePath, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	defer fmt.Printf("CSV file has been saved to path: [ %s ]\n", csvFilePath)

	fmt.Printf("CSV file will be saved to path: [ %s ]\n", csvFilePath)

	csvWriter := csv.NewWriter(csvFile)
	_ = csvWriter.Write([]string{"title", "url", "abstract", "gene", "pmid", "doi", "keyword"})

	g.Log().Info("Downloading page 1...")
	firstPageUrl, csrfToken, cookie, currentPage, totalPageCount := DownloadFirstSearchPage(keyword, csvWriter, limit)
	// g.Log().Info(currentPage)
	if currentPage < totalPageCount {
		for {
			currentPage = currentPage + 1
			if currentPage < totalPageCount {
				fmt.Printf("Downloading page %d...\n", currentPage)
				hasNext := DownloadFollowingSearchPage(keyword, firstPageUrl, csrfToken, cookie, currentPage, csvWriter, limit)
				if !hasNext {
					break
				}
			}
		}
	}
}

func main() {
	g.Log().Info("Please input name of the disease you want to search: ")
	inputReader := bufio.NewReader(os.Stdin)
	keyword, _ := inputReader.ReadString('\n')
	g.Log().Info(keyword)
	g.Log().Info("Please input count of papers you want to download: ")
	var limit int
	_, _ = fmt.Scanln(&limit)
	g.Log().Infof("Downloading... Please wait... Expect to download %d pages...\n", int(math.Ceil(float64(limit) / 10)))  // 10 abstracts a page.
	GetGeneResult(strings.TrimSpace(keyword), limit)
	g.Log().Info("Task finished!")
}
