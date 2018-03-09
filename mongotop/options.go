package mongotop

var Usage = `<options> <polling interval in seconds>

Monitor basic usage statistics for each collection.

See http://docs.mongodb.org/manual/reference/program/mongotop/ for more information.`

// Output defines the set of options to use in displaying data from the server.
type Output struct {
	Locks    bool
	RowCount int32
	Json     bool
}

// Name returns a human-readable group name for output options.
func (_ *Output) Name() string {
	return "output"
}
