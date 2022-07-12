package mqr

//go:generate stringer -type=State

type State int

const (
	Unacked State = iota
	Acked
	Rejected
	Pushed
)
