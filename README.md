# SMTP TLS Test

A simple Go program that demonstrates sending emails with and without TLS encryption.

## Installation

```bash
# Clone the repository
git clone https://github.com/mrichman/smtp_tls_test.git
cd smtp_tls_test

# Build the application
go build -o smtp_tls_test

# Or use the provided task runner
task build
```

## Usage

```bash
# Show help
smtp_tls_test --help

# Send email without TLS (uses STARTTLS if available)
smtp_tls_test send

# Send email with TLS (forces direct TLS connection)
smtp_tls_test send --tls

# Enable verbose mode to log the entire SMTP conversation
smtp_tls_test send --verbose

# Combine flags
smtp_tls_test send --tls --verbose

# Create a default configuration file
smtp_tls_test config create

# Show current configuration
smtp_tls_test config show

# Validate configuration
smtp_tls_test config validate

# Specify a custom configuration file
smtp_tls_test send --config /path/to/config.json

# Show version information
smtp_tls_test version

# Using Task runner (if installed)
task run           # Run without TLS
task run:tls       # Run with TLS
task run:verbose   # Run in verbose mode
task run:tls-verbose # Run with TLS and verbose mode
```

## Configuration

You can configure the application using a JSON configuration file:

```json
{
  "smtp": {
    "host": "smtp.example.com",
    "port": 587,
    "username": "username@example.com",
    "from": "sender@example.com",
    "password": "your_password",
    "to": ["recipient@example.com"]
  }
}
```

To create a default configuration file:

```bash
smtp_tls_test config create
```

You can also specify a custom configuration file path:

```bash
smtp_tls_test --config /path/to/config.json
```

## How It Works

The program automatically detects the appropriate TLS method based on the port number:

- Port 465: Uses direct TLS connection (implicit TLS)
- Port 587: Uses STARTTLS if supported by the server (explicit TLS)

Additionally, the program accepts a command line flag `--tls` that forces a direct TLS connection regardless of port:

- When `--tls` is not provided (default), the program follows the port-based detection logic.
- When `--tls` is provided, the program establishes a direct TLS connection to the SMTP server.

## Security Note

The example code includes `InsecureSkipVerify: true` for demonstration purposes. In a production environment, you should set this to `false` to properly verify the server's certificate.

## Dependencies

This program uses the following dependencies:

- Go standard library
- [Cobra](https://github.com/spf13/cobra) for command line interface

## Task Runner

A Taskfile.yml is included for easy execution of common tasks. If you have [Task](https://taskfile.dev/) installed, you can use the following commands:

```bash
task                # Show available tasks
task build          # Build the application
task run            # Run without TLS
task run:tls        # Run with TLS
task run:verbose    # Run in verbose mode
task run:tls-verbose # Run with TLS and verbose mode
task test           # Run tests
task lint           # Run linters
task fmt            # Format code
task vet            # Run go vet
task clean          # Remove build artifacts
task all            # Run format, lint, vet, test and build
task help           # Show help for the program
```

To install Task, follow the instructions at [taskfile.dev](https://taskfile.dev/installation/) or run:

```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
```

Once Task is installed, you can use all the task commands listed above.
