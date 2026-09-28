# cli

A Go CLI framework that makes building command-line applications simple and enjoyable.

## Why cli?

cli eliminates the boilerplate and complexity of building CLI applications, letting you focus on your application logic. With type-safe configuration, automatic help generation, and clear error handling, you can create robust CLIs in minutes, not hours.

### What Makes cli Different

- **Zero Boilerplate**: Define once, use anywhere - no repetitive setup code
- **Type Safety**: Compile-time guarantees with generics - no more string-based configuration
- **Transparent by Default**: Clean output, silent overrides, and clear error messages
- **Production Features**: File configuration, secret protection, middleware, and comprehensive validation
- **Performance Optimized**: 68% faster template rendering, 79% fewer allocations

## Quick Start

### Configuration-Only Applications

Perfect for services, daemons, and tools that need configuration without commands:

```go
package main

import (
    "fmt"
    "os"
    "strings"
    
    "github.com/fernandezvara/cli"
)

func main() {
    cfg := cli.New()

    // Define your configuration - fluent and readable
    cfg.Define("PORT").
        Int64().
        Env("PORT").
        Flag("port").
        Default(8080).
        Range(1, 65535).
        Description("HTTP server port")

    cfg.Define("DATABASE_URL").
        String().
        Env("DATABASE_URL").
        Required().
        Secret().
        Description("Database connection string")

    // Add empty string command for config-only mode
    cfg.Command("").
        Func(func(ctx *cli.CommandContext) error {
            // Type-safe access to configuration
            port := cli.MustGet[int64](ctx, "PORT")
            dbURL := cli.MustGet[string](ctx, "DATABASE_URL")
            
            fmt.Printf("Server starting on port %d\n", port)
            fmt.Printf("Database: %s\n", maskSecret(dbURL))
            return nil
        }).
        ShortHelp("Start the server").
        LongHelp("Starts the web server with the specified configuration.").
        Config(func(cc *cli.CommandConfig) {
            // Add configuration to the default command
            cc.Define("PORT").
                Int64().
                Env("PORT").
                Flag("port").
                Default(8080).
                Range(1, 65535).
                Description("HTTP server port")

            cc.Define("DATABASE_URL").
                String().
                Env("DATABASE_URL").
                Required().
                Secret().
                Description("Database connection string")
        })

    // One line to process everything
    if err := cfg.Execute(os.Args); err != nil {
        os.Exit(1)
    }
    defer cfg.Destroy()
}

func maskSecret(secret string) string {
    if len(secret) <= 8 {
        return strings.Repeat("*", len(secret))
    }
    return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}
```

**Usage:**
```bash
./app                           # Starts with defaults
./app --port 9000               # Starts on port 9000
./app --help                    # Shows help for default command
DATABASE_URL=... ./app          # Starts with environment variable
```

## Empty String Command - The Magic

cli introduces an elegant solution for configuration-only applications: the **empty string command** (`cfg.Command("")`). 

When you define an empty string command, it becomes the **default action** that executes when no command is provided. This creates a seamless experience for:

- **Services & Daemons** - Run directly with configuration flags
- **Simple Tools** - No need for subcommands, just configure and run  
- **Configuration Management** - Perfect for apps that just need to load config and start

**How it works:**
```bash
./app                    # Runs empty string command
./app --port 9000        # Empty string command gets the flag
./app --help             # Shows help for empty string command
```

The empty string command has access to all cli features:
- Type-safe configuration access
- Environment variable support  
- Flag parsing and validation
- Secret management
- Help generation
- Error handling

This approach eliminates the need for separate APIs or special cases - configuration-only apps use the exact same command system as complex CLIs!

### Command-Based Applications

Perfect for CLI tools with multiple commands and subcommands:

```go
package main

import (
    "fmt"
    "github.com/fernandezvara/cli"
)

func main() {
    cfg := cli.New()

    // Global configuration
    cfg.Define("VERBOSE").
        Bool().
        Flag("verbose").
        Default(false).
        Description("Enable verbose output")

    // Define commands with their own configuration
    cfg.Command("deploy").
        Func(deployCommand).
        ShortHelp("Deploy the application").
        LongHelp("Deploy the application to the specified environment.").
        Config(func(cc *cli.CommandConfig) {
            cc.Define("ENVIRONMENT").
                String().
                Flag("env").
                Required().
                OneOf("dev", "staging", "prod").
                Description("Target environment")
            
            cc.Define("DRY_RUN").
                Bool().
                Flag("dry-run").
                Default(false).
                Description("Show what would be deployed")
        })

    cfg.Command("status").
        Func(statusCommand).
        ShortHelp("Show application status").
        Aliases("st", "info")

    // Execute with professional help and error handling
    cfg.Execute(os.Args)
}

func deployCommand(ctx *cli.CommandContext) error {
    env := cli.MustGet[string](ctx, "ENVIRONMENT")
    dryRun := cli.MustGet[bool](ctx, "DRY_RUN")
    
    if dryRun {
        fmt.Printf("Would deploy to %s (dry run)\n", env)
    } else {
        fmt.Printf("Deploying to %s\n", env)
    }
    return nil
}

func statusCommand(ctx *cli.CommandContext) error {
    fmt.Println("Application is running")
    return nil
}
```

## Real-World Usage

### Clear Error Handling

cli provides clear, actionable error messages that help users fix problems:

```bash
$ go run app.go --port 99999
Usage: app [options]

Configuration errors:
  --port int64 (default: 8080) -> value 99999 is greater than maximum 65535

Flags:
  --port int64 (default: 8080) (valid: 1-65535)
        HTTP server port
```

### File Configuration

Load configuration from JSON, YAML, or TOML files with flexible key mapping:

```go
cfg.Define("PORT").
    Int64().
    Flag("port").
    File("server_port").  // Maps to "server_port" in files
    Default(8080)

cfg.LoadFile("config.json")  // Load once, use everywhere
```

**config.json:**
```json
{
  "server_port": 3000,
  "database_url": "postgres://localhost/mydb",
  "log_level": "debug"
}
```

### Silent Override Behavior

CLI tools don't warn about expected behavior:

```bash
# Environment variable (8080) -> Flag (3000) -> Works silently
PORT=8080 go run app.go --port 3000
Server starting on port 3000

# No confusing warning messages cluttering the output
```

### Secret Protection

Sensitive data gets special treatment:

```go
cfg.Define("API_KEY").
    String().
    Required().
    Secret().
    Description("API authentication key")

// Access safely
secret := cfg.GetSecret("API_KEY")
if secret.IsSet() {
    fmt.Printf("API key configured (%d bytes)\n", secret.Size())
    // Use secret.String() or secret.Bytes() when actually needed
}
```

## File Configuration

cli supports multiple file formats with flexible key mapping:

### Basic Usage

```go
cfg.Define("PORT").
    Int64().
    Flag("port").
    File("port_in_file").  // Look for this key in files
    Default(8080)

cfg.Define("DATABASE_URL").
    String().
    File("db_connection").
    Required().
    Secret()

// Load from environment variable containing file path
cfg.LoadFileFromEnv("CONFIG_FILE")

// Or load directly
cfg.LoadFile("config.json")
```

### Priority System

Configuration sources resolve in priority order. The default is:

```
Flag > Environment > File > Default
```

You can change the order globally or per definition:

```go
// Global default
cfg.SetDefaultPriority(cli.PriorityFileEnvFlagDefault)

// Per definition
cfg.Define("PORT").
    Int64().
    Flag("port").
    Env("PORT").
    Default(8080).
    Priority(cli.PriorityEnvFlagDefault)
```

Built-in presets: `PriorityFlagEnvFileDefault`, `PriorityFlagEnvDefault`, `PriorityEnvFlagDefault`, `PriorityFileEnvFlagDefault`, `PriorityDefaultOnly`. Priorities referencing sources a definition doesn't have are reported as configuration errors.

### Multiple Files

```go
// Load multiple files (later files override earlier ones)
cfg.LoadFiles("config.json", "secrets.json", "local.json")
```

## Configuration Types

### Supported Types

cli has seven core types:

```go
cfg.Define("NAME").String().Default("app")
cfg.Define("PORT").Int64().Default(8080)
cfg.Define("RATE").Float64().Default(100.0)
cfg.Define("ENABLED").Bool().Default(true)
cfg.Define("TIMEOUT").Duration().Default(30 * time.Second)
cfg.Define("TAGS").StringSlice().Default([]string{"v1", "api"})
cfg.Define("NUMBERS").Int64Slice().Default([]int64{1, 2, 3})
```

### Rich Validation

```go
cfg.Define("PORT").
    Int64().
    Range(1, 65535).                    // Numeric range
    Required()                          // Required field

cfg.Define("EMAIL").
    String().
    Regexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`). // Email format
    MinLength(5).                      // Minimum length
    MaxLength(100).                    // Maximum length

cfg.Define("LOG_LEVEL").
    String().
    OneOf("debug", "info", "warn", "error"). // Enum validation
    Default("info")
```

## Command System

### Commands with Configuration

```go
cfg.Command("deploy").
    Func(deployCommand).
    ShortHelp("Deploy the application").
    LongHelp("Deploy the application to the specified environment.").
    Config(func(cc *cli.CommandConfig) {
        cc.Define("ENVIRONMENT").
            String().
            Flag("env").
            Required().
            OneOf("dev", "staging", "prod").
            Description("Target environment")
    })
```

### Subcommands and Aliases

```go
cfg.Command("docker").
    ShortHelp("Docker operations").
    Config(func(cc *cli.CommandConfig) {
        cc.Command("run").
            Func(dockerRunCommand).
            ShortHelp("Run Docker container")
        
        cc.Command("stop").
            Func(dockerStopCommand).
            ShortHelp("Stop Docker container")
    })

cfg.Command("start").
    Func(startCommand).
    ShortHelp("Start the service").
    Aliases("run", "up")  // Multiple aliases
```

## Middleware System

Add cross-cutting concerns to your commands:

```go
// Global middleware - applies to all commands
cfg.UseMiddleware(cli.RecoveryMiddleware())
cfg.UseMiddleware(cli.LoggingMiddleware(func(ctx *cli.CommandContext, d time.Duration) {
    log.Printf("%s completed in %v", ctx.Command, d)
}))
cfg.UseMiddleware(cli.MetricsMiddleware(collectMetrics))

// Command-specific middleware
cfg.UseMiddlewareForCommands([]string{"admin", "shutdown"}, authMiddleware)

// Custom middleware - plain functions matching CommandMiddleware
// (see a full example in examples/cli-tool: tokenAuthMiddleware)
func authMiddleware(next cli.CommandFunc) cli.CommandFunc {
    return func(ctx *cli.CommandContext) error {
        // ...validate...
        return next(ctx)
    }
}
```

Built-in middleware:

| Middleware | Description |
| ------ | ----------- |
| `RecoveryMiddleware()` | Recover from panics in commands |
| `LoggingMiddleware(func(ctx, duration))` | Log command execution with timing |
| `MetricsMiddleware(func(ctx, duration, err))` | Collect command metrics |
| `AuthMiddleware(func(ctx) error)` | Validate authentication before execution |
| `TimingMiddleware()` | Measure and store execution timing in the context |
| `ConditionalMiddleware(cond, mw)` | Apply middleware only when a condition holds |

## Clear Help

cli automatically generates helpful help:

### Global Help
```bash
$ go run myapp --help
Usage: myapp <command> [options]

Available commands:
  deploy       Deploy the application
  start        Start the service (aliases: run, up)
  status       Show application status

Use 'myapp <command> --help' for command-specific help
```

### Command Help
```bash
$ go run myapp deploy --help
Usage: deploy [options]

Deploy the application to the specified environment.

Flags:
  --env string (required) (oneOf: dev staging prod)
        Target environment
  --dry-run bool (default: false)
        Show what would be deployed
```

## Examples

cli includes complete examples:

### Web Server Example
**Location:** `examples/web-server/`

A  web server demonstrating the **empty string command** approach for configuration-only applications:

```bash
cd examples/web-server

# Run with defaults (empty string command executes)
go run main.go

# Use environment variables
DATABASE_URL="postgres://user:pass@localhost/db" \
JWT_SIGNING_KEY="your-32-character-secret-key-here" \
go run main.go

# Override with flags
go run main.go --port 3000 --host 0.0.0.0 --log-level debug

# Get help for the default command
go run main.go --help

# Full help with all options
go run main.go --full-help
```

**Features demonstrated:**
- Empty string command for config-only mode
- Environment variable configuration
- Flag-based configuration
- Secret management
- Validation and error handling
- Help generation

### CLI Tool Example  
**Location:** `examples/cli-tool/`

A full-featured CLI tool with commands and middleware:
```bash
cd examples/cli-tool

# Deploy with validation
go run main.go deploy --env staging --dry-run=true

# Show system status
go run main.go status --detailed=true

# Manage configuration
go run main.go config --show-secrets=true
```

## Performance

cli is optimized for production use:

- **68% faster** template rendering
- **79% fewer** memory allocations
- **Silent overrides** for transparent CLI behavior
- **Zero boilerplate** configuration access
- **Thread-safe** concurrent operations

## API Reference

### Configuration Builder

| Method | Description |
| ------ | ----------- |
| `Define(key)` | Start defining a configuration key |
| `String()`, `Int64()`, `Float64()`, `Bool()`, `Duration()`, `StringSlice()`, `Int64Slice()` | Set value type |
| `Env(name)` | Set environment variable name |
| `Flag(name)` | Set command-line flag name |
| `File(key)` | Set file key name |
| `Default(value)` | Set default value |
| `Delimiter(d)` | Set delimiter for parsing slice flag values |
| `Required()` | Mark as required |
| `Secret()` | Mark as secret (memory protected) |
| `Description(text)` | Set description for help |
| `Priority(priority)` | Set per-definition source priority |

### Validation

| Method | Description |
| ------ | ----------- |
| `Min(n)`, `Max(n)` | Set numeric minimum / maximum |
| `Range(min, max)` | Set numeric range validation |
| `OneOf(values...)` | Set enum validation (string values) |
| `Regexp(pattern)` | Set regex validation |
| `MinLength(n)`, `MaxLength(n)` | Set string length bounds |
| `MinItems(n)`, `MaxItems(n)` | Set slice item count bounds |
| `MinDuration(d)`, `MaxDuration(d)` | Set duration bounds |
| `Custom(name, check)` | Add a custom validation function |

### Access Methods

| Method | Description |
| ------ | ----------- |
| `Get[T](ctx, key)` | Get value with type T (returns T, error) |
| `MustGet[T](ctx, key)` | Get value or panic on error |
| `GetSecret(key)` | Get a secret value (`IsSet`, `Size`, `String`, `Bytes`, `Destroy`) |
| `IsSecret(key)` | Check if a key is defined as secret |
| `Dump()` | Map of all configuration values (secrets masked) |
| `Execute(args)` | Execute with command routing |
| `Destroy()` | Securely wipe all secrets from memory |
| `LoadFile(path)` / `LoadFiles(paths...)` | Load configuration file(s) |
| `LoadFileFromEnv(envVar)` | Load config file whose path is in an env var |
| `SetDefaultPriority(p)` | Set global default source priority |
| `GenerateHelp()` | Return generated global help text |
| `ShowGlobalHelp()` / `ShowCommandHelp(name)` | Print help output |

### Command Builder

| Method | Description |
| ------ | ----------- |
| `Command(name)` | Define a new command |
| `Func(fn)` | Set command function |
| `ShortHelp(text)` | Set short help text |
| `LongHelp(text)` | Set long help text |
| `Aliases(names...)` | Set command aliases |
| `Config(fn)` | Define command-specific config (`cc.Define`, `cc.Command` for subcommands) |
| `SubCommand(name)` | Add a subcommand builder |
| `Middleware(fn)` | Add command-specific middleware |

Config-level middleware registration:

| Method | Description |
| ------ | ----------- |
| `UseMiddleware(fn)` | Apply middleware to all commands |
| `UseMiddlewareForCommands(names, fn)` | Apply middleware to specific commands |
| `UseMiddlewareForSubcommands(cmd, names, fn)` | Apply middleware to specific subcommands |

## Getting Started

1. **Install**: `go get github.com/fernandezvara/cli`
2. **Try Examples**: `cd examples/web-server && go run main.go`
3. **Read Documentation**: Check the examples for real-world patterns
4. **Build**: Start with configuration-only mode, add commands as needed

## License

MIT License - feel free to use cli in your projects!

---

**cli**: Professional CLI applications, simplified.
