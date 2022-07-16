package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sempr/cf/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

// Submission submit state
type Submission struct {
	name   string
	id     uint64
	status string
	passed uint64
	judged uint64
	points uint64
	time   uint64
	memory uint64
	lang   string
	when   string
	end    bool
}

func isWait(verdict string) bool {
	return verdict == "null" || verdict == "TESTING" || verdict == "SUBMITTED"
}

// ParseStatus with color
func (s *Submission) ParseStatus() string {
	status := strings.ReplaceAll(s.status, "${f-points}", fmt.Sprintf("%v", s.points))
	status = strings.ReplaceAll(status, "${f-passed}", fmt.Sprintf("%v", s.passed))
	status = strings.ReplaceAll(status, "${f-judged}", fmt.Sprintf("%v", s.judged))
	for k, v := range colorMap {
		tmp := strings.ReplaceAll(status, k, "")
		if tmp != status {
			status = color.New(v).Sprint(tmp)
		}
	}
	return status
}

// ParseID formatter
func (s *Submission) ParseID() string {
	return fmt.Sprintf("%v", s.id)
}

// ParseMemory formatter
func (s *Submission) ParseMemory() string {
	if s.memory > 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(s.memory)/1024.0/1024.0)
	} else if s.memory > 1024 {
		return fmt.Sprintf("%.2f KB", float64(s.memory)/1024.0)
	}
	return fmt.Sprintf("%v B", s.memory)
}

// ParseTime formatter
func (s *Submission) ParseTime() string {
	return fmt.Sprintf("%v ms", s.time)
}

// ParseProblemIndex get problem's index
func (s *Submission) ParseProblemIndex() string {
	p := strings.Index(s.name, " ")
	if p == -1 {
		return ""
	}
	return strings.ToLower(s.name[:p])
}

func refreshLine(n int, maxWidth int) {
	for i := 0; i < n; i++ {
		ansi.Printf("%v\n", strings.Repeat(" ", maxWidth))
	}
	ansi.CursorUp(n)
}

func updateLine(line string, maxWidth *int) string {
	*maxWidth = len(line)
	return line
}

func (s *Submission) display(first bool, maxWidth *int) {
	if !first {
		ansi.CursorUp(7)
	}
	ansi.Printf("      #: %v\n", s.ParseID())
	ansi.Printf("   when: %v\n", s.when)
	ansi.Printf("   prob: %v\n", s.name)
	ansi.Printf("   lang: %v\n", s.lang)
	refreshLine(1, *maxWidth)
	ansi.Printf(updateLine(fmt.Sprintf(" status: %v\n", s.ParseStatus()), maxWidth))
	ansi.Printf("   time: %v\n", s.ParseTime())
	ansi.Printf(" memory: %v\n", s.ParseMemory())
}

func display(submissions []Submission, problemID string, first bool, maxWidth *int, line bool) {
	if line {
		submissions[0].display(first, maxWidth)
		return
	}
	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"#", "when", "problem", "lang", "status", "time", "memory"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, sub := range submissions {
		if problemID != "" && sub.ParseProblemIndex() != problemID {
			continue
		}
		table.Append([]string{
			sub.ParseID(),
			sub.when,
			sub.name,
			sub.lang,
			sub.ParseStatus(),
			sub.ParseTime(),
			sub.ParseMemory(),
		})
	}
	table.Render()

	if !first {
		ansi.CursorUp(len(submissions) + 2)
	}
	refreshLine(len(submissions)+2, *maxWidth)

	scanner := bufio.NewScanner(io.Reader(&buf))
	for scanner.Scan() {
		line := scanner.Text()
		*maxWidth = len(line)
		ansi.Println(line)
	}
}

func findEndStatus(text string) bool {
	if text == "WJ" {
		return false
	}
	if strings.Contains(text, "/") && !strings.Contains(text, " ") {
		return false
	}
	return true
}

func findSubmission(body []byte, n int) (submissions []Submission, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	vv := doc.Find("tbody tr")

	var sub Submission
	for i := 0; i < vv.Size(); i++ {
		z := vv.Eq(i).Children()
		sidStr, _ := z.Eq(4).Attr("data-id")
		sid, _ := strconv.ParseInt(sidStr, 10, 64)
		sub.id = uint64(sid)
		stText := strings.Trim(z.Eq(6).Text(), " ")
		sub.end = findEndStatus(stText)
		sub.status = stText
		if z.Length() == 10 {
			memStr := strings.Split(z.Eq(8).Text(), " ")[0]
			mem, _ := strconv.ParseInt(memStr, 10, 64)
			sub.memory = uint64(mem) * 1024
			timeStr := strings.Split(z.Eq(7).Text(), " ")[0]
			time_, _ := strconv.ParseInt(timeStr, 10, 64)
			sub.time = uint64(time_)
		}
		sub.when = z.Eq(0).Text()
		sub.name = z.Eq(1).Text()
		sub.lang = z.Eq(3).Text()
		submissions = append(submissions, sub)
	}

	return
}

// var ruTime = "DD.MM.YYYY HH:mm";
// var enTime = "MMM/DD/YYYY HH:mm";
// https://github.com/go-shadow/moment/blob/master/moment_parser.go
const ruTime = "02.01.2006 15:04 Z07:00"
const enTime = "Jan/02/2006 15:04 Z07:00"

func parseWhen(raw, cfOffset string) string {
	data := fmt.Sprintf("%v %v", raw, cfOffset)
	tm, err := time.Parse(ruTime, data)
	if err != nil {
		tm, _ = time.Parse(enTime, data)
	}
	return tm.In(time.Local).Format("2006-01-02 15:04")
}

func parseSubmission(body []byte, cfOffset string) (ret Submission, err error) {
	data := fmt.Sprintf("<table><tr %v</table>", string(body))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return
	}
	get := func(sel string) string {
		return strings.TrimSpace(doc.Find(sel).Text())
	}
	reg := regexp.MustCompile(`\d+`)
	getInt := func(sel string) uint64 {
		if tmp := reg.FindString(doc.Find(sel).Text()); tmp != "" {
			t, _ := strconv.Atoi(tmp)
			return uint64(t)
		}
		return 0
	}
	sub := doc.Find(".submissionVerdictWrapper")
	end := false
	if verdict, exist := sub.Attr("submissionverdict"); exist && !isWait(verdict) {
		end = true
	}
	status, _ := sub.Html()
	numReg := regexp.MustCompile(`\d+`)
	fmtReg := regexp.MustCompile(`<span\sclass=["']?verdict-format-([\S^>]+?)["']?>`)
	colReg := regexp.MustCompile(`<span\sclass=["']?verdict-([\S^>]+?)["']?>`)
	tagReg := regexp.MustCompile(`<[\s\S]*?>`)
	status = fmtReg.ReplaceAllString(status, "")
	status = colReg.ReplaceAllString(status, `${c-$1}`)
	status = tagReg.ReplaceAllString(status, "")
	status = strings.TrimSpace(status)
	when := get(".format-time")
	if when != "" {
		when = parseWhen(when, cfOffset)
	} else {
		when = strings.TrimSpace(doc.Find("td").First().Next().Text())
	}
	if status == "" {
		status = "Unknown"
	}
	var num uint64
	if s := numReg.FindString(status); s != "" {
		n, _ := strconv.Atoi(s)
		num = uint64(n)
	}
	return Submission{
		id:     getInt(".id-cell"),
		name:   get("td[data-problemId]"),
		lang:   get("td:not([class])"),
		status: status,
		time:   getInt(".time-consumed-cell"),
		memory: getInt(".memory-consumed-cell") * 1024,
		when:   when,
		passed: num,
		judged: num,
		points: num,
		end:    end,
	}, nil
}

func (c *Client) getSubmissions(URL string, n int) (submissions []Submission, err error) {
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	submissions, err = findSubmission(body, n)

	if err != nil {
		return
	}

	if len(submissions) < 1 {
		return nil, errors.New("cannot find any submission")
	}

	return
}

// WatchSubmission n is the number of submissions
func (c *Client) WatchSubmission(info Info, n int, line bool) (submissions []Submission, err error) {
	URL := fmt.Sprintf("%v/contests/%v/submissions/me", c.host, info.ContestID)

	maxWidth := 0
	first := true
	for {
		st := time.Now()
		submissions, err = c.getSubmissions(URL, n)
		if err != nil {
			return
		}
		display(submissions, info.ProblemID, first, &maxWidth, line)
		first = false
		endCount := 0
		for _, submission := range submissions {
			if submission.end {
				endCount++
			}
		}
		if endCount == len(submissions) {
			return
		}
		sub := time.Since(st)
		if sub < time.Second {
			time.Sleep(time.Duration(time.Second - sub))
		}
	}
}

var colorMap = map[string]color.Attribute{
	"${c-waiting}":  color.FgWhite,
	"${c-failed}":   color.FgRed,
	"${c-accepted}": color.FgGreen,
	"${c-rejected}": color.FgBlue,
}
