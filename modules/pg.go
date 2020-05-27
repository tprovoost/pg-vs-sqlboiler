package modules

import (
	"fmt"
	"time"

	"github.com/tprovoost/pg-vs-sqlboiler/modules/pgmodels"
	models "github.com/tprovoost/pg-vs-sqlboiler/modules/shared"

	"github.com/go-pg/pg/v10"
)

func wrapPGFunction(db *pg.DB, fn func(*pg.DB) error) {
	startTime := time.Now()
	if err := fn(db); err != nil {
		fmt.Printf("error while running pg function %v\n", err)
	}
	fmt.Printf("<%d ns>\n", time.Now().Sub(startTime))
}

// RunPG executes all PG commands
func RunPG() error {
	fmt.Println("-----------")
	fmt.Println("Starts PG")

	db := pg.Connect(&pg.Options{
		User:      "thomasprovoost",
		Database:  "pgguide",
		TLSConfig: nil,
	})
	defer db.Close()

	wrapPGFunction(db, pgReadOneProduct)
	wrapPGFunction(db, pgReadOnePurchaseItem)
	wrapPGFunction(db, pgFetchIn)
	wrapPGFunction(db, pgReadAll)
	wrapPGFunction(db, pgComplexQuery)

	return nil
}

func pgReadOneProduct(db *pg.DB) error {
	fmt.Println("-----------")
	fmt.Println("Read one element")

	product := &pgmodels.ProductPG{ID: 2}
	err := db.Select(product)
	if err != nil {
		return fmt.Errorf("error while selecting product %v", err)
	}

	fmt.Printf("Product with id 2 is: %+v\n", product)

	return nil
}

func pgFetchIn(db *pg.DB) error {
	fmt.Println("-----------")
	productIds := []int{1, 5, 10}
	fmt.Printf("Fetches items with known product IDs: %v\n", productIds)

	var products []pgmodels.ProductPG
	if err := db.Model(&products).WhereIn("id IN (?)", productIds).Select(); err != nil {
		return fmt.Errorf("Error while fetching multiple products: %v", err)
	}

	for _, p := range products {
		fmt.Println(p)
	}

	return nil
}

func pgReadOnePurchaseItem(db *pg.DB) error {
	fmt.Println("-----------")
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
	fmt.Println("-----------")
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
	fmt.Println("-----------")
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
