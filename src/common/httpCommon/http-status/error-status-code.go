package httpStatusCode

type HttpStatusCode int

const (
	OK          HttpStatusCode = 200
	BAD_REQUEST HttpStatusCode = 400
	CONFLICT    HttpStatusCode = 409
)
