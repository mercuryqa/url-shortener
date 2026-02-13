package errs

var (
	// ErrInternal - внутренняя ошибка сервера
	ErrInternal = NewError(500, "Internal error")

	// ErrNotFound - ресурс не найден
	ErrNotFound = NewError(404, "Not found")

	// ErrBadRequest - ошибка клиента (некорректный запрос)
	ErrBadRequest = NewError(400, "Bad request")

	// ErrDatabaseError - ошибка работы с базой данных
	ErrDatabaseError = NewError(500, "Database error")
)

type Error struct {
	HTTPCode int
	message  string
}

func (e *Error) Error() string {
	return e.message
}

func NewError(httpCode int, message string) *Error {
	return &Error{
		HTTPCode: httpCode,
		message:  message,
	}
}
