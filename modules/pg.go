package modules

import (
	"fmt"
	"orm_compare/modules/pgmodels"

	models "orm_compare/modules/shared"

	"github.com/go-pg/pg/v10"
)

// RunPG executes all PG commands
func RunPG() error {
	fmt.Println("Starts PG")

	db := pg.Connect(&pg.Options{
		User:      "thomasprovoost",
		Database:  "pgguide",
		TLSConfig: nil,
	})
	defer db.Close()

	if err := pgReadOneProduct(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading value %v", err)
	}

	if err := pgReadOnePurchaseItem(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading value %v", err)
	}

	if err := pgReadAll(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading all %v", err)
	}

	if err := pgComplexQuery(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while doing complex query %v", err)
	}

	return nil
}

func pgReadOneProduct(db *pg.DB) error {
	fmt.Println("Read one element")

	product := &pgmodels.ProductPG{ID: 2}
	err := db.Select(product)
	if err != nil {
		return fmt.Errorf("error while selecting product %v", err)
	}

	fmt.Printf("Product with id 2 is: %+v\n", product)

	return nil
}

func pgReadOnePurchaseItem(db *pg.DB) error {
	fmt.Println("Read one element with wrong columns order")

	pItem := &pgmodels.PurchaseItemPG{ID: 2}
	err := db.Select(pItem)
	if err != nil {
		return fmt.Errorf("error while selecting product %v", err)
	}

	fmt.Printf("Purchase item is correct: %+v\n", pItem)

	return nil
}

func pgReadAll(db *pg.DB) error {
	fmt.Println("Read data")
	var products []pgmodels.ProductPG

	cpt, err := db.Model(&products).Count()
	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Products in database: %d\n", cpt)

	return nil
}

func pgComplexQuery(db *pg.DB) error {
	var qpps []models.QuantityPerProduct
	_, err := db.Query(&qpps, `SELECT product_id, SUM(quantity) quantity FROM purchase_items GROUP BY product_id ORDER BY product_id`)

	if err != nil {
		return fmt.Errorf("error while fetching results %v", err)
	}

	// Now calculate manually all quantities together
	sum := 0
	for _, qpp := range qpps {
		sum += qpp.Quantity
	}

	fmt.Printf("Amount of purchased items, ever: %d\n", sum)

	return nil
}
