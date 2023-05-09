package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	url := flag.String("url", "https://www.amazon.in/s?i=electronics&bbn=1805560031&rh=n%3A976419031%2Cn%3A1389401031%2Cn%3A1389432031%2Cn%3A1805560031%2Cp_36%3A1000100-2500000%2Cp_n_condition-type%3A8609960031&pd_rd_r=605acd42-5051-4ec3-a144-084102cf7ef6&pd_rd_w=BZvVI&pd_rd_wg=6mROC&pf_rd_p=6e9c5ebb-d370-421b-8375-bf50155e0300&pf_rd_r=GJBQTS6X1AP98KR1H799&ref=tile2_10to20K", "URL to scrape")

	flag.Parse()

	phoneList := scraper(*url)
	jsonifier(phoneList)

}

func err_handler(msg string) {
	fmt.Printf("Error: %s\n", msg)
	os.Exit(1)
}

func scraper(url string) []Phone {
	resp, err := http.Get(url)
	// makes a GET request to the url

	fmt.Println(resp.Status)

	if err != nil {
		err_handler("Can't fetch the url")
	}

	var phoneList []Phone

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	// goquery.NewDocumentFromReader takes a reader as an argument

	if err != nil {
		err_handler("Goquery can't make up the document")
	}

	doc.Find(".sg-col-4-of-24").Each(func(_ int, class *goquery.Selection) {
		name := class.Find("span .a-size-base-plus").Text()
		name_first, err := name_modifier(name)

		if err != nil {
			fmt.Println("Can't handle empty text for splitting")
		} else {
			price := class.Find("span .a-price-whole").Text()
			//fmt.Println(price)
			// find function finds the specified tag from the buffer
			price_int := price_modifier(price)

			crr_phone := Phone{
				Name:  name_first,
				Price: price_int,
			}
			phoneList = append(phoneList, crr_phone)
		}

	})

	defer resp.Body.Close()
	// always remember to close the body of the response
	return phoneList
}

func name_modifier(s string) (string, error) {
	string_list := strings.Split(s, ")")
	// splits the string into a slice where 1st part is the string before delimitter and the 2nd part is after it

	if len(string_list) == 1 {
		return "", errors.New("no text to split")
	}

	string_first := strings.TrimSpace(string_list[0])

	return string_first + ")", nil
}

func price_modifier(s string) int {
	string_list := strings.Split(s, ",")

	var price_string string

	if len(string_list) > 1 {

		price_string = string_list[0] + string_list[1]
	} else {
		price_string = string_list[0]
	}
	price, _ := strconv.Atoi(price_string)
	//strconv.Atoi converts a string type to int
	return price
}

func jsonifier(phoneList []Phone) {

	//creating the output json file
	file, err := os.Create("output.json")
	defer file.Close()

	if err != nil {
		err_handler("Can't create the file")
	}

	jsonData, err := json.Marshal(phoneList)
	// this will jsonify the struct data
	if err != nil {
		err_handler("can't jsonify the given data")
	}
	n, err := io.WriteString(file, string(jsonData))
	if err != nil {
		err_handler("can't write to the file")
	}
	fmt.Printf("Successfully written %d bytes\n", n)
}

type Phone struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
