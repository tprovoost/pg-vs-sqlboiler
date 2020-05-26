package shared

// QuantityPerProduct is a structure thaat shows how many items were sold
// for a specific product
type QuantityPerProduct struct {
	ProductID int `db:"product_id"`
	Quantity  int `db:"quantity"`
}
