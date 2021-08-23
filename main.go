package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func main() {
	var query string
	flag.StringVar(&query, "query", "", "Search query")
	flag.Parse()

	if query == "" {
		fmt.Println("[ERR]", "Query value should be provided")
		return
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.Flag("headless", false), chromedp.NoFirstRun)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	searchBox := "input[name='q']"
	var attrs []map[string]string

	var next string
	var ok bool

	tasks := chromedp.Tasks{
		chromedp.WaitReady("body"),
		chromedp.AttributesAll(".yuRUbf > a", &attrs, chromedp.ByQueryAll),
		chromedp.AttributeValue("#pnnext", "href", &next, &ok, chromedp.AtLeast(0)),
	}

	chromedp.Run(ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.Focus(searchBox),
		chromedp.SendKeys(searchBox, query),
		chromedp.SendKeys(searchBox, kb.Enter),
		tasks,
	)

	for ok {
		time.Sleep(time.Second * 1)
		chromedp.Run(ctx,
			chromedp.Navigate("https://www.google.com"+next),
			tasks,
		)
	}

	for _, attr := range attrs {
		fmt.Println(attr["href"])
	}
}
