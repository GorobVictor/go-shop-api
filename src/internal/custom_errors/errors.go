package customerrors

type BadRequestError struct {
	Message string `json:"message"`
}

func (e *BadRequestError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string `json:"message"`
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string `json:"message"`
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

type InternalServerError struct {
	Message string `json:"message"`
}

func NewInternalServerError() *InternalServerError {
	return &InternalServerError{Message: "Internal server error"}
}

func (e *InternalServerError) Error() string {
	return e.Message
}
