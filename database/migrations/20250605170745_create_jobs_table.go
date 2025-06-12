package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605170745CreateJobsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605170745CreateJobsTable) Signature() string {
	return "20250605170745_create_jobs_table"
}

// Up Run the migrations.
func (r *M20250605170745CreateJobsTable) Up() error {
	if !facades.Schema().HasTable("jobs") {
		return facades.Schema().Create("jobs", func(table schema.Blueprint) {
			table.ID()
			table.String("queue")
			table.Index("queue")
			table.UnsignedTinyInteger("attempts")
			table.UnsignedInteger("reserved_at").Nullable()
			table.UnsignedInteger("available_at")
			table.UnsignedInteger("created_at")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605170745CreateJobsTable) Down() error {
	return facades.Schema().DropIfExists("jobs")
}
