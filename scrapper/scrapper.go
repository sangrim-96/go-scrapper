package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	company  string
	location string
	date     string
	url      string
}

// Strip string
func Strip(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// Check Error when http.Get
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check response status code
func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}

func getTotalPages(url string) int {
	totalPages := 0
	// Request the HTML page.
	res, err := http.Get(url)
	checkErr(err)
	checkStatusCode(res)

	// Prevent memory leak
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection) {
		totalPages = s.Find("li").Length()
	})

	return totalPages
}

func getPage(page int, baseURL string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := fmt.Sprintf("%s&start=%d", baseURL, page*10)
	fmt.Println("Requesting... ", page)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	cardList := doc.Find(".jobsearch-SerpJobCard")
	cardList.Each(func(i int, card *goquery.Selection) {
		go getCardInfo(card, pageURL, c)
	})
	for i := 0; i < cardList.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func getCardInfo(card *goquery.Selection, pageURL string, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	url := fmt.Sprintf("%s&vjk=%s", pageURL, id)
	title := Strip(card.Find(".title>a").Text())
	company := Strip(card.Find(".company").Text())
	location := Strip(card.Find(".location").Text())
	date := Strip(card.Find(".date").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		date:     date,
		url:      url,
	}
}

func saveJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Title", "Location", "date", "URL"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.title, job.location, job.date, job.url}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

// Scrape Indeed by query
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=30"
	var jobs []extractedJob

	totalPages := getTotalPages(baseURL)
	c := make(chan []extractedJob)
	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	saveJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}
