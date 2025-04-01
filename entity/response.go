package entity

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponsePagination struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

type ResponseError struct {
	Error string `json:"error"`
}
