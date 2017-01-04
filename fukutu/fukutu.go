package fukutu

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	gq "github.com/puerkitobio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetStatus(numbers []string) ([][]string, error) {
	j, _ := cookiejar.New(nil)
	c := &http.Client{Jar: j}
	res, err := c.Get("https://corp.fukutsu.co.jp/situation/tracking_no")
	if err != nil {
		return nil, err
	}
	hidden, err := getHidenItem(res)
	if err != nil {
		return nil, err
	}
	value := url.Values{}
	for k, v := range hidden {
		value.Add(k, v)
	}
	r := [][]string{}
	for i, v := range numbers {
		value.Add("data[TrackingNo][tracking_no"+fmt.Sprint(i%10+1)+"]", v)
		if (i+1)%10 == 0 && i > 0 {
			if d, e := post(c, value); e != nil {
				return nil, e
			} else {
				r = append(r, parseResult(d)...)
			}
		}
	}
	if d, e := post(c, value); e != nil {
		return nil, e
	} else {
		r = append(r, parseResult(d)...)
	}
	return r, nil
}

func getHidenItem(res *http.Response) (map[string]string, error) {
	r := map[string]string{}
	doc, err := gq.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	doc.Find("form input").Each(func(_ int, s *gq.Selection) {
		if a, b := s.Attr("type"); b && a == "hidden" {
			if a, b := s.Attr("value"); b {
				if n, b2 := s.Attr("name"); b2 {
					r[n] = a
				}
			}
		}
	})
	return r, nil
}

func post(client *http.Client, value url.Values) (*gq.Document, error) {
	req, err := http.NewRequest("POST", "https://corp.fukutsu.co.jp/situation/tracking_no", strings.NewReader(value.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
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
	doc.Find("#contents>table.table01").Each(func(_ int, s *gq.Selection) {
		s.Find(".address_content").Each(func(_ int, s *gq.Selection) {
			tmp := make([]string, 0)
			s.Find("td").Each(func(i int, s *gq.Selection) {
				if i == 0 || i == 2 || i == 3 {
					t := strings.Replace(s.Text(), "\u00a0", " ", -1) // No Breaking Space
					t = strings.Replace(t, " ", "", -1)
					t = strings.Replace(t, "\n", "", -1)
					t = strings.Replace(t, "\b", "", -1)
					t = strings.Replace(t, "\r", "", -1)
					t = strings.Replace(t, "\t", "", -1)
					tmp = append(tmp, t)
				} else if i == 4 {
					t := strings.Replace(s.Text(), "\u00a0", " ", -1) // No Breaking Space
					t = strings.Replace(t, " ", "", -1)
					t = strings.Replace(t, "\n", "", -1)
					t = strings.Replace(t, "\b", "", -1)
					t = strings.Replace(t, "\r", "", -1)
					t = strings.Replace(t, "\t", "", -1)
					tmp[len(tmp)-1] = tmp[len(tmp)-1] + t
				}
			})
			r = append(r, tmp)
		})
	})
	return r
}
