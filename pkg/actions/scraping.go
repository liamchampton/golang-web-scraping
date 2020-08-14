package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gocolly/colly"
	"github.com/golang-web-scraping/pkg/utils"
)

// Product struct will hold the information we want to print out from the web page
type Product struct {
	Name  string
	Stars string
	Price string
}

/* Scrape is an exported function
** This lets functions outside of the "actions" package call it (makes it globally available)
 */
func Scrape(w http.ResponseWriter, r *http.Request) {
	// Write the status code 200
	w.WriteHeader(http.StatusOK)

	// create a new collector with the colly HTTP framework
	c := colly.NewCollector()

	// Execute a new request on every call from the collector
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Create an empty slice based on the Product struct (Name, Stars, Price) that you will populate with data
	var dataSlice []Product

	// Uses goquery to find an element of the HTML that we want to extract
	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {

		// There are many elements in the above HTML so use a ForEach to loop over them and then call the callback function
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			var productName, stars, price string

			// ChildText extracts the sanaitsed string from the within the matched element
			// In this case... the product name
			productName = e.ChildText("span.a-size-medium.a-color-base.a-text-normal")

			if productName == "" {
				// If we can't get any name, we return and go directly to the next element
				return
			}

			// In this case... the stars
			stars = e.ChildText("span.a-icon-alt")

			// Call the helper function to format the stars into a float (decimal) e.g 4.8
			utils.FormatStars(&stars)

			// In this case... the price
			price = e.ChildText("span.a-price > span.a-offscreen")
			if price == "" {
				// If we can't get any price, we return and go directly to the next element
				return
			}

			// Format the price so it is readable - some prices may have a 'was' and a 'now' price
			// This helper function will strip the old price from the string
			utils.FormatPrice(&price)

			//fmt.Printf("Product Name: %s \nStars: %s \nPrice: %s \n", productName, stars, price)

			// Append the collected data (Name, Stars, Price) from the element into the empty slice
			dataSlice = append(dataSlice, Product{
				Name:  productName,
				Stars: stars,
				Price: price,
			})
		})

		// Marshal the JSON into a readable format using json.MarshalIndent
		result, err := json.MarshalIndent(dataSlice, "", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Write the response to the byte array - Sprintf formats and returns a string without printing it anywhere
		w.Write([]byte(fmt.Sprintf(string(result))))

		// **NOTE** You could also write this response to a file instead of a stream if you wanted to store the data locally or into a database...
	})

	// Start the collecting job at the following URL
	c.Visit("https://www.amazon.co.uk/s?k=gopro&ref=nb_sb_noss_2")

}
