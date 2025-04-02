package response

type Shop struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}
