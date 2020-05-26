package pgmodels

// PurchaseItemPG represents a purchased items in database.
// Properties are defined in the wrong order to show
// that PG is able to detect with naming convention
// which one is correct.
//
// If name is too different though, it's always possible to
// define the proper column name with magic quotes.
type PurchaseItemPG struct {
	tableName  struct{} `pg:"select:purchase_items"`
	State      string
	ID         int64 `pg:",pk"`
	PurchaseID int64
	Price      float64
	ProID      int64 `pg:"product_id"`
	Quantity   int64
}
