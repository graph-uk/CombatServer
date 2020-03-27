package jobs

// Job represents API response for acquire job method
type Job struct {
	CaseID      int
	CommandLine string
	Content     string
}
