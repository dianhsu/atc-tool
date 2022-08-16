package client

type Info struct {
	ContestID    string `json:"contest_id"`
	ProblemID    string `json:"problem_id"`
	SubmissionID string `json:"submission_id"`
	RootPath     string
}
