package model

type Result[T model] struct {
	Data   *T    `json:"data,omitempty"`
	Status int   `json:"status,omitempty"`
	Error  error `json:"error,omitempty"`
}

type model interface {
	Post | User | []User | []Post
}
