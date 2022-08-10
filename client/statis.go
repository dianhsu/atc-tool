package client

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/dianhsu/atc/util"
)

// StatisInfo statis information
type StatisInfo struct {
	ID     string
	Name   string
	IO     string
	Limit  string
	Passed string
	State  string
}

func findProblems(body []byte) ([]StatisInfo, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var ret []StatisInfo
	doc.Find("tbody").Children().Each(func(i1 int, s1 *goquery.Selection) {
		var p StatisInfo
		s1.Find("td").Each(func(i2 int, s2 *goquery.Selection) {
			switch i2 {
			case 0:
				problemUrl, _ := s2.Find("a").Attr("href")
				p.ID = problemUrl[len(problemUrl)-1:]
			case 1:
				p.Name = s2.Text()
			case 2:
				p.Limit = s2.Text()
			}
		})
		if p.ID != "" {
			ret = append(ret, p)
		}
	})
	return ret, nil
}

// Statis get statis
func (c *Client) Statis(info Info) (problems []StatisInfo, err error) {
	URL, err := info.ProblemSetURL(c.host)
	if err != nil {
		return
	}
	if info.ProblemType == "acmsguru" {
		return nil, errors.New(ErrorNotSupportAcmsguru)
	}

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	if _, err = findHandle(body); err != nil {
		return
	}

	return findProblems(body)
}
