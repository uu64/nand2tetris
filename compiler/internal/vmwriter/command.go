package vmwriter

//go:generate go run golang.org/x/tools/cmd/stringer -type=CommandType -linecomment
type CommandType int

const (
	Add CommandType = iota // add
	Sub                    // sub
	Neg                    // neg
	EQ                     // eq
	GT                     // gt
	LT                     // lt
	And                    // and
	Or                     // or
	Not                    // not
)
