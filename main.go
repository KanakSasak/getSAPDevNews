package main

import (
	"bytes"
	"fmt"
	"getSAPDevNews/model"
	"getSAPDevNews/slackbuilder"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var listnews []model.News
	news := model.News{}

	webPage := "https://blogs.sap.com/tags/8077228b-f0b1-4176-ad1b-61a78d61a847/"
	resp, err := http.Get(webPage)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	f := func(i int, s *goquery.Selection) bool {

		link, _ := s.Attr("href")
		return strings.HasPrefix(link, "https")
	}

	doc.Find("ul.dm-contentList li.dm-contentListItem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".dm-user__name").Text()
		date := s.Find(".dm-user__date").Text()
		datearr := strings.Split(date, "\n")
		datestr := strings.TrimSpace(datearr[2])
		category := s.Find(".dm-user__category").Text()
		desc := s.Find(".dm-content-list-item__text").Text()
		title := s.Find(".dm-contentListItem__title").Text()

		news.Title = title
		news.Author = name
		news.Date = datestr
		news.Category = category
		news.Desc = desc

		s.Find(".dm-contentListItem__title a").FilterFunction(f).Each(func(_ int, tag *goquery.Selection) {

			link, _ := tag.Attr("href")
			//linkText := tag.Text()
			//fmt.Printf("%s\n", link)
			news.Link = link
		})

		listnews = append(listnews, news)
	})

	fmt.Println("-----------------------------------")

	fmt.Println(listnews[0].Link)
	fmt.Println(listnews[0].Date)
	fmt.Println(listnews[0].Author)
	fmt.Println(listnews[0].Title)
	fmt.Println(listnews[0].Category)
	fmt.Println(listnews[0].Desc)

	payload := slackbuilder.Build(listnews)

	httpposturl := os.Getenv("SLACK_WEBHOOK")
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(payload)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

}
