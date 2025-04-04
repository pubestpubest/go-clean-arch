package entity

type ResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}
