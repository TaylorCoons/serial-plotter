package serial

type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}
