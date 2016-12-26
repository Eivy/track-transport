package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"

	gq "github.com/puerkitobio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	if len(os.Args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			value := url.Values{
				"number00": {"2"},
				"number01": {scanner.Text()},
			}
			if err := getStatus(value); err != nil {
				exitWithError(err)
			}
		}
		if err := scanner.Err(); err != nil {
			exitWithError(err)
		}
	} else {
		for i, v := range os.Args {
			if i > 0 {
				fp, err := os.Open(v)
				if err != nil {
					exitWithError(err)
				}
				defer fp.Close()
				scanner := bufio.NewScanner(fp)
				i := 1
				value := url.Values{
					"number00": {"2"},
				}
				for scanner.Scan() {
					value.Add("number"+fmt.Sprintf("%02d", i), scanner.Text())
					if i == 10 {
						if err := getStatus(value); err != nil {
							exitWithError(err)
						}
						value = url.Values{
							"number00": {"2"},
						}
					}
					i = i%10 + 1
				}
				if err := scanner.Err(); err != nil {
					exitWithError(err)
				}
				if i != 10 {
					if err := getStatus(value); err != nil {
						exitWithError(err)
					}
				}
			}
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func getStatus(value url.Values) error {
	res, err := http.PostForm("https://toi.kuronekoyamato.co.jp/cgi-bin/tneko", value)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := transform.NewReader(res.Body, japanese.ShiftJIS.NewDecoder())
	doc, err := gq.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	doc.Find(".ichiran").First().Find("tr").Each(func(_ int, s *gq.Selection) {
		if a, b := s.Attr("align"); b && a == "middle" {
			if _, b := s.Find("input").First().Attr("value"); b {
				fmt.Print(s.Find(".denpyo").First().Text())
				fmt.Print(",")
				fmt.Print(s.Find(".hiduke").First().Text())
				fmt.Print(",")
				fmt.Println(s.Find(".ct").First().Text())
			}
		}
	})
	return nil
}
