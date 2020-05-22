package modules

import (
	"context"
	"database/sql"
	"fmt"
	models "orm_compare/database_models"
)

func RunSqlBoiler() error {
	ctx := context.Background()

	// Open handle to database like normal
	db, err := sql.Open("postgres", "dbname=pgguide user=thomasprovoost host=localhost sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	read(ctx, db)

	return nil
}

func read(ctx context.Context, db *sql.DB) error {
	cpt, err := models.Products().Count(ctx, db)

	if err != nil {
		return fmt.Errorf("Error while prompting count %v\n", err)
	}
	fmt.Println(cpt)

	return nil
}
