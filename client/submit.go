package client

import (
	"fmt"
	"net/url"

	"github.com/sempr/cf/util"

	"github.com/fatih/color"
)

// Submit submit (block while pending)
func (c *Client) Submit(info Info, langID, source string) (err error) {
	color.Cyan("Submit " + info.Hint())

	URL, err := info.SubmitURL(c.host)
	if err != nil {
		return
	}

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	handle, err := findHandle(body)
	if err != nil {
		return
	}

	fmt.Printf("Current user: %v\n", handle)

	csrf, err := findCsrf(body)
	if err != nil {
		return
	}

	taskScreenName := fmt.Sprintf("%v_%v", info.ContestID, info.ProblemID)
	fmt.Println(taskScreenName, URL)
	_, err = util.PostBody(c.client, URL, url.Values{
		"csrf_token":          {csrf},
		"data.TaskScreenName": {taskScreenName},
		"data.LanguageId":     {langID},
		"sourceCode":          {source},
	})
	if err != nil {
		return
	}

	// errMsg, err := findErrorMessage(body)
	// if err == nil {
	// 	return errors.New(errMsg)
	// }

	// msg, err := findMessage(body)
	// if err != nil {
	// 	return errors.New("Submit failed")
	// }
	// if !strings.Contains(msg, "submitted successfully") {
	// 	return errors.New(msg)
	// }

	color.Green("Submitted")

	// if err != nil {
	// 	return errors.New("No submission Id")
	// }

	submissions, err := c.WatchSubmission(info, 1, true)
	if err != nil {
		return
	}

	info.SubmissionID = submissions[0].ParseID()
	c.Handle = handle
	c.LastSubmission = &info
	return c.save()
}
