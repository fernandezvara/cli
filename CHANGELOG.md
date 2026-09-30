# Changelog

## v0.2.0

Made for programs that must control their exit code and be tested in-process.

### Added

- `cli.Exit(code, err)`, `cli.ExitCode(err)`, `cli.IsReported(err)` and `ExitOK`/`ExitFailure`/`ExitUsage`.
- `Config.SetEnv`, `Config.SetIO` and `Config.SetName`; `CommandContext.Getenv`, `Stdin`, `Stdout`, `Stderr`.
- `CommandContext.Positional()`: the arguments left after a command's flags.
- `help` as a word (`app help`, `app help <command>`, `app <command> help`).
- Real boolean flags: `--flag` and `--flag=false`.

### Changed

- **`Execute` never calls `os.Exit`.** Usage errors (unknown command or flag, missing or invalid flag value) return exit code 2 instead of exiting with 1; command errors are 1. `CommandResult.Handle` was removed.
- **Boolean flags no longer take the next argument**: `--verbose false` is now `--verbose` plus a leftover argument; write `--verbose=false`.
- Usage lines name the program and the full command: `Usage: app user list [options]`.
- A word after a group command that is not one of its subcommands is a usage error instead of showing help.
- Help is written to the configured stdout (default `os.Stdout`), errors to the configured stderr.

### Fixed

- README: subcommands are declared with `SubCommand`, not `Config(... cc.Command ...)`.
