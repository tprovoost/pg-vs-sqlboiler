package modules

import (
	"context"
	"fmt"

	bmodels "orm_compare/modules/bmodels"
	models "orm_compare/modules/shared"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func wrapCall(ctx context.Context, exec boil.ContextExecutor, fn func(context.Context, boil.ContextExecutor) error) {
	if err := fn(ctx, exec); err != nil {
		fmt.Printf("error while running function %T: %v\n", fn, err)
	}
}

// RunSQLBoiler runs a few methods with the SQL Boiler module
func RunSQLBoiler() {
	fmt.Println("-----------")
	fmt.Println("Starts SQL Boiler")
	ctx := context.Background()

	// Open handle to database like normal
	db, err := sqlx.Connect("postgres", "dbname=pgguide user=thomasprovoost host=localhost sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	wrapCall(ctx, db, boilerRead)
	wrapCall(ctx, db, boilerFetchIn)
	wrapCall(ctx, db, boilerInsert)
	wrapCall(ctx, db, boilerComplexQuery)
	wrapCall(ctx, db, boilerSqlxQuery) // to show the compatibility
}

func test() error {
	fmt.Println("Hallo")
	return nil
}

func boilerRead(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("----------\nRead")
	cpt, err := bmodels.Products().Count(ctx, exec)

	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Amount of products: %d\n", cpt)

	secondProduct, err := bmodels.FindProduct(ctx, exec, 2)
	fmt.Printf("Also, this is the product with id 2:\n%v\n", secondProduct)

	return nil
}

func boilerFetchIn(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("----------\nFetches items with known product IDs.")
	productIds := []int{1, 5, 10}
	products, err := bmodels.Products(bmodels.ProductWhere.ID.IN(productIds)).All(ctx, exec)

	if err != nil {
		return fmt.Errorf("error while fetching multiple products %v", err)
	}

	for i := 0; i < len(products); i++ {
		product := products[i]
		fmt.Printf("%v\n", product)
	}

	return nil
}

func boilerInsert(ctx context.Context, exec boil.ContextExecutor) error {
	// var product bmodels.Product

	// product.Insert(ctx, exec, boil.Infer())
	return nil
}

// Fetch product Ids and their quantity
func boilerComplexQuery(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("----------\nMake a complex query with Boiler QueryMod:")
	productsAndQuantities := make([]models.QuantityPerProduct, 0)

	rows, err := bmodels.NewQuery(
		qm.Select("product_id, sum(quantity)"),
		qm.From("purchase_items"),
		qm.GroupBy("product_id"),
		qm.OrderBy("product_id"),
	).QueryContext(ctx, exec)

	for rows.Next() {
		var qpp models.QuantityPerProduct
		if err = rows.Scan(&qpp.ProductID, &qpp.Quantity); err != nil {
			return fmt.Errorf("error while fetching quantity join %v", err)
		}

		productsAndQuantities = append(productsAndQuantities, qpp)
	}

	sum := 0
	for _, qpp := range productsAndQuantities {
		sum += qpp.Quantity
	}
	fmt.Printf("Average quantity per product is: %d\n", sum/len(productsAndQuantities))

	return nil
}

func boilerSqlxQuery(ctx context.Context, exec boil.ContextExecutor) error {
	// It is actually easier with SQLX as the binding can be done directly to an array of a struct.
	db := exec.(*sqlx.DB)
	fmt.Println("Or use the SQLX equivalent:")

	productsAndQuantities := []models.QuantityPerProduct{}

	err := db.Select(&productsAndQuantities, "SELECT product_id, SUM(quantity) quantity FROM purchase_items GROUP BY product_id ORDER BY product_id")

	if err != nil {
		return fmt.Errorf("error while fetching results %v", err)
	}

	sum := 0
	for _, qpp := range productsAndQuantities {
		sum += qpp.Quantity
	}
	fmt.Printf("Average per product: %d\n", sum/len(productsAndQuantities))
	return nil
}

func wrapString(s string) null.String {
	return null.NewString(s, false)
}

func wrapDecimal(d float64) null.Float64 {
	return null.NewFloat64(d, false)
}
