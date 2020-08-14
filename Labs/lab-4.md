# Lab 4 - Scrape the web :page_facing_up:

In this lab you are going to write a function to scrape a web page. It is important to understand that this is different to web crawling. Web scraping will copy the HTML and you can then programmatically target individual elements, whereas web crawling will follow embedded URLs until it reaches a page where no URLs are present.

This is the product you will be analysing. You will be scraping 3 parts, the name, number of stars and the price.
![ProductElements](../images/ProductElements.png)

To do this, the code will search through the `HTML` and find the individual `<span>` tags we outline. It will then scrape the contents.

You can see the `HTML` elements by using the browser inspector since this is client-side code.
![ProductName](../images/ProductName.png)


![ProcuctStars](../images/ProductStars.png)


![ProductPrice](../images/ProductPrice.png)

### Step 1

1. Navigate to the `/pkg/actions/scaping.go` file
2. In this file, add the code below (read the code carefully, so you understand what is happening)

```golang

package actions

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	//"github.com/golang-web-scraping/pkg/utils"
	"net/http"
	"os"
)

// Product struct will hold the information we want to print out from the web page
type Product struct {
	Name string
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

			// ChildText extracts the sanitised string from the within the matched element
			// In this case... the product name
			productName = e.ChildText("span.a-size-medium.a-color-base.a-text-normal")

			if productName == "" {
				// If we can't get any name, we return and go directly to the next element
				return
			}

			// In this case... the stars
			stars = e.ChildText("span.a-icon-alt")

			// Call the helper function to format the stars into a float (decimal) e.g 4.8
			//format.FormatStars(&stars)

			// In this case... the price
			price = e.ChildText("span.a-price > span.a-offscreen")
			if price == "" {
				// If we can't get any price, we return and go directly to the next element
				return
			}

			// Format the price so it is readable - some prices may have a 'was' and a 'now' price
			// This helper function will strip the old price from the string
			//format.FormatPrice(&price)

			//fmt.Printf("Product Name: %s \nStars: %s \nPrice: %s \n", productName, stars, price)

			// Append the collected data (Name, Stars, Price) from the element into the empty slice
			dataSlice = append(dataSlice, Product{
				Name: productName,
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

```

> **Note**: You will notice the import `"github.com/gocolly/colly"`. This is a 3rd party import that has been written to make web scraping much easier.

> **Note**: You may have also noticed the use of `&` and `*` operators throughout the code and wondered what they are. In Golang, these represent pointers (memory address locations). For some extra reading, code examples and understanding, more on this can be found at https://www.golang-book.com/books/intro/8

### Step 2

Now the code is written, you need to add a function call into a route handler on the sever.

Just like you did in [Lab 2](./lab-2.md), you need to create a `http.HandleFunc()` in your `main()` function.

Below the line `http.HandleFunc("/", home)` in your `main()` function in `main.go`, add the line, `http.HandleFunc("/scrape", actions.Scrape)`. Make sure the import for the actions package is also visible. If it is not, in the `import` block  add `"github.com/golang-web-scraping/pkg/actions"` (make sure the system path is correct for your project).

This route handler will call the `Scrape()` function when the route `/scrape` is hit and print out the results to the screen.

Push the application code up to Cloud Foundry just like you did in [Lab 3](./lab-3.md). Use the command `ibmcloud cf push` from within the root directory of the project.

### Step 3

If you run the application as it is now, you will have a fairly jumbled output on some items. Items with a 'was', and a 'now' price may look like "£12.99£15.99", and the stars will be more than one float number e.g  "4.6 out of 5 stars".
 
We want the output to look something like...
{
"Name": "Camera",
"Stars": "4.6",
"Price": "£29.99"
}

To do this, you need to write 2 helper functions. One for the price formatting and one for the stars formatting.

This is really easy to do. Follow the next steps in order to do this for your project:
1. Navigate to the `/pkg/utils/format.go` file
2. Add the following code snippet and read it carefully to understand what is happening. Make sure the import for the `utils` package is added in the `scraping.go` imports. If it is not there, the code will fail to compile - (`"github.com/web-scraping/pkg/utils"`). 

```golang

package format

import (
	"strings"
)

func FormatPrice(price *string) {
    
        // Count the number of £'s present in the string 
    	r := strings.Count(*price, "£")
    
    	// If >= 1 £'s in the string, split it and return the 'now' price
    	if r >= 1 {
        	splitStr := strings.Split(*price, "£")
        	*price = "£" + splitStr[1]
        }

}

func FormatStars(stars *string) {
    // Take the first 3 chars of the stars string (e.g 4.8)
    if len(*stars) >= 3 {
    		*stars = (*stars)[0:3]
    } else {
    		*stars = "Unknown"
    }
}

```

### Step 4

Now the helper functions have been written, go back to `Scrape()` function and uncomment the following two lines:

1. `format.FormatStars(&stars)`
2. `format.FormatPrice(&price)`
3. `"github.com/golang-web-scraping/pkg/utils"` (this line is in the imports)

### Step 5

Push the application code up to the cloud again...

1. Ensure all the project code is saved.
2. Ensure you are still signed in to your `ibmcloud` account from the terminal. If your are not logged in, follow the login instructions on [Lab 3](./lab-3.md) step 3.
2. In your terminal, make sure you are in the root directory of your project and enter the command `ibmcloud cf push`. This will push the new code up to the project you already have running in `ibmcloud`.
3. In the URL address bar of a web browser, enter the project address which can be found in the terminal output once the app has started. It will be the value of `routes`. This will show you the home route of the sever.
4. Append `/scrape` to the end of the URL to see the output of the web scraping function.
5. You should see the JSON output on the page for each product you are scraping.

The next stage of this workshop will be writing a function to crawl and show you how they differ. Continue onto [Lab 5](./lab-5.md)
