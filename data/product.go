package data

type Product struct {
	Id          int64
	Name        string
	Description string
}

var Products = map[int64]Product{
	1: {
		Id:          1,
		Name:        "product 1",
		Description: "description 1",
	},
	2: {
		Id:          2,
		Name:        "product 2",
		Description: "description 2",
	},
	3: {
		Id:          3,
		Name:        "product 3",
		Description: "description 3",
	},
}
