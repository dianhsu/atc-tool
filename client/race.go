package client

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/dianhsu/atc/util"

	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
)

func findCountdown(body []byte) (time.Time, error) {

	reg := regexp.MustCompile(`var startTime = moment\("(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2})"\);`)
	tmp := reg.FindSubmatch(body)
	if tmp == nil {
		return time.Time{}, errors.New("cannot find any countdown")
	}
	return time.Parse("2006-01-02T15:04:05-07:00", string(tmp[1]))
}

// RaceContest wait for contest starting
func (c *Client) RaceContest(info Info) (err error) {
	color.Cyan("Race " + info.Hint())

	URL := c.host + "/contests/" + info.ContestID

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	_, err = findHandle(body)
	if err != nil {
		return
	}

	startTime, err := findCountdown(body)
	if err != nil {
		return err
	}
	now := time.Now()
	count := int(startTime.Sub(now).Seconds())
	if count > 0 {
		color.Green("Countdown: ")
	}
	for count > 0 {
		h := count / 60 / 60
		m := count/60 - h*60
		s := count - h*60*60 - m*60
		fmt.Printf("%02d:%02d:%02d\n", h, m, s)
		ansi.CursorUp(1)
		count--
		time.Sleep(time.Second)
	}
	time.Sleep(900 * time.Millisecond)

	return
}
