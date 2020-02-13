package comm

// TODO: A bit of polymorphism here to not send unnecessary args would be nice.
type Request struct {
	Type          string
	Command       string
	Args          []string
	JobId         int
	Start         int
	StartScale    string
	Continue      int
	ContinueScale string
}

// TODO: No longer used?
func NewRequest(args []string) *Request {
	return &Request{
		Command: args[1],
		Args:    args[2:],
	}
}
