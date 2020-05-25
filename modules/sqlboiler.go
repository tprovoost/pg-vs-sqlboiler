package modules

import (
	"context"
	"fmt"
	dbmodels "orm_compare/database_models"
	models "orm_compare/models"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// RunSQLBoiler runs a few methods with the SQL Boiler module
func RunSQLBoiler() error {
	fmt.Println("Starts SQL Boiler")
	ctx := context.Background()

	// Open handle to database like normal
	db, err := sqlx.Connect("postgres", "dbname=pgguide user=thomasprovoost host=localhost sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	boilerRead(ctx, db)
	err = complexQuery(ctx, db)
	if err != nil {
		return fmt.Errorf("error while running RunSQLBoiler: %v", err)
	}

	return nil
}

func boilerRead(ctx context.Context, exec boil.ContextExecutor) error {
	cpt, err := dbmodels.Products().Count(ctx, exec)

	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Println(cpt)

	secondProduct, err := dbmodels.FindProduct(ctx, exec, 2)
	fmt.Printf("Product with id 2 is: %v\n", secondProduct)

	return nil
}

func boilerInsert(ctx context.Context, exec boil.ContextExecutor) error {
	return nil
}

// Fetch product Ids and their quantity
func complexQuery(ctx context.Context, db *sqlx.DB) error {

	productsAndQuantities := make([]models.QuantityPerProduct, 0)

	fmt.Println("Normal SQL still usable for  more complex queries")

	rows, err := dbmodels.NewQuery(
		qm.Select("product_id, sum(quantity)"),
		qm.From("purchase_items"),
		qm.GroupBy("product_id"),
		qm.OrderBy("product_id"),
	).Query(db)

	for rows.Next() {
		var qpp models.QuantityPerProduct
		if err = rows.Scan(&qpp.ProductID, &qpp.Quantity); err != nil {
			return fmt.Errorf("error while fetching quantity join %v", err)
		}

		productsAndQuantities = append(productsAndQuantities, qpp)
	}

	for i := 0; i < len(productsAndQuantities); i++ {
		fmt.Printf("%v\n", productsAndQuantities[i])
	}

	// It is actually easier with SQLX as the binding can be done directly to an array of a struct.

	fmt.Println("Or use the SQLX equivalent")

	productsAndQuantities2 := []models.QuantityPerProduct{}

	err = db.Select(&productsAndQuantities2, "SELECT product_id, SUM(quantity) quantity FROM purchase_items GROUP BY product_id ORDER BY product_id")

	if err != nil {
		return fmt.Errorf("error while fetching results %v", err)
	}

	for i := 0; i < len(productsAndQuantities2); i++ {
		fmt.Printf("%v\n", productsAndQuantities2[i])
	}

	return nil
}
