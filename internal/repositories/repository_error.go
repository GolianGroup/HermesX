package repositories

import "errors"

type RepositoryErr struct {
	Err error
	Msg string
}

func (r *RepositoryErr) Error() string {
	return r.Err.Error()
}

func (r *RepositoryErr) Message() string {
	return r.Msg
}

func (r *RepositoryErr) Unwrap() error {
	return r.Err
}

var (
	ErrTimeout = &RepositoryErr{
		Msg: "Context timeout occured",
		Err: errors.New("context timeout")}
	ErrProfileNotFound = &RepositoryErr{
		Msg: "Profile with this credentials not found",
		Err: errors.New("record not found"),
	}
)
