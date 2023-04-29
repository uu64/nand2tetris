// Code generated by "stringer -type=CommandType -linecomment"; DO NOT EDIT.

package vmwriter

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Add-0]
	_ = x[Sub-1]
	_ = x[Neg-2]
	_ = x[EQ-3]
	_ = x[GT-4]
	_ = x[LT-5]
	_ = x[And-6]
	_ = x[Or-7]
	_ = x[Not-8]
}

const _CommandType_name = "addsubnegeqgtltandornot"

var _CommandType_index = [...]uint8{0, 3, 6, 9, 11, 13, 15, 18, 20, 23}

func (i CommandType) String() string {
	if i < 0 || i >= CommandType(len(_CommandType_index)-1) {
		return "CommandType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommandType_name[_CommandType_index[i]:_CommandType_index[i+1]]
}