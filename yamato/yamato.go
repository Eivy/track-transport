package yamato

import (
	"fmt"
	"net/http"
	"net/url"

	gq "github.com/puerkitobio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetStatus(numbers []string) ([][]string, error) {
	var value url.Values
	r := [][]string{}
	for i, v := range numbers {
		if i%10 == 0 {
			value = url.Values{
				"number00": {"2"},
			}
		}
		value.Add("number"+fmt.Sprintf("%02d", i%10+1), v)
		if (i+1)%10 == 0 && i > 0 {
			if d, e := post(value); e != nil {
				return nil, e
			} else {
				r = append(r, parseResult(d)...)
			}
		}
	}
	if len(value) > 1 {
		if d, e := post(value); e != nil {
			return nil, e
		} else {
			r = append(r, parseResult(d)...)
		}
	}
	return r, nil
}

func post(value url.Values) (*gq.Document, error) {
	res, err := http.PostForm("https://toi.kuronekoyamato.co.jp/cgi-bin/tneko", value)
	if err != nil {
		return nil, err
	}
	r := transform.NewReader(res.Body, japanese.ShiftJIS.NewDecoder())
	doc, err := gq.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func parseResult(doc *gq.Document) [][]string {
	r := [][]string{}
	doc.Find(".ichiran").First().Find("tr").Each(func(_ int, s *gq.Selection) {
		if a, b := s.Attr("align"); b && a == "middle" {
			if _, b := s.Find("input").First().Attr("value"); b {
				tmp := []string{}
				tmp = append(tmp, s.Find(".denpyo").First().Text())
				tmp = append(tmp, s.Find(".hiduke").First().Text())
				tmp = append(tmp, s.Find(".ct").First().Text())
				r = append(r, tmp)
			}
		}
	})
	return r
}
