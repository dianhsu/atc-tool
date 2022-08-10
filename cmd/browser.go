package cmd

import (
	"github.com/dianhsu/atc/client"
	"github.com/dianhsu/atc/config"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
)

func openURL(url string) error {
	color.Green("Open %v", url)
	return open.Run(url)
}

// Open command
func Open() (err error) {
	URL, err := Args.Info.OpenURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Stand command
func Stand() (err error) {
	URL, err := Args.Info.StandingsURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Sid command
func Sid() (err error) {
	info := Args.Info
	if info.SubmissionID == "" && client.Instance.LastSubmission != nil {
		info = *client.Instance.LastSubmission
	}
	URL, err := info.SubmissionURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}
