package status

// Status ...
type Status int

// Statuses
const (
	Pending    Status = 1
	Processing Status = 2
	Success    Status = 3
	Failed     Status = 4
	Incomplete Status = 5
)

func (s Status) String() string {
	statuses := map[Status]string{
		1: "Pending",
		2: "Processing",
		3: "Success",
		4: "Failed",
		5: "Incomplete"}

	return statuses[s]
}
