package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	var query string
	var out string
	flag.StringVar(&query, "query", "", "The search query to use")
	flag.StringVar(&out, "out", "", "File to write the output")
	flag.Parse()

	printScriptName()

	if query == "" {
		fmt.Println("[ERR]", "The query parameter must be provided")
		return
	}

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

	if out != "" {
		saveOutput(out, urls)
	}

	fmt.Printf("\n[INFO] Found %d domains for (%s)\n", len(urls), query)
}

// Parse URLs from the page and print them to stdout
func getURLs(page *rod.Page, urls *[]string) {
	anchors, _ := page.Elements(".yuRUbf > a")
	for _, anchor := range anchors {
		href, _ := anchor.Attribute("href")
		fmt.Println(*href)
		*urls = append(*urls, *href)
	}
}

// Write the extracted URLs to the specified file
func saveOutput(out string, urls []string) {
	f, err := os.Create(out)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	for _, url := range urls {
		f.WriteString(url + "\n")
	}
}

// Print the script name as big text
func printScriptName() {
	fmt.Println(`
░██████╗░░█████╗░░█████╗░░██████╗░██╗░░░██╗██████╗░██╗░░░░░
██╔════╝░██╔══██╗██╔══██╗██╔════╝░██║░░░██║██╔══██╗██║░░░░░
██║░░██╗░██║░░██║██║░░██║██║░░██╗░██║░░░██║██████╔╝██║░░░░░
██║░░╚██╗██║░░██║██║░░██║██║░░╚██╗██║░░░██║██╔══██╗██║░░░░░
╚██████╔╝╚█████╔╝╚█████╔╝╚██████╔╝╚██████╔╝██║░░██║███████╗
░╚═════╝░░╚════╝░░╚════╝░░╚═════╝░░╚═════╝░╚═╝░░╚═╝╚══════╝          											 
	`)
}
