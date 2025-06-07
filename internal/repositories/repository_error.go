package repositories

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
