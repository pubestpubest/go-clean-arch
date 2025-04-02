package response

type Product struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
}
