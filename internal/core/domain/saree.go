package domain

type Saree struct {
	Id               string     `json:"id"`
	FabricType       string     `json:"fabric_type" `
	Category         string     `json:"category"`
	Color            string     `json:"color"`
	DesignPattern    string     `json:"design_pattern"`
	SareeDimensions  Dimensions `json:"saree_dimensions"`
	BlouseDimensions Dimensions `json:"blouse_dimensions"`
	Price            int        `json:"price"`
	Availability     bool       `json:"availability"`
	Stock            int        `json:"stock"`
	ImagesUrl        []string   `json:"images_url"`
}

type Dimensions struct {
	Length  string `json:"length"`
	Breadth string `json:"breadth"`
}
