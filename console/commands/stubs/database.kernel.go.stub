package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"goravel/database/migrations"
	"goravel/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20240915060148CreateUsersTable{},
		&migrations.M20250605170712CreateCacheTable{},
		&migrations.M20250605170745CreateJobsTable{},
		&migrations.M20250605180614CreateCacheLocksTable{},
		&migrations.M20250605180755CreateJobBatchesTable{},
		&migrations.M20250605180830CreateFailedJobsTable{},
		&migrations.M20250605181954CreatePasswordResetTokensTable{},
		&migrations.M20250605182035CreateSessionsTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
