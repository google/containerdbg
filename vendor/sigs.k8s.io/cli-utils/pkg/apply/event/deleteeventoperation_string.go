// Code generated by "stringer -type=DeleteEventOperation"; DO NOT EDIT.

package event

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DeleteUnspecified-0]
	_ = x[Deleted-1]
	_ = x[DeleteSkipped-2]
}

const _DeleteEventOperation_name = "DeleteUnspecifiedDeletedDeleteSkipped"

var _DeleteEventOperation_index = [...]uint8{0, 17, 24, 37}

func (i DeleteEventOperation) String() string {
	if i < 0 || i >= DeleteEventOperation(len(_DeleteEventOperation_index)-1) {
		return "DeleteEventOperation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DeleteEventOperation_name[_DeleteEventOperation_index[i]:_DeleteEventOperation_index[i+1]]
}
