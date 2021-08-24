package main

import (
	"flag"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	flag.Parse()
	query := flag.Arg(0)

	page := rod.New().MustConnect().MustPage("https://www.google.com")

	searchBox, _ := page.Element("input[name='q']")
	searchBox.Focus()
	searchBox.Input(query)
	searchBox.Press(input.Enter)
	var urls []string

	for {
		page.WaitLoad()

		selector := "#pnnext"
		getURLs(page, &urls)

		exists := page.MustHas(selector)

		if exists {
			url := page.MustElement(selector).MustAttribute("href")
			page.Navigate("https://www.google.com" + *url)
		} else {
			break
		}
	}

	for _, u := range urls {
		fmt.Println(u)
	}
	fmt.Printf("\n[INFO] Found %d domains for (%s)\n", len(urls), query)
}

func getURLs(page *rod.Page, urls *[]string) {
	anchors, _ := page.Elements(".yuRUbf > a")
	for _, anchor := range anchors {
		href, _ := anchor.Attribute("href")
		*urls = append(*urls, *href)
	}
}
