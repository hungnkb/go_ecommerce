package httpMessage

type ErrorMessage string

const (
	ERROR_WRONG         ErrorMessage = "Something went wrong"
	ERROR_ACCOUNT_EXIST ErrorMessage = "Account exist"
)
