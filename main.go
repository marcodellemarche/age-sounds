package main

import (
	"log"
	"mime"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	crawl()
}

func crawl() {
	_ = os.Mkdir("./download", os.ModePerm)
	
	mainCollector := colly.NewCollector(
		colly.AllowedDomains("ageofempires.fandom.com", "www.ageofempires.fandom.com"),
		colly.CacheDir("./cache"),
	)

	soundsCollector := colly.NewCollector(
		colly.AllowedDomains("static.wikia.nocookie.net", "www.static.wikia.nocookie.net"),
		colly.CacheDir("./cache"),
	)

	mainCollector.OnRequest(func(r *colly.Request) {
		log.Printf("Scraping URL: %s", r.URL.String())
	})

	soundsCollector.OnResponse(func(r *colly.Response) {
		log.Printf("Downloading sound: %s", r.Request.URL.String())

		contentType := r.Headers.Get("Content-Type")

		if contentType != "application/ogg" {
			return
		}

		contentDisposition := r.Headers.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)

		if err != nil {
			log.Printf("Error Content-Disposition header: %v", err)
			return
		}

		filename := params["filename"]

		path := "./download/" + params["filename"]
		err = r.Save(path)
		
		if err != nil {
			log.Printf("Error downloading sound: %v", err)
			return
		}

		log.Printf("Downloaded sound %s into %s", filename, path)
	})

	mainCollector.OnHTML("h3 > span > a[href]", func(e *colly.HTMLElement) {
		if e.Attr("title") == "" {
			return
		}

		link := e.Attr("href")

		if !strings.HasPrefix(link, "/wiki/") {
			return
		}

		absoluteLink := e.Request.AbsoluteURL(link)
		mainCollector.Visit(absoluteLink)
	})

	mainCollector.OnHTML("audio", func(e *colly.HTMLElement) {
		source := e.Attr("src")
		soundsCollector.Visit(source)
	})

	mainCollector.Visit("https://ageofempires.fandom.com/wiki/Civilizations_(Age_of_Empires_II)")

}
