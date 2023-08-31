package api

type ResponseStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() ResponseStatus {
	return ResponseStatus{
		Status: StatusOK,
	}
}

func Error(msg string) ResponseStatus {
	return ResponseStatus{
		Status: StatusError,
		Error:  msg,
	}
}
