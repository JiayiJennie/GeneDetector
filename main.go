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

//// DataList is a [][]string that holds the data from pubmed abstract
//type DataList [][]string
//
////CreateDataList make a DataList and return a pointer of it.
//func CreateDataList() *DataList{
//	dataList := make(DataList,0)
//	return &dataList
//}

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


	// make a datalist to hold all the data
	dataList := make([][]string,0)

	// Write headers in csv files
	csvWriter := csv.NewWriter(csvFile)
	err2 := csvWriter.Write([]string{"title", "url", "abstract", "gene", "pmid", "doi", "keyword"})
	if err2 != nil{
		panic(err2)
	}

	// download first search page
	fmt.Println("Downloading page 1...")
	firstPageUrl, csrfToken, cookie, totalPageCount := DownloadFirstSearchPage(keyword, &dataList, paperNumber)

	// download following search pages
	for currentPage := 2; currentPage <= totalPageCount; currentPage++ {
		fmt.Printf("Downloading page %d...\n", currentPage)
		hasNext := DownloadFollowingSearchPage(keyword, firstPageUrl, csrfToken, cookie, currentPage, &dataList, paperNumber)
		if !hasNext {
			break
		}
	}

	//write dataList to csv file

	err3 := csvWriter.WriteAll(dataList[0:paperNumber])
	if err3 != nil {
		panic(err3)
	}
}


func main() {
	var keyword string
	var paperNumber int

	flag.StringVar(&keyword, "disease", "Alzheimer's", "disease name")
	flag.IntVar(&paperNumber, "n", 10, "paper number")
	flag.Parse()

	GetGeneResult(strings.TrimSpace(keyword), paperNumber)

	fmt.Println("GeneDetector finished!")

	// Example Command:
	// ./GeneDetector
	// ./GeneDetector -disease diabetes -n 20
}

