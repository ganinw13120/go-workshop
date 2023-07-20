package apperror

type AppErrorCode string

const (
	NotFound         AppErrorCode = "NOT_FOUND"
	MySQLSyntaxError AppErrorCode = "MYSQL_SYNTAX_ERROR"
)

type AppError struct {
	Message     string
	Description string
	Code        AppErrorCode
}

func (err AppError) Error() string {
	return err.Message
}

func NewError(message, description string, code AppErrorCode) AppError {
	return AppError{
		Message:     message,
		Description: description,
		Code:        code,
	}
}
