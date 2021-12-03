package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"net/url"
	"strings"
	"time"
)

var ExistPaperCount2 = 0

// DownloadFirstSearchPage2 downloads the first search page in file.
func DownloadFirstSearchPage2(keyword string, csvWriter *csv.Writer, limit int) (string, string, string, int, int) {
	// targetUrl is the target URL including the keyword.
	targetUrl := fmt.Sprintf("https://www.genecards.org/Search/Keyword?queryString=%s", url.QueryEscape(keyword))

	client := g.Client() //client is an HTTP client.
	client.SetHeaderRaw(`
		accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
		accept-language: en-US,en;q=0.9
		cache-control: no-cache
		pragma: no-cache
		upgrade-insecure-requests: 1
		user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36
	`)

	response, err := client.Get(targetUrl) // The response contains status information about the request.
	if err != nil {
		panic(err)
	}
	body := response.ReadAllString()  // body is the <body> element containing all the contents of an HTML document.

	currentPage := 1
	// csrfTokens are used to send requests from an authenticated user to a web application.
	// csrfTokens can be gotten from html document.
	csrfToken, totalPageCount := ParseFirstPage2(body)
	if totalPageCount > 1000 {  // The maximum value of totalPageCount is 1000 (limited by the website)
		totalPageCount = 1000
	}

	var cookieList []string
	dataList := make([][]string, 0)
	for key, value := range response.GetCookieMap() {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", key, value))
	}
	cookie := strings.Join(cookieList, "; ")

	// range all the url of single page of abstract and get all the data from each abstract.
	for i, geneCardUrl := range ParseGeneCardUrlList(body) {
		g.Log().Infof("Begin to write data line: %d", i)
		geneCardDetailBody := DownloadGeneCardDetail(geneCardUrl, targetUrl)
		geneCard := CreateGeneCard()
		geneCard.ParseGeneCard(geneCardUrl, geneCardDetailBody,targetUrl)
		//paper := CreatePaper()
		//paper.ParsePaper(paperUrl, paperDetailBody, keyword)
		//dataList = append(dataList, []string{paper.title, paper.url, paper.abstract, paper.gene, paper.pmid, paper.doi, paper.keyword})
		ExistPaperCount += 1
		if ExistPaperCount >= limit {
			currentPage = OverflowPage
			break
		}
	}

	// write data in cvs file.
	_ = csvWriter.WriteAll(dataList)
	return targetUrl, csrfToken, cookie, currentPage, totalPageCount
}


// DownloadFollowingSearchPage2 downloads the following pages after the first page.
func DownloadFollowingSearchPage2(keyword string, referer string, csrfToken string, cookie string, page int, csvWriter *csv.Writer, limit int) bool {  // function to get content of the following search result page
	// the target urls of the following pages are different from the target url of the first page.
	targetUrl := "https://www.genecards.org/more/"
	client := g.Client()
	client.SetHeaderRaw(fmt.Sprintf(` 
		accept: */*
		accept-language: en-US,en;q=0.9
		cache-control: no-cache
		content-type: application/x-www-form-urlencoded; charset=UTF-8
		cookie: %s
		origin: https://pubmed.ncbi.nlm.nih.gov
		pragma: no-cache
		referer: %s
		user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36
		x-requested-with: XMLHttpRequest
	`, cookie, referer))
	data := g.Map{
		"term": keyword,
		"filter": "simsearch1.fha",
		"no_cache": "yes",
		"page": page,
		"no-cache": time.Now().UnixMilli(),
		"csrfmiddlewaretoken": csrfToken,
	}
	response, err := client.Post(targetUrl, data)
	if err != nil {
		panic(err)
	}
	body := response.ReadAllString()

	// Referer is the name of an optional HTTP header field that identifies the address of the web page
	// which is linked to the resource being requested.
	referer = fmt.Sprintf("https://www.genecards.org/Search/Keyword?queryString=%s&startPage=%d", url.QueryEscape(keyword), page)

	// range all the url of single page of abstract and get all the data from each abstract.
	// then write them into file.
	for i, paperUrl := range ParsePaperUrlList(body) {
		g.Log().Infof("Begin to write data line: %d", i)
		paperDetailBody := DownloadPaperDetail(paperUrl, referer)
		paper := CreatePaper()
		paper.ParsePaper(paperUrl, paperDetailBody, keyword)
		_ = csvWriter.Write([]string{paper.title, paper.url, paper.abstract, paper.gene, paper.pmid, paper.doi, paper.keyword})
		ExistPaperCount += 1
		if ExistPaperCount >= limit {
			return false
		}
	}

	return true
}

// DownloadGeneCardDetail takes the targetUrl and referer and returns the body.
func DownloadGeneCardDetail(targetUrl string, referer string) string {
	client := g.Client()
	client.SetHeaderRaw(fmt.Sprintf(`
		accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
		accept-language: en-US,en;q=0.9
		cache-control: no-cache
		pragma: no-cache
		referer: %s
		upgrade-insecure-requests: 1
		user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36
	`, referer))

	response, err := client.Get(targetUrl)
	if err != nil {
		panic(err)
	}

	body := response.ReadAllString()
	return body
}
