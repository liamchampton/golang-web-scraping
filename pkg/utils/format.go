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
