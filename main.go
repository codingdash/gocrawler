/**
 * Gocrawler
 *
 * @author Ritesh Shrivastav
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

// To track crawled URL
var processedLinks = make(map[string]int)

// Will keep reference to the given target URL
var targetUrl *url.URL

// Default file permission
var filePermission = os.FileMode(0777)

// Performs boot check wheather the cli arguments are acceptable, if invalid
// prints help message.
func InitialCheck() bool {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gocrawl http://example.com")
		return false
	}
	return true
}

// Based on content-type returns path to save the content
func GetSavePath(target string, contentType string) string {
	u, _ := url.Parse(target)
	newFilePath := "./" + targetUrl.Host + u.Path

	// creating dirs, in case if not present already
	if strings.Contains(contentType, "text/directory") {
		os.MkdirAll(newFilePath, filePermission)
	} else if strings.Contains(contentType, "text/html") {
		os.MkdirAll(newFilePath, filePermission)
		dirSep := "/"
		if newFilePath[len(newFilePath)-1:] == "/" {
			dirSep = ""
		}
		return newFilePath + dirSep + "index.html"
	} else {
		os.MkdirAll((newFilePath[:strings.LastIndex(newFilePath, "/")]), filePermission)
	}
	return newFilePath
}

// Performs check if given URL is already processed
func IsProcessed(target string) bool {
	if processedLinks[target] == 0 {
		// not processed
		return false
	} else {
		// processed
		return true
	}
}

// Saves the processed target Url with response code(in case if we need in
// future)
func SaveProcessed(target string, status int) {
	processedLinks[target] = status
}

// Saves given file(content) to the path
func StoreFile(path string, content []byte) {
	ioErr := ioutil.WriteFile(path, content, filePermission)
	if ioErr != nil {
		fmt.Println("Error while writing file " + path + "; " + ioErr.Error())
		os.Exit(2)
	}
}

// Process all css files from the given document
func ProcessStyle(doc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("href")
		if exist {
			u, _ := url.Parse(val)
			if u.Host == targetUrl.Host || u.Host == "" {
				newTarget := targetUrl.Scheme + "://" + targetUrl.Host + u.Path
				ProcessDoc(newTarget, 0)
			}
		}
	})
}

// Process all the javascript files from given document
func ProcessScript(doc *goquery.Document) {
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("src")
		if exist {
			u, _ := url.Parse(val)
			if u.Host == targetUrl.Host || u.Host == "" {
				newTarget := targetUrl.Scheme + "://" + targetUrl.Host + u.Path
				ProcessDoc(newTarget, 0)
			}
		}
	})
}

// Process all image from the given document
func ProcessImg(doc *goquery.Document) {
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("src")
		if exist {
			u, _ := url.Parse(val)
			if u.Host == targetUrl.Host || u.Host == "" {
				newTarget := targetUrl.Scheme + "://" + targetUrl.Host + u.Path
				ProcessDoc(newTarget, 0)
			}
		}
	})
}

// Process all the links from the given document
func ProcessLink(doc *goquery.Document, level int) {
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("href")
		if exist {
			u, _ := url.Parse(val)
			if u.Host == targetUrl.Host || u.Host == "" {
				newTarget := targetUrl.Scheme + "://" + targetUrl.Host + u.Path
				ProcessDoc(newTarget, level + 1)
			}
		}
	})
}

// Process and save the files, the main method which at the end calls other
// node processing methods. Currently [a, img, link, script] tags are being
// explored.
func ProcessDoc(target string, level int) {
	// check if processed already
	if IsProcessed(target) || level > 5 {
		return
	}
	// printing to console
	fmt.Printf("Fetching level %v from %v\n", level, target)
	// fetching target URL
	resp, fetchErr := http.Get(target)
	if fetchErr != nil {
		fmt.Println("Error while fetching " + target + "; " + fetchErr.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Got response code %v\n", resp.StatusCode)
		return
	}
	contentType := resp.Header.Get("Content-Type")
	// just save the file in case content type is not HTML
	if !strings.Contains(contentType, "text/html") {
		pathToSave := GetSavePath(target, contentType)
		content, readErr := ioutil.ReadAll(resp.Body)
		if readErr == nil {
			StoreFile(pathToSave, []byte(content))
		}
		SaveProcessed(target, resp.StatusCode)
	} else {
		// need to look further
		doc, docErr := goquery.NewDocument(target)
		if docErr != nil {
			SaveProcessed(target, resp.StatusCode)
			fmt.Println("Error while fetching " + target)
		} else {
			SaveProcessed(target, resp.StatusCode)
		}
		// saving doc
		pathToSave := GetSavePath(target, contentType)
		htmlStr, _ := doc.Html()
		StoreFile(pathToSave, []byte(htmlStr))
		// process child nodes
		ProcessImg(doc)
		ProcessLink(doc, level)
		ProcessScript(doc)
		ProcessStyle(doc)
	}
}

// Gocrawler main thread.
func main() {
	if !InitialCheck() {
		os.Exit(2)
	}
	target := os.Args[1]
	u, err := url.Parse(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	targetUrl = u
	ProcessDoc(target, 0)
}
