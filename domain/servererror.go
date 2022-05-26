package domain

import "github.com/google/uuid"

type ServerError struct {
	Message string `json:"message"`
	Ticket  string `json:"ticket"`
}

func GenerateServerError(msg string) ServerError {
	var serverError = ServerError{}

	serverError.Message = msg
	uuid, err := uuid.NewUUID()
	if err == nil { // notice this is the uuid generator error
		serverError.Ticket = uuid.String()
	}

	return serverError
}
