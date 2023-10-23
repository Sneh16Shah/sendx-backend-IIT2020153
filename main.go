package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pq "github.com/Sneh16Shah/sendx-backend-IIT2020153/priority_queue"
	"github.com/gocolly/colly"
)

type PageData struct {
	CrawledURLs []string
}

type CachedPage struct {
	HTML        string
	CrawledURLs []string
}

var requestQueue = pq.NewPriorityQueue()
var wg sync.WaitGroup

type CrawlRequest struct {
	URL          string
	CustomerType string
	Priority     int
	UniqueURLs   []string 
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := PageData{
			CrawledURLs: []string{},
		}
		tmpl.Execute(w, data)
	})

	cleanupCache()
	go crawlCustomers()
	http.HandleFunc("/crawl", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		urlStr := r.FormValue("path")
		customerType := r.FormValue("customer_type")
		priority := 0 // Non-paying customer
		if customerType == "paying" {
			priority = 1 // Paying customer
		}

		request := &CrawlRequest{
			URL:          urlStr,
			CustomerType: customerType,
		}

		// Push the request into the priority queue
		obj := &pq.Element{
			Value: request,
			Priority: priority,
		}
		requestQueue.Push(obj)

		// // Inside this goroutine, we crawl the URL and then render the results.
		wg.Wait()
		crawlURL(request)
		// After crawling is done, render the results.
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := PageData{
			CrawledURLs: request.UniqueURLs, // Use unique URLs from the request
		}
		for _,link := range request.UniqueURLs{
			fmt.Println(link)
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			// Handle the error, e.g., log it
			fmt.Println("Template rendering error:", err)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func crawlCustomers() {
	for {
		if requestQueue.Len() > 0 {
			item := requestQueue.Pop()
			request := item.Value.(*CrawlRequest)
			fmt.Println(request.CustomerType, request.URL)
			wg.Add(1)
			crawlURL(request)
			wg.Done()
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func crawlURL(request *CrawlRequest) {
	urlStr := request.URL
	customerType := request.CustomerType
	cachedPage, err := readCachedPage(urlStr)
	if err == nil {
		request.UniqueURLs = cachedPage.CrawledURLs
		updateCacheModificationTime(urlStr)
		fmt.Println("cache usage")
		return
	}

	crawledURLs := make(map[string]struct{})
	visitedURLs := make(map[string]struct{})
	crawledHTML := ""
	maxRetries := 3
	if customerType == "paying" {
		maxRetries = 5
	}
	urlqueue := list.New()
	urlqueue.PushBack(urlStr)
	cntURLs := 0
	for {
		if urlqueue.Len() == 0 || cntURLs>30 {
			break
		}
		front := urlqueue.Front()
		urlqueue.Remove(front)
		parsedURL, err := url.Parse(front.Value.(string))
		if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || err != nil  {
			continue
		}
		// fmt.Println(front.Value)
		cntURLs++
		for retry := 0; retry < maxRetries; retry++ {
			c := colly.NewCollector()
			c.OnHTML("html", func(e *colly.HTMLElement) {
				crawledHTML = e.Text
			})
			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				link := e.Attr("href")
				if _, exists := crawledURLs[link]; !exists {
					crawledURLs[link] = struct{}{}
					if _, exists := visitedURLs[link]; !exists {
						urlqueue.PushBack(link)
						visitedURLs[link] = struct{}{}
					}
				}
			})
			err := c.Visit(front.Value.(string))
			if err != nil {
				// fmt.Println("Error visiting URL:", front.Value, err)
				continue
			} else {
				break
			}
		}
	}

	// Set unique URLs in the request
	request.UniqueURLs = make([]string, 0, len(crawledURLs))
	for link := range crawledURLs {
		request.UniqueURLs = append(request.UniqueURLs, link)
		// fmt.Println(link)
	}

	saveCachedPage(urlStr, crawledHTML, request.UniqueURLs)
}

func readCachedPage(urlStr string) (CachedPage, error) {
	filePath := "cache/" + urlToFileName(urlStr)
	if _, err := os.Stat(filePath); err == nil {
		fileStat, _ := os.Stat(filePath)
		if time.Since(fileStat.ModTime()) <= 60*time.Minute {
			html, err := os.ReadFile(filePath)
			if err == nil {
				var cachedPage CachedPage
				if err := json.Unmarshal(html, &cachedPage); err == nil {
					return cachedPage, nil
				}
			}
		}
	}
	return CachedPage{}, fmt.Errorf("no cached page available or outdated")
}

func saveCachedPage(urlStr string, html string, crawledURLs []string) {
	filePath := "cache/" + urlToFileName(urlStr)
	cachedPage := CachedPage{
		HTML:        html,
		CrawledURLs: crawledURLs,
	}
	data, _ := json.Marshal(cachedPage)
	os.WriteFile(filePath, data, 0644)
}

func urlToFileName(urlStr string) string {
	return strings.Replace(urlStr, "/", "_", -1)
}

func cleanupCache() {
	cacheDir := "cache"
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}
	for _, file := range files {
		fileStat, _ := file.Info()
		if time.Since(fileStat.ModTime()) > 60*time.Minute {
			os.Remove(filepath.Join(cacheDir, file.Name()))
		}
	}
}

func updateCacheModificationTime(urlStr string) {
	filePath := "cache/" + urlToFileName(urlStr)
	if _, err := os.Stat(filePath); err == nil {
		os.Chtimes(filePath, time.Now(), time.Now())
	}
}