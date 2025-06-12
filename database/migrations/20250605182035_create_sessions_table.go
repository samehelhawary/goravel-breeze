package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605182035CreateSessionsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605182035CreateSessionsTable) Signature() string {
	return "20250605182035_create_sessions_table"
}

// Up Run the migrations.
func (r *M20250605182035CreateSessionsTable) Up() error {
	if !facades.Schema().HasTable("sessions") {
		return facades.Schema().Create("sessions", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id").Nullable()
			table.Foreign("user_id").References("id").On("users")
			table.String("ip_address", 45).Nullable()
			table.LongText("payload")
			table.Integer("last_activity")
			table.Index("last_activity")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605182035CreateSessionsTable) Down() error {
	return facades.Schema().DropIfExists("sessions")
}
