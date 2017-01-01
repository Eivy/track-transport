package sagawa

import (
	"strings"

	gq "github.com/puerkitobio/goquery"
)

func GetStatus(numbers []string) ([][]string, error) {
	r := [][]string{}
	for _, v := range numbers {
		doc, err := gq.NewDocument("http://k2k.sagawa-exp.co.jp/p/web/okurijosearch.do?okurijoNo=" + v)
		if err != nil {
			return nil, err
		}
		r = append(r, parseResult(doc)...)
	}
	return r, nil
}

func parseResult(doc *gq.Document) [][]string {
	r := [][]string{}
	doc.Find(".table_okurijo_index").First().Each(func(_ int, s *gq.Selection) {
		tmp := []string{}
		tmp = append(tmp, s.Find("strong").First().Text())
		date := s.Find("tr").Last().Find("td").Text()
		date = strings.Split(date, "　")[0]
		date = strings.Replace(date, "月", "/", -1)
		date = strings.Replace(date, "日", "", -1)
		date = strings.Replace(date, "\t", "", -1)
		date = strings.Replace(date, "\n", "", -1)
		date = strings.Replace(date, " ", "", -1)
		tmp = append(tmp, date)
		status := s.Find("dt").First().Text()
		status = strings.Replace(status, "\n", "", -1)
		status = strings.Replace(status, " ", "", -1)
		status = strings.Replace(status, "\t", "", -1)
		tmp = append(tmp, status)
		r = append(r, tmp)
	})
	return r
}
