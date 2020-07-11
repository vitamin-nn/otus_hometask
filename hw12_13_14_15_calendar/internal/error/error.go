package error

type OutError string

func (e OutError) Error() string {
	return string(e)
}

var (
	ErrEventNotFound = OutError("event not found")
	ErrDateBusy      = OutError("time is busy")
	ErrUnknownUserID = OutError("Unknown user_id")
	ErrEmptyEvent    = OutError("Empty event param")
)
