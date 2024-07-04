package httpMessage

type ErrorMessage string

const (
	ERROR_WRONG                  ErrorMessage = "Something went wrong"
	ERROR_ACCOUNT_EXIST          ErrorMessage = "Account exist"
	ERROR_ACCOUNT_WRONG_PASSWORD              = "Wrong password"
	ERROR_ACCOUNT_NOT_FOUND                   = "Account not found"
	ERROR_MISSING_INVALID_TOKEN               = "Missing or invalid token"
	ERROR_EXPIRED_TOKEN                       = "Token expired"
	ERROR_PERMISSION_KEY_EXIST                = "Permission key exist"
	ERROR_UNAUTHORIZED                        = "Unauthorized"
)
