package modal

type Saree struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Image byte   `json:"image"`
	Type string `json:"type"`
	Color string `json:"color"`
}
