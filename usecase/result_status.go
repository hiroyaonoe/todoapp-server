package usecase

type ResultStatus struct {
	StatusCode int
	StatusMessage string
}

func NewResultStatus(statusCode int, msg string) *ResultStatus {
	return &ResultStatus{
		StatusCode: statusCode,
		StatusMessage: msg,
	}
}


