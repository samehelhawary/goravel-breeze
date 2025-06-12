package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605170712CreateCacheTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605170712CreateCacheTable) Signature() string {
	return "20250605170712_create_cache_table"
}

// Up Run the migrations.
func (r *M20250605170712CreateCacheTable) Up() error {
	if !facades.Schema().HasTable("cache") {
		return facades.Schema().Create("cache", func(table schema.Blueprint) {
			table.String("key")
			table.Primary("key")
			table.MediumText("value")
			table.Integer("expiration")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605170712CreateCacheTable) Down() error {
	return facades.Schema().DropIfExists("cache")
}
