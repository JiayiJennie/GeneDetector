package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/dlclark/regexp2"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"strings"
)

// Paper is a type that includes all the information that we want according to the input keyword.
type Paper struct {
	title    string
	url      string
	abstract string
	gene     string
	pmid     string
	doi      string
	keyword  string
}

// ParseFirstPage parses the body of fist search result page.
// It takes the body and returns csrfToken and totalPageCount.
func ParseFirstPage(body string) (string, int) {
	//Load HTML documents from string.
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	// get csrfToken
	expr, _ := xpath.Compile("string(//input[@name='csrfmiddlewaretoken']/@value)")
	navigator := htmlquery.CreateXPathNavigator(doc)
	csrfToken := expr.Evaluate(navigator).(string)

	// Get total page count from the search result
	expr2, _ := xpath.Compile("string(//span[@class='total-pages'])")
	totalPageString := strings.TrimSpace(expr2.Evaluate(navigator).(string))
	totalPageCount, _ := strconv.Atoi(strings.Replace(totalPageString, ",", "", -1))

	return csrfToken, totalPageCount
}

// ParsePaperUrlList takes a string of html and returns a list of url of every single abstract.
func ParsePaperUrlList(body string) []string {
	var result []string

	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	// put html nodes of abstract search result into a list.
	// On pubmed, one search page has 10 abstract results.
	// so the list has 10 html nodes.
	list, err2 := htmlquery.QueryAll(doc, "//div[@class='search-results-chunk results-chunk']/article[@class='full-docsum']/div[@class='docsum-wrap']/div[@class='docsum-content']")
	if err2 != nil {
		panic(err2)
	}

	// get paperUrl from href values.
	// href value can lead to a new single page of abstract.
	for _, node := range list {
		titleNode, _ := htmlquery.Query(node, "./a[@class='docsum-title']")
		paperUrl := fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov%s", htmlquery.SelectAttr(titleNode, "href"))
		result = append(result, paperUrl)
	}

	return result
}

// CreatePaper makes a Paper.
func CreatePaper() *Paper{
	paper := new(Paper)
	return paper
}

// ParsePaper takes a paperUrl, a body and the input keyword, and returns a pointer of a Paper.
func (paper *Paper)ParsePaper(paperUrl string, body string, keyword string)  {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	// Get all the information we need from the paper.
	paper.GetTitle(doc)
	paper.url = paperUrl
	paper.GetAbstract(doc)
	paper.GetGeneName()
	paper.GetPmid(doc)
	paper.GetDoi(doc)
	paper.keyword = keyword
}

// GetTitle is a method getting the title of the paper.
func (paper *Paper)GetTitle(doc *html.Node){
	titleNode, _ := htmlquery.Query(doc, "//h1[@class='heading-title']")
	paper.title = strings.TrimSpace(htmlquery.InnerText(titleNode))
}

// GetAbstract is a method getting the content of the abstract of the paper.
func(paper *Paper)GetAbstract(doc *html.Node)  {
	abstractNode, _ := htmlquery.Query(doc, "//div[@id='enc-abstract']")
	pattern, _ := regexp.Compile(`\n\s*`)
	abstract := pattern.ReplaceAllString(htmlquery.InnerText(abstractNode), "\n")
	paper.abstract = strings.TrimSpace(abstract)
}

// GetGeneName is a method getting the gene name from the abstract.
// Here use pattern to match the gene name.
func (paper *Paper)GetGeneName() {
	genePattern := regexp2.MustCompile(`[A-Z][A-Z\d-]{1,5}(?<![-])\b`, 0)

	tempMap := map[string]string{}
	m, _ := genePattern.FindStringMatch(paper.abstract)
	for m != nil {
		tempMap[m.String()] = ""
		m, _ = genePattern.FindNextMatch(m)
	}

	var geneList []string
	for key, _ := range tempMap {
		geneList = append(geneList, key)
	}
	paper.gene = strings.Join(geneList, "|")
}

// GetPmid is a method getting pmid of the paper.
func (paper *Paper)GetPmid(doc *html.Node) {
	pmidNode, _ := htmlquery.Query(doc, "//strong[@title='PubMed ID']")
	if pmidNode == nil {
		paper.pmid = ""
	} else {
		paper.pmid = strings.TrimSpace(htmlquery.InnerText(pmidNode))
	}
}

// GetDoi is a method getting the doi of the paper.
func (paper *Paper)GetDoi(doc *html.Node){
	doiNode, _ := htmlquery.Query(doc, "//a[@data-ga-action='DOI']")
	if doiNode == nil {
		paper.doi = ""
	} else {
		paper.doi = strings.TrimSpace(htmlquery.InnerText(doiNode))
	}
}
