package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605180614CreateCacheLocksTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605180614CreateCacheLocksTable) Signature() string {
	return "20250605180614_create_cache_locks_table"
}

// Up Run the migrations.
func (r *M20250605180614CreateCacheLocksTable) Up() error {
	if !facades.Schema().HasTable("cache_locks") {
		return facades.Schema().Create("cache_locks", func(table schema.Blueprint) {
			table.String("key")
			table.Primary("key")
			table.String("Owner")
			table.Integer("expiration")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605180614CreateCacheLocksTable) Down() error {
	return facades.Schema().DropIfExists("cache_locks")
}
