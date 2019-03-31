package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/WarpRat/scrape/aws"
	"github.com/aws/aws-lambda-go/lambda"

	. "github.com/WarpRat/scrape/config"

	"github.com/PuerkitoBio/goquery"
)

const table string = "frenchy"

func scrapeFrency() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	request, err := http.NewRequest("GET", "http://rockaway.seatbytext.com/mobilewait/", nil)
	if err != nil {
		log.Panic("Error building request", err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 9; Pixel 2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.90 Mobile Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		log.Panic("Error around the http client", err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Panic("Error loading http body", err)
	}

	var parties []Res

	document.Find("li").Filter(".menu").Each(func(i int, s *goquery.Selection) {

		var newres Res
		s.Find("td").Each(func(ii int, se *goquery.Selection) {

			if ii == 1 {
				re := regexp.MustCompile(`READY![\r\n]`)
				name := re.ReplaceAllLiteralString(se.Text(), "")
				newres.Name = name
			} else if ii == 3 {
				newres.Party = se.Text()
			}
		})

		parties = append(parties, newres)

	})

	//fmt.Printf("%+v\n\n", parties)

	aws.LoadDynamo(parties, table)
}

//HandleRequest is a lambda requirement
func HandleRequest() {
	scrapeFrency()
}

func main() {
	lambda.Start(HandleRequest)
}
