package comm

type Request struct {
	Command string
	Args    []string
}

func NewRequest(args []string) *Request {
	return &Request{
		Command: args[1],
		Args:    args[2:],
	}
}
