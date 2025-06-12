package commands

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

// Group all embedded stubs together for clarity.
var (
	//go:embed stubs/database.kernel.go.stub
	kernelStub []byte
	//go:embed stubs/app.http.kernel.go.stub
	httpKernelStub []byte
	//go:embed stubs/app.http.controllers.auth.auth_controller.go.stub
	authControllerStub []byte
	//go:embed stubs/app.http.middleware.auth_functions.go.stub
	middlewareAuthFunctionsStub []byte
	//go:embed stubs/app.http.middleware.remember_me.go.stub
	rememberMeMiddlewareStub []byte
	//go:embed stubs/app.providers.validation_service_provider.go.stub
	validationServiceProviderStub []byte
	//go:embed stubs/routes.web.go.stub
	routesWebStub []byte
)

type Install struct {
}

func (receiver *Install) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *Install) Signature() string {
	return "breeze:install"
}

func (receiver *Install) Description() string {
	return "Install all of the Breeze resources"
}

// Handle orchestrates the entire installation process.
func (receiver *Install) Handle(ctx console.Context) error {
	ctx.Info("Executing command: breeze:install")
	ctx.Info("---------------------------------")

	// 1. Prepare Environment
	if err := runShellCommand(ctx, "go", "mod", "tidy"); err != nil {
		return err
	}
	if err := removeWelcomeView(ctx); err != nil {
		return err
	}
	if err := copyEnvFile(ctx); err != nil {
		return err
	}
	if err := facades.Artisan().Call("key:generate"); err != nil {
		ctx.Error(fmt.Sprintf("Error generating app key: %v", err))
		return err
	}

	// 2. Publish Vendor Assets (views, configs, etc.)
	ctx.Info("Publishing Breeze assets...")
	publishCmd := "vendor:publish --package=github.com/samehelhawary/goravel-breeze --force"
	if err := facades.Artisan().Call(publishCmd); err != nil {
		ctx.Error(fmt.Sprintf("Error publishing Breeze assets: %v", err))
		return err
	}
	ctx.Info("Assets published successfully.")

	// 3. Publish Stub Files Programmatically
	// A map centralizes all stubs and their destinations.
	stubsToPublish := map[string][]byte{
		"database/kernel.go":                           kernelStub,
		"app/http/kernel.go":                           httpKernelStub,
		"app/http/controllers/auth/auth_controller.go": authControllerStub,
		"app/http/middleware/auth_functions.go":        middlewareAuthFunctionsStub,
		"app/http/middleware/remember_me.go":           rememberMeMiddlewareStub,
		"app/providers/validation_service_provider.go": validationServiceProviderStub,
		"routes/web.go":                                routesWebStub,
	}

	// Loop through the map and publish each file using a single helper.
	for dest, content := range stubsToPublish {
		if err := publishStub(ctx, dest, content); err != nil {
			return err
		}
	}

	// 4. Finalize
	ctx.Info("Finalizing installation...")
	if err := runShellCommand(ctx, "go", "mod", "tidy"); err != nil {
		return err
	}

	ctx.Info("---------------------------------")
	ctx.Info("Breeze installation completed successfully.")
	return nil
}

// --- Helper Functions ---

// publishStub is a single, reusable function that replaces all the previous `publish...` functions.
func publishStub(ctx console.Context, destPath string, stubContent []byte) error {
	ctx.Info(fmt.Sprintf("Generating %s...", destPath))

	// Ensure the parent directory exists before writing the file.
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		ctx.Error(fmt.Sprintf("Error creating directory for %s: %v", destPath, err))
		return err
	}

	// Write the stub content to the destination file.
	if err := os.WriteFile(destPath, stubContent, 0644); err != nil {
		ctx.Error(fmt.Sprintf("Error writing %s: %v", destPath, err))
		return err
	}

	ctx.Info(fmt.Sprintf("Created %s successfully.", destPath))
	return nil
}

// removeWelcomeView safely removes the default welcome view.
func removeWelcomeView(ctx console.Context) error {
	path := "resources/views/welcome.tmpl"
	ctx.Info("Removing default welcome view...")
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		ctx.Error(fmt.Sprintf("Error removing '%s': %v", path, err))
		return err
	}
	if os.IsNotExist(err) {
		ctx.Warning(fmt.Sprintf("'%s' does not exist, skipping.", path))
	} else {
		ctx.Info("Welcome view removed.")
	}
	return nil
}

// copyEnvFile copies .env.example to .env if .env doesn't exist.
func copyEnvFile(ctx console.Context) error {
	sourceFile, destFile := ".env.example", ".env"

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
	defer source.Close()

	destination, err := os.Create(destFile)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error creating '%s': %v", destFile, err))
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		ctx.Error(fmt.Sprintf("Error copying contents to '%s': %v", destFile, err))
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
	err := cmd.Run()
	if err != nil {
		ctx.Error(fmt.Sprintf("Command failed: %v", err))
	}
	return err
}
