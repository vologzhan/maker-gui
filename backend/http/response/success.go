package response

type Success struct {
	Success bool `json:"success"`
}

func NewSuccess() *Success {
	return &Success{true}
}
