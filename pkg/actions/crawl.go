package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func Crawl(w http.ResponseWriter, r *http.Request) {

	// Array containing all the known URLs in a sitemap
	var knownUrls []string

	// Create a Collector specifically for Google using a whitelist for the domain
	c := colly.NewCollector(colly.AllowedDomains("www.google.com"))

	// Create a callback on the XPath query searching for the URLs in the sitemap
	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		knownUrls = append(knownUrls, e.Text)
	})

	// Start the collecting job at the following URL
	c.Visit("https://www.google.com/sitemap.xml")

	numURLs := strconv.Itoa(len(knownUrls))

	// Push URL found to the slice
	knownUrls = append(knownUrls, "Collected "+numURLs+" URLs")

	urlStr := strings.Join(knownUrls, "\n")

	// Write the response to the byte array - Sprintf formats and returns a string without printing it anywhere
	w.Write([]byte(fmt.Sprintf(urlStr)))
}
