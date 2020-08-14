// Define the package - this will help maintain scope in your application
package main

// Define the fmt import for formatting
import (
	"fmt"
	"net/http"

	// You can prefix imports to make it easier to reference them in your code, like this one
	"github.com/golang-web-scraping/pkg/actions"
	logr "github.com/sirupsen/logrus"
)

func main() {
	// Print a formatted string to the terminal
	fmt.Println("Hello Fellow Gopher!")

	// Create the first route handler listening on '/'
	http.HandleFunc("/", home)
	http.HandleFunc("/scrape", actions.Scrape)
	http.HandleFunc("/crawl", actions.Crawl)
	logr.Info("Starting up on 8080")

	// Start the sever
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	// Assign the 'msg' variable with a string value
	msg := "Hello, welcome to your app. Use the following suffix's on the URL to show the different results.\n1)'/scrape' to show results of web scraping.\n2)'/crawl' to show results of a web crawler"

	// Logs a message to the terminal using the 3rd party import logrus
	logr.Info("Received request for the home page")

	// Write the response to the byte array - Sprintf formats and returns a string without printing it anywhere
	w.Write([]byte(fmt.Sprintf(msg)))
}
