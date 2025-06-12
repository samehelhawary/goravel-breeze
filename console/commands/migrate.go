package commands

import (
	"fmt"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type Migrate struct {
}

func (receiver *Migrate) Extend() command.Extend {
	return command.Extend{}
}

// Signature the name and signature of hte console command.
func (receiver *Migrate) Signature() string {
	return "breeze:migrate"
}

// Description the console command description.
func (receiver *Migrate) Description() string {
	return "Migrate breeze tables"
}

// Handle Execute the console command.
func (receiver *Migrate) Handle(ctx console.Context) error {
	ctx.Info("Executing command: breeze:migrate")
	ctx.Info("---------------------------------")

	ctx.Info("Migrate database...")
	if err := facades.Artisan().Call("migrate"); err != nil {
		ctx.Error(fmt.Sprintf("Error migrating database: %v", err))
		return err
	}
	ctx.Info("Database migrated successfully.")

	ctx.Info("---------------------------------")
	ctx.Info("Breeze migration completed successfully.")

	return nil
}
