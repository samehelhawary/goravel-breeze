package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240915060148CreateUsersTable struct {
}

// Signature The unique signature for the migration.
func (r *M20240915060148CreateUsersTable) Signature() string {
	return "20240915060148_create_users_table"
}

// Up Run the migrations.
func (r *M20240915060148CreateUsersTable) Up() error {
	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name")
		table.String("email")
		table.Unique("email")
		table.Timestamp("email_verified_at").Nullable()
		table.String("password")
		table.String("remember_token", 100).Nullable()
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20240915060148CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
