package parsing

type ManifestError struct {
	Message string
	Line    int
	Column  int
}

func (e ManifestError) Error() string {
	return e.Message
}
