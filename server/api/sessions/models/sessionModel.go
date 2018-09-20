package sessions

// SessionModel ...
type SessionModel struct {
	ID                  string
	Status              string
	SessionError        string
	CasesCount          int
	CasesProcessedCount int
	CasesFailed         []string
}
