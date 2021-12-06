// Thanks for Robin L. having good discussions about the tips of web scrap!
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// GetGeneResult takes a keyword and a paper count, and output all the information in csv file.
func GetGeneResult(keyword string, paperNumber int)  {
	// create a csv file to hold output

	// Use the time right now to make a unique name for the output csv file.
	timeObj := time.Now()
	now := fmt.Sprintf("%d%d%d_%d%d%d", timeObj.Year(), timeObj.Month(), timeObj.Day(), timeObj.Hour(), timeObj.Minute(), timeObj.Second())
	fileName := fmt.Sprintf("%s_%s.csv", keyword, now)
	csvFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666) // O_CREATE creates a csv file. perm 0666: everyone can read and write this file.
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	defer fmt.Println("CSV file has been saved")

	// Write headers in csv files
	csvWriter := csv.NewWriter(csvFile)
	err2 := csvWriter.Write([]string{"title", "url", "abstract", "gene", "pmid", "doi", "keyword"})
	if err2 != nil{
		panic(err2)
	}

	// download first search page
	fmt.Println("Downloading page 1...")
	firstPageUrl, csrfToken, cookie, totalPageCount := DownloadFirstSearchPage(keyword, csvWriter, paperNumber)

	// download following search pages
	for currentPage := 2; currentPage < totalPageCount; currentPage++ {
		fmt.Printf("Downloading page %d...\n", currentPage)
		hasNext := DownloadFollowingSearchPage(keyword, firstPageUrl, csrfToken, cookie, currentPage, csvWriter, paperNumber)
		if !hasNext {
			break
		}
	}
}


func main() {
	var keyword string
	var paperNumber int

	flag.StringVar(&keyword, "disease", "Alzheimer's", "disease name")
	flag.IntVar(&paperNumber, "n", 10, "paper number")
	flag.Parse()

	GetGeneResult(strings.TrimSpace(keyword), paperNumber)

	fmt.Println("Task finished!")

	// Example Command:
	// ./GeneDetector
	// ./GeneDetector -disease diabetes 20
}

