package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/url"
	"strings"
	"time"
)

var ExistPaperCount = 0

// DownloadFirstSearchPage downloads the first search page in file.
func DownloadFirstSearchPage(keyword string, csvWriter *csv.Writer, paperNumber int) (string, string, string, int) {
	// targetUrl is the target URL including the keyword.
	targetUrl := fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov/?term=%s&filter=simsearch1.fha", url.QueryEscape(keyword))

	// Request and response
	response := RequestAndResponseForFirstPage(targetUrl)

	// get body
	body := response.ReadAllString()  // body is the <body> element containing all the contents of an HTML document.

	// Capture the CSRF token value and get totalSearchPageCount from html body.
	csrfToken, totalSearchPageCount := GetCsrfToken(body)

	// Get cookie from response
	cookie := GetCookie(response)

	// parse body to get urls of every paper in the first search page
	// range all the urls
	// and get the data from each abstract.
	dataList := make([][]string, 0)
	for i, paperUrl := range ParsePaperUrlList(body) {
		fmt.Printf("Begin to write data line: %d\n", i)
		paperDetailBody := DownloadPaperDetail(paperUrl, targetUrl)
		paper := CreatePaper()
		paper.ParsePaper(paperUrl, paperDetailBody, keyword)
		dataList = append(dataList, []string{paper.title, paper.url, paper.abstract, paper.gene, paper.pmid, paper.doi, paper.keyword})
		ExistPaperCount ++
		if ExistPaperCount >= paperNumber{
			break
		}
	}

	// write data in cvs file.
	err := csvWriter.WriteAll(dataList)
	if err != nil{
		panic(err)
	}

	return targetUrl, csrfToken, cookie, totalSearchPageCount
}

// DownloadFollowingSearchPage downloads the following pages after the first page.
func DownloadFollowingSearchPage(keyword string, referer string, csrfToken string, cookie string, currPage int, csvWriter *csv.Writer, paperNumber int) bool {
	targetUrl := "https://pubmed.ncbi.nlm.nih.gov/more/" // the target urls of the following pages are different from the target url of the first page.

	// request and response
	response := RequestAndResponseForFollowingSearchPages(targetUrl, cookie, keyword, csrfToken, currPage)

	//get body
	body := response.ReadAllString()

	// Referer is the name of an optional HTTP header field that identifies the address of the web page
	// which is linked to the resource being requested.
	referer = fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov/?term=%s&filter=simsearch1.fha&page=%d", url.QueryEscape(keyword), currPage)

	// parse body to get urls of every paper in a single search page
	// range all the urls
	// and get the data from each abstract.
	// then write them into file.
	for i, paperUrl := range ParsePaperUrlList(body) {
		paperDetailBody := DownloadPaperDetail(paperUrl, referer)
		paper := CreatePaper()
		paper.ParsePaper(paperUrl, paperDetailBody, keyword)
		fmt.Printf("Begin to write data line: %d\n", i)
		err := csvWriter.Write([]string{paper.title, paper.url, paper.abstract, paper.gene, paper.pmid, paper.doi, paper.keyword})
		if err != nil{
			panic(err)
		}
		ExistPaperCount ++
		if ExistPaperCount >= paperNumber {
			return false
		}
	}
	return true
}

// RequestAndResponseForFirstPage takes targetUrl to request and returns a response.
func RequestAndResponseForFirstPage(targetUrl string) *ghttp.ClientResponse {

	// create a new client
	client := g.Client()

	// Set request headers
	client.SetHeaderRaw(fmt.Sprintf(`
		accept: */*
		accept-language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7
		user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36
	`))

	// send GET request and returns the response object.
	response, err := client.Get(targetUrl) // The response contains body in html.
	if err != nil {
		panic(err)
	}

	return response
}

// RequestAndResponseForFollowingSearchPages takes targetUrl, cookie, keyword, csrfToken and currPage and returns a response.
// cookie and csrfToken are got from parse fist search page.
func RequestAndResponseForFollowingSearchPages(targetUrl, cookie, keyword, csrfToken string, currPage int)*ghttp.ClientResponse{
	client := g.Client()
	client.SetHeaderRaw(fmt.Sprintf(` 
		accept: */*
		accept-language: en-US,en;q=0.9
		cache-control: no-store, no-cache, must-revalidate, max-age=0
		content-type: application/x-www-form-urlencoded
		cookie: %s
		origin: https://pubmed.ncbi.nlm.nih.gov
		pragma: no-cache
		referer: https://pubmed.ncbi.nlm.nih.gov/
		user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36
		x-requested-with: XMLHttpRequest
	`, cookie))

	// data are parameters used for post request
	data := g.Map{
		"term": keyword,
		"filter": "simsearch1.fha",
		"no_cache": "yes",
		"page": currPage,
		"no-cache": time.Now().UnixMilli(),
		"csrfmiddlewaretoken": csrfToken,
	}

	// post request
	// Note that here use post request, not get request as the first search page.
	response, err := client.Post(targetUrl, data)
	if err != nil {
		panic(err)
	}

	return response
}

// GetCsrfToken takes a body and returns csrfToken and totalSearchPageCount
// parse first search page to get csrfToken, totalPageCount
// csrfToken is a random key for safe access.
// csrfTokens are used to send requests from an authenticated user to a web application.
// totalPageCount is the number of total search page.
func GetCsrfToken(body string) (string, int){
	csrfToken, totalSearchPageCount := ParseFirstPage(body)
	if totalSearchPageCount > 1000 {  // The maximum value of totalPageCount is 1000 (limit set by myself)
		totalSearchPageCount = 1000
	}
	return csrfToken, totalSearchPageCount
}

// GetCookie takes a response and returns a cookie.
// A cookie is often used to identify a user.
// A cookie is a small file that the server embeds on the user's computer.
// Each time the same computer requests a page with a browser, it will send the cookie too.
func GetCookie(response *ghttp.ClientResponse) string {
	var cookieList []string
	for key, value := range response.GetCookieMap() {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", key, value))
	}
	cookie := strings.Join(cookieList, "; ")

	return cookie
}

// DownloadPaperDetail takes the targetUrl and referer and returns the body.
func DownloadPaperDetail(targetUrl string, referer string) string {
	client := g.Client()
	client.SetHeaderRaw(fmt.Sprintf(`
		accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
		accept-language: en-US,en;q=0.9
		cache-control: no-cache
		pragma: no-cache
		referer: %s
		upgrade-insecure-requests: 1
		user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36
	`, referer))

	response, err := client.Get(targetUrl)
	if err != nil {
		panic(err)
	}

	body := response.ReadAllString()
	return body
}
