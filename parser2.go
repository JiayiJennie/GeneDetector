package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

type GeneCard struct {
	url    string
	symbol string
	keyword string
}

// ParseFirstPage2 parses the fist search result page.
// It takes the body and returns csrfToken and totalPageCount.
func ParseFirstPage2(body string) (string, int) {
	//Load HTML documents from string.
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	expr, _ := xpath.Compile("string(//input[@name='csrfmiddlewaretoken']/@value)") //***************************
	navigator := htmlquery.CreateXPathNavigator(doc)
	csrfToken := expr.Evaluate(navigator).(string)

	expr, _ = xpath.Compile("string(//span[@class='total-pages'])") //*********************88
	totalPageString := strings.TrimSpace(expr.Evaluate(navigator).(string))
	totalPageCount, _ := strconv.Atoi(strings.Replace(totalPageString, ",", "", -1))

	return csrfToken, totalPageCount
}

// ParseGeneCardUrlList takes a string of html and returns a list of url of every single abstract.
func ParseGeneCardUrlList(body string) []string {
	var result []string

	doc, err1 := htmlquery.Parse(strings.NewReader(body))
	if err1 != nil {
		panic(err1)
	}

	// Find href value that lead to each single page of abstract.
	// put them into a list of href value.
	list, err2 := htmlquery.QueryAll(doc, "//div[@class='search-results-chunk results-chunk']/article[@class='full-docsum']/div[@class='docsum-wrap']/div[@class='docsum-content']")
	if err2 != nil {
		fmt.Println("Error: cannot find matched html nodes.")
	}

	// get paperUrl from href value.
	for _, node := range list {
		titleNode, _ := htmlquery.Query(node, "./a[@class='docsum-title']") //*****************************
		paperUrl := fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov%s", htmlquery.SelectAttr(titleNode, "href"))
		result = append(result, paperUrl)
	}

	return result
}

// CreateGeneCard makes a Paper.
func CreateGeneCard() *GeneCard{
	geneCard := new(GeneCard)
	return geneCard
}

// ParseGeneCard takes a paperUrl, a body and the input keyword, and returns a pointer of a Paper.
func (geneCard *GeneCard)ParseGeneCard(paperUrl string, body string, keyword string)  {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	// Get all the information we need from the paper.
	geneCard.url = paperUrl
	geneCard.GetSymbol(doc)
	geneCard.keyword = keyword
}

// GetSymbol is a methode getting the title of the paper.
func (geneCard *GeneCard)GetSymbol(doc *html.Node) {
	symbolNode, _ := htmlquery.Query(doc, "//h1[@class='heading-title']") // *************************
	geneCard.symbol = strings.TrimSpace(htmlquery.InnerText(symbolNode))
}





