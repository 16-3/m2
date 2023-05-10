package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PostBody struct {
	URL string `json:"url"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var postData PostBody
		err = json.Unmarshal(body, &postData)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		if postData.URL == "" {
			http.Error(w, "Missing URL", http.StatusBadRequest)
			return
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var html string
		err = chromedp.Run(ctx,
			chromedp.Navigate(postData.URL),
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

	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var postData PostBody
		err = json.Unmarshal(body, &postData)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		if postData.URL == "" {
			http.Error(w, "Missing URL", http.StatusBadRequest)
			return
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var screenshot []byte
		err = chromedp.Run(ctx,
			chromedp.Navigate(postData.URL),
			chromedp.ScrollIntoView(`document.body`, chromedp.ByJSPath),
			chromedp.CaptureScreenshot(&screenshot),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(screenshot)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
