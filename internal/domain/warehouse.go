package domain

type Warehouse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"adress"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

type WarehouseReport struct {
	WarehouseName string `json:"warehouse_name"`
	ProductCount  int    `json:"product_count"`
}
