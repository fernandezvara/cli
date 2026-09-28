// cli/command_result.go
package cli

import (
	"fmt"
	"os"
)

// CommandResult represents the result of command execution with unified error handling
type CommandResult struct {
	Error      error
	ExitCode   int
	ShouldExit bool
	Message    string
}

// Handle processes the command result according to its state
func (r *CommandResult) Handle() {
	if r.ShouldExit {
		r.displayAndExit()
	}
}

// displayAndExit displays the result message and exits with the appropriate code
func (r *CommandResult) displayAndExit() {
	if r.Message != "" {
		fmt.Fprintln(os.Stderr, r.Message)
	}
	os.Exit(r.ExitCode)
}

// success creates a successful command result
func success() *CommandResult {
	return &CommandResult{
		ExitCode:   0,
		ShouldExit: false,
	}
}

// errorResult creates an error command result
func errorResult(err error) *CommandResult {
	return &CommandResult{
		Error:      err,
		ExitCode:   1,
		ShouldExit: false,
	}
}

// validationError creates a validation error result that should exit
func validationError(message string) *CommandResult {
	return &CommandResult{
		Error:      fmt.Errorf("validation error: %s", message),
		ExitCode:   1,
		ShouldExit: true,
		Message:    message,
	}
}

// configErrorResult creates a configuration error result that should exit
func configErrorResult(message string) *CommandResult {
	return &CommandResult{
		Error:      fmt.Errorf("configuration error: %s", message),
		ExitCode:   1,
		ShouldExit: true,
		Message:    message,
	}
}
