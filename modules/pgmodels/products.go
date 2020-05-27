package pgmodels

import (
	"time"

	"github.com/go-pg/pg/v10"
)

// ProductPG represents a product in database
type ProductPG struct {
	// if struct has different name than the table, it is necessary
	// to define the correct name.
	tableName struct{} `pg:"select:products"`
	ID        int64    `pg:",pk"`
	Title     string
	Price     float64
	CreatedAt time.Time
	DeletedAt pg.NullTime `pg:",soft_delete"`
	Tags      []string    `pg:",array"`
}
