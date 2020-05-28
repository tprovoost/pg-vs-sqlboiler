package modules

import (
	"context"
	"fmt"
	"os"

	bmodels "github.com/tprovoost/pg-vs-sqlboiler/modules/bmodels"
	shared "github.com/tprovoost/pg-vs-sqlboiler/modules/shared"

	"github.com/ericlagergren/decimal"
	"github.com/jmoiron/sqlx"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

// BoilerRunBenchmark executes N times the function read
// and updates result
func BoilerRunBenchmark(fn func(context.Context, boil.ContextExecutor) error, N int) shared.Benchmark {
	benchmark := shared.Benchmark{N: N}

	ctx := context.Background()
	exec, err := sqlx.Connect("postgres", "dbname=pgguide user=thomasprovoost host=localhost sslmode=disable")
	if err != nil {
		fmt.Printf("Could not open DB: %v", err)
		os.Exit(2)
	}
	defer exec.Close()

	benchmark.StartTimer()
	for i := 0; i < benchmark.N; i++ {
		if err := fn(ctx, exec); err != nil {
			benchmark.Failed = true
			os.Exit(2)
		}
	}
	benchmark.StopTimer()

	return benchmark
}

// BoilerCount counts the amount of products in database
func BoilerCount(ctx context.Context, exec boil.ContextExecutor) error {
	cpt, err := bmodels.Products().Count(ctx, exec)

	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Amount of products: %d\n", cpt)

	return nil
}

// BoilerReadOne gets the product with id 2.
func BoilerReadOne(ctx context.Context, exec boil.ContextExecutor) error {
	secondProduct, err := bmodels.FindProduct(ctx, exec, 2)
	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}

	if shared.DebugMode {
		fmt.Printf("Product with id 2:\n%v\n", secondProduct)
	}
	return nil
}

// BoilerFetchIn gets the products with specific IDs.
func BoilerFetchIn(ctx context.Context, exec boil.ContextExecutor) error {
	productIds := []int{1, 5, 10}

	if shared.DebugMode {
		fmt.Printf("Fetches items with known product IDs: %v\n", productIds)
	}

	products, err := bmodels.Products(bmodels.ProductWhere.ID.IN(productIds)).All(ctx, exec)

	if err != nil {
		return fmt.Errorf("error while fetching multiple products %v", err)
	}

	if shared.DebugMode {
		for i := 0; i < len(products); i++ {
			product := products[i]
			fmt.Printf("%v\n", product)
		}
	}

	return nil
}

// BoilerCleanUp cleans up database from all inserted products.
func BoilerCleanUp(ctx context.Context, exec boil.ContextExecutor) error {
	count, err := bmodels.Products(qm.Where("id>?", 20)).DeleteAll(ctx, exec)

	if err != nil {
		return fmt.Errorf("could not clean: %v", err)
	}
	if shared.DebugMode {
		fmt.Printf("Sucessfully deleted %d rows\n", count)
	}

	return nil
}

// BoilerInsert creates a single product and inserts it into the database.
func BoilerInsert(ctx context.Context, exec boil.ContextExecutor) error {
	title := "Smartphone"
	var product bmodels.Product

	product.Title = wrapS(title)
	product.Price = wrapD(12345, 2)
	product.Tags = []string{"Technology"}

	if err := product.Insert(ctx, exec, boil.Infer()); err != nil {
		return fmt.Errorf("error while inserting %v", err)
	}

	if shared.DebugMode {
		fmt.Println("Product inserted, now reading...")

		product.Reload(ctx, exec)

		p, err := bmodels.Products(qm.Where("title=?", title)).One(ctx, exec)
		if err != nil {
			return fmt.Errorf("error while reading after inserting %v", err)
		}

		fmt.Println(p)
	}

	return nil
}

// BoilerComplexQuery fetch product Ids and their quantity
func BoilerComplexQuery(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("----------")
	fmt.Println("Make a complex query with Boiler QueryMod:")
	productsAndQuantities := make([]shared.QuantityPerProduct, 0)

	rows, err := bmodels.NewQuery(
		qm.Select("product_id, sum(quantity)"),
		qm.From("purchase_items"),
		qm.GroupBy("product_id"),
		qm.OrderBy("product_id"),
	).QueryContext(ctx, exec)

	for rows.Next() {
		var qpp shared.QuantityPerProduct
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

	productsAndQuantities := []shared.QuantityPerProduct{}

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

func wrapS(s string) null.String {
	return null.NewString(s, true)
}

func wrapF(d float64) null.Float64 {
	return null.NewFloat64(d, true)
}

func wrapD(value int64, scale int) types.NullDecimal {
	big := decimal.New(value, scale)
	return types.NewNullDecimal(big)
}
