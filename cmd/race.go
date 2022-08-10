package cmd

import (
	"fmt"
	"time"

	"github.com/dianhsu/atc/client"
	"github.com/dianhsu/atc/config"
)

// Race command
func Race() (err error) {
	cfg := config.Instance
	cln := client.Instance
	info := Args.Info
	if err = cln.RaceContest(info); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = cln.RaceContest(info)
		}
	}
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	URL, err := info.ProblemSetURL(cfg.Host)
	if err != nil {
		return
	}
	fmt.Printf("Open ProblemSetURL %s\n", URL)
	// openURL(URL)
	// openURL(URL + "/problems")
	return Parse()
}
