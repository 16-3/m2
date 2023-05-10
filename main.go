package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing URL", http.StatusBadRequest)
			return
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.ScrollIntoView(`document.body`, chromedp.ByJSPath),
			chromedp.OuterHTML("html", &html),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		strings.NewReader(html).WriteTo(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
