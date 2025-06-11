package commands

import (
	"fmt"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/goravel/framework/contracts/console"
)

type Install struct {
}

func (receiver *Install) Extend() command.Extend {
	return command.Extend{}
}

// Signature the name and signature of hte console command.
func (receiver *Install) Signature() string {
	return "breeze:install"
}

// Description the console command description.
func (receiver *Install) Description() string {
	return "Install all of the Breeze resources"
}

// Handle Execute the console command.
func (receiver *Install) Handle(ctx console.Context) error {
	ctx.Info("Executing command: breeze:install")
	ctx.Info("---------------------------------")

	ctx.Info("Running go mod tidy...")
	if err := runShellCommand(ctx, "go", "mod", "tidy"); err != nil {
		ctx.Error(fmt.Sprintf("Error running go mod tidy: %v", err))
		return err
	}
	ctx.Info("'go mod tidy' finished.")

	// Step 1: remove welcome.tmpl
	if err := removeWelcomeView(ctx); err != nil {
		return err
	}

	// Step 2: Copy .env.example to .env
	if err := copyEnvFile(ctx); err != nil {
		return err
	}

	// Step 3: Generate app key
	ctx.Info("Generating app key...")
	if err := facades.Artisan().Call("key:generate"); err != nil {
		ctx.Error(fmt.Sprintf("Error generating app key: %v", err))
		return err
	}
	ctx.Info("App key generated.")

	// Step 4: Publish vendor assets
	ctx.Info("Publishing Breeze assets...")
	publishCmd := "vendor:publish --package=github.com/samehelhawary/goravel-breeze -f"
	if err := facades.Artisan().Call(publishCmd); err != nil {
		ctx.Error(fmt.Sprintf("Error publishing Breeze asset: %v", err))
		return err
	}
	ctx.Info("Assets published successfully.")

	// Step 5: Run go mod tidy
	ctx.Info("Running go mod tidy...")
	if err := runShellCommand(ctx, "go", "mod", "tidy"); err != nil {
		ctx.Error(fmt.Sprintf("Error running go mod tidy: %v", err))
		return err
	}
	ctx.Info("'go mod tidy' finished.")

	ctx.Info("---------------------------------")
	ctx.Info("Breeze installation completed successfully.")

	return nil
}

// removeWelcomeView safely removes the default welcome view.
func removeWelcomeView(ctx console.Context) error {
	path := "resources/views/welcome.tmpl"
	ctx.Info("Remove default welcome view...")
	err := os.Remove(path)
	if err != nil {
		// If the file doesn't exist, it's not an error in this context.
		if os.IsNotExist(err) {
			ctx.Warning(fmt.Sprintf("'%s' does not exist, skipping.", path))
			return nil
		}
		ctx.Error(fmt.Sprintf("Error removing '%s': %v", path, err))
	}
	ctx.Info("Welcome view removed.")
	return nil
}

// copyEnvFile copies .env.example to .env if .env doesn't exist.
func copyEnvFile(ctx console.Context) error {
	sourceFile := ".env.example"
	destFile := ".env"

	ctx.Info("Checking for .env file...")
	if _, err := os.Stat(destFile); err == nil {
		ctx.Warning(fmt.Sprintf("'%s' already exists, skipping.", destFile))
		return nil
	}

	ctx.Info("Creating .env file...")
	source, err := os.Open(sourceFile)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error opening '%s': %v", sourceFile, err))
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			ctx.Error(fmt.Sprintf("Error closing '%s': %v", sourceFile, err))
		}
	}(source)

	destination, err := os.Create(destFile)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error creating '%s': %v", destFile, err))
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			ctx.Error(fmt.Sprintf("Error closing '%s': %v", destFile, err))
		}
	}(destination)

	_, err = io.Copy(destination, source)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error copying '%s': %v", sourceFile, err))
		return err
	}

	ctx.Info(".env file created.")
	return nil
}

// runShellCommand executes an external shell command and streams its output.
func runShellCommand(ctx console.Context, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ctx.Info(fmt.Sprintf("↳ %s", strings.Join(cmd.Args, " ")))
	return cmd.Run()
}
