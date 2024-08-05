package types

type AllForMain struct {
	Products   []Product
	Categories []Category
}
type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
type Category struct {
	Id   int
	Name string
}