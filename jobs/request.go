package jobs

// TODO: A bit of polymorphism here to not send unnecessary args would be nice.
type Request struct {
	Username string
	Type     string

	// Most command argument for request types
	JobId int

	// Basic Command
	Cmd  string
	Args []string

	Start      int
	StartScale *TimeScale
	Continue   int
	Scale      *TimeScale

	// Mark down the user's time
	Time int64
}
