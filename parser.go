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

// Paper is a struct that includes all the information that we want according to the input keyword.
type Paper struct {
	title    string
	url      string
	abstract string
	gene     string
	pmid     string
	doi      string
	keyword  string
}

// ParseFirstPage parses the fist search result page.
// It takes the body and returns csrfToken and totalPageCount.
func ParseFirstPage(body string) (string, int) {
	//Load HTML document from string.
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	expr, _ := xpath.Compile("string(//input[@name='csrfmiddlewaretoken']/@value)")
	navigator := htmlquery.CreateXPathNavigator(doc)
	csrfToken := expr.Evaluate(navigator).(string)

	expr, _ = xpath.Compile("string(//span[@class='total-pages'])")
	totalPageString := strings.TrimSpace(expr.Evaluate(navigator).(string))
	totalPageCount, _ := strconv.Atoi(strings.Replace(totalPageString, ",", "", -1))

	return csrfToken, totalPageCount
}

// ParsePaperUrlList takes a string of html and returns a list of url of every single abstract.
func ParsePaperUrlList(body string) []string {
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

// GetTitle is a methode getting the title of the paper.
func (paper *Paper)GetTitle(doc *html.Node) {
	titleNode, _ := htmlquery.Query(doc, "//h1[@class='heading-title']")
	paper.title = strings.TrimSpace(htmlquery.InnerText(titleNode))
}

// GetAbstract is a methode getting the content of the abstract of the paper.
func(paper *Paper)GetAbstract(doc *html.Node)  {
	abstractNode, _ := htmlquery.Query(doc, "//div[@id='enc-abstract']")
	pattern, _ := regexp.Compile(`\n\s*`)
	abstract := pattern.ReplaceAllString(htmlquery.InnerText(abstractNode), "\n")
	paper.abstract = strings.TrimSpace(abstract)
}

// GetGeneName is a methode getting the gene name from the abstract.
func (paper *Paper)GetGeneName() {
	//pattern, _ := regexp.Compile(`[A-Z][\w-]{0,5}\b`) // ****************************** here need to be revised!!!
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

// GetPmid is a methode getting pmid of the paper.
func (paper *Paper)GetPmid(doc *html.Node) {
	pmidNode, _ := htmlquery.Query(doc, "//strong[@title='PubMed ID']")
	if pmidNode == nil {
		paper.pmid = ""
	} else {
		paper.pmid = strings.TrimSpace(htmlquery.InnerText(pmidNode))
	}
}

// GetDoi is a methode getting the doi of the paper.
func (paper *Paper)GetDoi(doc *html.Node){
	doiNode, _ := htmlquery.Query(doc, "//a[@data-ga-action='DOI']")
	if doiNode == nil {
		paper.doi = ""
	} else {
		paper.doi = strings.TrimSpace(htmlquery.InnerText(doiNode))
	}
}
