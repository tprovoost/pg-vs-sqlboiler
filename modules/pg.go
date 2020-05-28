package modules

import (
	"fmt"
	"os"

	"github.com/tprovoost/pg-vs-sqlboiler/modules/pgmodels"
	shared "github.com/tprovoost/pg-vs-sqlboiler/modules/shared"

	"github.com/go-pg/pg/v10"
)

// PGRunBenchmark executes N times the fn function
// and returns a Benchmark.
func PGRunBenchmark(fn func(*pg.DB) error, N int) shared.Benchmark {
	benchmark := shared.Benchmark{N: N}

	db := pg.Connect(&pg.Options{
		User:      "thomasprovoost",
		Database:  "pgguide",
		TLSConfig: nil,
	})
	defer db.Close()

	benchmark.StartTimer()
	for i := 0; i < benchmark.N; i++ {
		if err := fn(db); err != nil {
			benchmark.Failed = true
			os.Exit(2)
		}
	}
	benchmark.StopTimer()

	return benchmark
}

func pgCleanUp(db *pg.DB) error {
	minID := 20

	res, err := db.Model((*pgmodels.ProductPG)(nil)).Where("id>?", minID).Delete()
	if err != nil {
		return fmt.Errorf("could not delete products - %v", err)
	}

	fmt.Printf("Sucessfully deleted %d rows\n", res.RowsAffected())

	return nil
}

// PGReadOne reads one element in the database.
func PGReadOne(db *pg.DB) error {
	product := &pgmodels.ProductPG{ID: 2}
	err := db.Select(product)
	if err != nil {
		return fmt.Errorf("error while selecting product %v", err)
	}

	if shared.DebugMode {
		fmt.Printf("Product with id 2 is: %+v\n", product)
	}

	return nil
}

// PGFetchIn fetches products with specific IDs.
func PGFetchIn(db *pg.DB) error {
	productIds := []int{1, 5, 10}
	if shared.DebugMode {
		fmt.Printf("Fetches items with known product IDs: %v\n", productIds)
	}

	var products []pgmodels.ProductPG
	if err := db.Model(&products).WhereIn("id IN (?)", productIds).Select(); err != nil {
		return fmt.Errorf("Error while fetching multiple products: %v", err)
	}

	if shared.DebugMode {
		for _, p := range products {
			fmt.Println(p)
		}
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

func PGReadAll(db *pg.DB) error {
	fmt.Println("-----------")
	fmt.Println("Simply count products")
	var products []pgmodels.ProductPG

	cpt, err := db.Model(&products).Count()
	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Products in database: %d\n", cpt)

	return nil
}

func PGComplexQuery(db *pg.DB) error {
	fmt.Println("-----------")
	fmt.Println("Complex query")
	var qpps []shared.QuantityPerProduct
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

// PGInsert inserts one product into the database
func PGInsert(db *pg.DB) error {
	var product pgmodels.ProductPG
	title := "Smartphone"

	product.Title = title
	product.Price = 123.45
	product.Tags = []string{"Technology"}

	if err := db.Insert(&product); err != nil {
		return fmt.Errorf("Could not insert product %s. reason: %v", title, err)
	}

	if shared.DebugMode {
		fmt.Println("Inserted... now read")

		var p pgmodels.ProductPG

		if err := db.Model(&p).Where("title=?", title).Select(); err != nil {
			return fmt.Errorf("Could not read product %s. reason: %v", title, err)
		}

		fmt.Println(p)
	}

	return nil
}
