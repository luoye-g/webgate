package errors

const (
	PermissionError = "permission_error"
	NotFoundError   = "not_found_error"
	InternalError   = "internal_error"
	ValidationError = "validation_error"
)

func NewError(code string, err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		code:    code,
		message: err.Error(),
	}
}

func NewCustomError(code, message string) *AppError {
	return &AppError{
		code:    code,
		message: message,
	}
}

type AppError struct {
	code    string
	message string
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Code() string {
	return e.code
}
