# Contributing

We 💛 contributions! This document outlines the process for contributing to the Supabase Terraform Provider.

## Rules

1. **Don't be a jerk**
2. Search existing issues before opening new ones
3. Lint and test your changes locally before submitting PRs
4. Follow the project's existing code style and conventions
5. Include tests with your contributions when applicable

## Getting Started

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

### Setup

1. Fork and clone the repository
2. Enter the repository directory
3. Install dependencies:

```shell
go mod download
```

4. Build the provider:

```shell
go install
```

## Development Workflow

### Before You Start

1. Search existing [issues](https://github.com/shellscape/terraform-provider-supabase/issues) to see if your contribution is already being discussed
2. If not, open an issue to discuss your proposed changes
3. Fork the repository and create a feature branch

### Making Changes

1. Make your changes following the existing code style
2. Add or update tests as needed
3. Update documentation if applicable
4. Ensure all tests pass locally

### Building and Testing

#### Building the Provider

To compile the provider locally:

```shell
go install
```

This builds the provider and puts the binary in your `$GOPATH/bin` directory.

#### Generating Documentation

To generate or update documentation:

```shell
go generate
```

#### Running Tests

Run the full test suite before submitting your PR:

```shell
# Run unit tests
go test ./...

# Run acceptance tests (creates real resources - may cost money)
make testacc
# or manually:
TF_ACC=1 go test ./internal/provider -v -timeout 120s
```

**Note:** Acceptance tests create real Supabase resources and may incur costs.

### Submitting Changes

1. **Lint your changes**: Ensure code follows Go conventions
2. **Test locally**: All tests must pass
3. **Generate docs**: Run `go generate` to update documentation
4. **Create a Pull Request**: Include a clear description of changes

## Pull Request Guidelines

- Complete the PR template fully
- Explain what your code changes do and how to test them
- Link to related GitHub issues when applicable
- Resolve all automated check failures
- Be responsive to code review feedback

## Working with Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules). Please see the Go documentation for the most up-to-date information about using Go modules.

### Adding New Dependencies

To add a new dependency `github.com/author/dependency` to the provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Development Setup

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

The provider binary will be available in your `$GOPATH/bin` directory after building.

## Need Help?

- Check existing [issues](https://github.com/shellscape/terraform-provider-supabase/issues)  
- Review the [step-by-step tutorial](docs/tutorial.md)
- Ask questions in your issue or pull request
