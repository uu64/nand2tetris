package codewriter

//go:generate go run golang.org/x/tools/cmd/stringer -type=SegmentType -linecomment
type SegmentType int

const (
	Const   SegmentType = iota // constant
	Arg                        // argument
	Local                      // local
	Static                     // static
	This                       // this
	That                       // that
	Pointer                    // pointer
	Temp                       // temp
)
