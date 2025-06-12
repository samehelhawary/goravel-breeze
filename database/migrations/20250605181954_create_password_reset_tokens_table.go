package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605181954CreatePasswordResetTokensTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605181954CreatePasswordResetTokensTable) Signature() string {
	return "20250605181954_create_password_reset_tokens_table"
}

// Up Run the migrations.
func (r *M20250605181954CreatePasswordResetTokensTable) Up() error {
	if !facades.Schema().HasTable("password_reset_tokens") {
		return facades.Schema().Create("password_reset_tokens", func(table schema.Blueprint) {
			table.String("email")
			table.Primary("email")
			table.String("token")
			table.Timestamp("created_at").Nullable()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605181954CreatePasswordResetTokensTable) Down() error {
	return facades.Schema().DropIfExists("password_reset_tokens")
}
