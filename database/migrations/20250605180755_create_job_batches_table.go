package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250605180755CreateJobBatchesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250605180755CreateJobBatchesTable) Signature() string {
	return "20250605180755_create_job_batches_table"
}

// Up Run the migrations.
func (r *M20250605180755CreateJobBatchesTable) Up() error {
	if !facades.Schema().HasTable("job_batches") {
		return facades.Schema().Create("job_batches", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.Integer("total_jobs")
			table.Integer("pending_jobs")
			table.Integer("failed_jobs")
			table.LongText("failed_job_ids")
			table.MediumText("options").Nullable()
			table.Integer("cancelled_at").Nullable()
			table.Integer("created_at")
			table.Integer("finished_at").Nullable()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250605180755CreateJobBatchesTable) Down() error {
	return facades.Schema().DropIfExists("job_batches")
}
