package sessions

// CaseItem ...
type CaseItem struct {
	Title             string    `json:"title"`
	Status            string    `json:"status"`
	Tries             []TryItem `json:"tries"`
	LastSuccessfulRun TryItem   `json:"lastSuccessfulRun"`
}

// TryItem ...
type TryItem struct {
	Steps  []TryStepItem `json:"steps"`
	Output string        `json:"output"`
}

// TryStepItem ...
type TryStepItem struct {
	Source string `json:"source"`
	URL    string `json:"url"`
	Image  string `json:"image"`
}
