package responseType

type StorageReponseType struct {
	Data           interface{} `json:"data"`
	Error          string      `json:"error"`
	HttpStatusCode int         `json:"httpStatusCode"`
}
