package sessions

// View is model for the view page
type View struct {
	ProjectName  string
	Title        string
	Cases        string
	SilentTries  bool
	SessionError string
}
