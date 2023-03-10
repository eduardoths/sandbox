package main

type SaveRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Data  interface{}     `json:"data,omitempty"`
	Error []ErrorResponse `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
