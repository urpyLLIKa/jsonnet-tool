package exitcode

// Error is an error code which returns a exit code for the program.
// this approach helps with testing.
type Error struct {
	Msg      string
	ExitCode int
}

func (e *Error) Error() string {
	return e.Msg
}

// Failed creates an Error instance with exitCode=1.
func Failed() *Error {
	return &Error{ExitCode: 1, Msg: "failed"}
}

// Invalid creates an Error instance with exitCode=3.
func Invalid() *Error {
	return &Error{ExitCode: 3, Msg: "invalid"}
}
