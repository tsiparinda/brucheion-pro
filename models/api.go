package models


type User struct {
	Name string `json:"name"`
}

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}