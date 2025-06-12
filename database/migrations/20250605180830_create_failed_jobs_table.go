package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605180830CreateFailedJobsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605180830CreateFailedJobsTable) Signature() string {
	return "20250605180830_create_failed_jobs_table"
}

// Up Run the migrations.
func (r *M20250605180830CreateFailedJobsTable) Up() error {
	if !facades.Schema().HasTable("failed_jobs") {
		return facades.Schema().Create("failed_jobs", func(table schema.Blueprint) {
			table.ID()
			table.String("uuid")
			table.Unique("uuid")
			table.Text("connection")
			table.Text("queue")
			table.LongText("payload")
			table.LongText("exception")
			table.Timestamp("failed_at").UseCurrent()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605180830CreateFailedJobsTable) Down() error {
	return facades.Schema().DropIfExists("failed_jobs")
}
