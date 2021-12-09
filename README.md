
# GeneDetector
## Discription
Get the gene name that related to your interested disease from PubMed Abstract.
## package Used 
[GoFrame](https://github.com/gogf/gf) is an application development framework of Golang.\
[html](https://pkg.go.dev/golang.org/x/net/html) implements an HTML5-compliant tokenizer and parser. \
[htmlquery](https://github.com/antchfx/htmlquery) supports HTML document query.\
[XPath](https://github.com/antchfx/xpath) is Go package provides selecting nodes from HTML or other documents using XPath expression.\
[regexp2](https://github.com/dlclark/regexp2) is a regex engine in pure Go based on the .NET engine\
[csv](https://pkg.go.dev/encoding/csv) is a package related to read and write csv file\

## How to run the project
Type in disease name and abstract number, and you will get disease related gene name and other paper information in a csv file.
Default disease name is Alzheimer's and default abstract number is 10.
```go
$ go build
$ ./GeneDetector
```
OR
```go
$ go build
$ ./GeneDetector -disease diabetes -n 20
```



## Expected output
A CSV file with the related gene symbol.
There are other information in this csv file, including: paper title, url, abstract content, gene name, pmid, doi,
keyword(disease name)

## Changes
The main strategy to get the related gene name is use **regular expression** to match the gene symbol. 
I didn't use text mining strategies to get the gene name due to my limit knowledge background. 
But I do spend time to get to know more about the knowledge of text mining. 
I found that Named-entity recognition(NER) would be a good strategy for the next step of this project.

## Acknowledge
Thanks for **Robin L.** having good discussions about the tips of web scrap!\
Thanks for Professor **Carl Kingsford** importing all the great knowledge about golang!\
Thanks for TA **Siddharth Reed** for giving me lots of help when I found this project was really hard to go on!\
Thanks for TA **Jingjing Tang** for grading and giving feedback!


