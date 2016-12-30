package yupack

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	gq "github.com/puerkitobio/goquery"
)

func GetStatus(numbers []string) ([][]string, error) {
	var value url.Values
	r := [][]string{}
	for i, v := range numbers {
		if v == "" {
			continue
		}
		n := i%10 + 1
		if n == 1 {
			value = url.Values{
				"search.x": {"1"},
			}
		}
		b := true
		for _, temp := range value {
			for _, no := range temp {
				if v == no {
					b = false
				}
			}
		}
		if b {
			value.Add("requestNo"+fmt.Sprint(n), v)
		}
		if n == 10 && i > 0 {
			if d, e := get(value); e != nil {
				return nil, e
			} else {
				if len(value) > 2 {
					r = append(r, writeSome(d)...)
				} else {
					r = append(r, writeOne(d))
				}
			}
		}
	}
	if len(value) > 1 {
		if d, e := get(value); e != nil {
			return nil, e
		} else {
			if len(value) > 2 {
				r = append(r, writeSome(d)...)
			} else {
				r = append(r, writeOne(d))
			}
		}
	}
	return r, nil
}

func get(value url.Values) (*gq.Document, error) {
	res, err := http.Get("https://trackings.post.japanpost.jp/services/srv/search/?" + value.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := gq.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func writeOne(doc *gq.Document) []string {
	r := []string{}
	doc.Find(".tableType01").Each(func(_ int, s *gq.Selection) {
		if a, b := s.Attr("summary"); b && a == "配達状況詳細" {
			r = append(r, s.Find(".w_180").Last().Text())
		}
		if a, b := s.Attr("summary"); b && a == "履歴情報" {
			r = append(r, strings.Split(s.Find("tr .w_120").Last().Text(), " ")[0])
			r = append(r, s.Find("tr .w_150").Last().Text())
		}
	})
	return r
}

func writeSome(doc *gq.Document) [][]string {
	r := [][]string{}
	doc.Find(".tableType01").Each(func(_ int, s *gq.Selection) {
		if a, b := s.Attr("summary"); b && a == "照会結果" {
			s.Find("tr").Each(func(i int, s *gq.Selection) {
				if i%2 == 0 && i > 0 {
					tmp := []string{}
					tmp = append(tmp, s.Find(".w_120").First().Text())
					date := strings.Replace(s.Find(".w_80").First().Next().Text(), " ", "", -1)
					date = strings.Replace(date, "\n", "", -1)
					tmp = append(tmp, date[:10])
					tmp = append(tmp, s.Find(".w_180").First().Text())
					r = append(r, tmp)
				}
			})
		}
	})
	return r
}
