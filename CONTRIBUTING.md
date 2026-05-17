# Contributing to go-output

Thank you for your interest in contributing to go-output!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-output.git`
3. Add the upstream remote: `git remote add upstream https://github.com/larsartmann/go-output.git`
4. Create a feature branch: `git checkout -b feature/your-feature-name`

## Development Setup

This project uses a Go multi-module workspace with 8 independent modules.

### Initial Setup

```bash
# Create go.work for local development (gitignored)
cat > go.work << 'EOF'
go 1.26.2

use (
  .
  ./enum
  ./escape
  ./testhelpers
  ./cmdguard
  ./sort
  ./table
  ./integration
  ./examples
)
EOF
```

### Build & Test

```bash
# Build all modules (from project root with go.work)
go build ./...

# Test all modules
go test ./...

# Test a specific module
go test ./enum/...

# Test with race detector
go test -race ./...

# Test with coverage
go test -cover ./...

# Lint all modules
golangci-lint run ./...
```

### Per-Module Commands (without go.work)

Each module is standalone and can be built/tested independently:

```bash
cd table && go test ./... && cd ..
cd cmdguard && go test ./... && cd ..
```

### Tidy Dependencies

After changing imports, tidy each affected module:

```bash
go mod tidy                    # Root module
cd enum && go mod tidy && cd ..    # Sub-module
```

## Making Changes

1. Keep changes focused and atomic
2. Follow existing code style (enforced by golangci-lint)
3. Add tests for new functionality
4. Update documentation as needed
5. Keep commits small and descriptive

## Commit Messages

Follow the conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`

Examples:

- `feat(table): add streaming renderer support`
- `fix(d2): correct SQL table constraint rendering`
- `docs(readme): add registry system documentation`

## Testing

All changes must pass tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run race detector
go test -race ./...

# Run fuzz tests
go test -fuzz=FuzzParseOutputFormat -fuzztime=1m .
go test -fuzz=FuzzParseSortBy -fuzztime=1m .
```

## Code Quality

- All code must pass `golangci-lint run`
- File size limit: 350 lines
- Test coverage target: 90%+
- No code duplication (threshold: 30 tokens)
- Sub-modules must not import `internal/` packages from root

## Pull Request Process

1. Update the README.md if you add new features
2. Update CHANGELOG.md with your changes (following [Keep a Changelog](https://keepachangelog.com/) format)
3. Ensure all tests pass and linter is clean
4. Push your branch and create a pull request
5. Respond to review feedback

## Reporting Issues

When reporting issues, include:

- Go version (`go version`)
- Library version (git commit or tag)
- Minimal reproduction case
- Expected vs actual behavior

## Questions?

Feel free to open an issue for questions or discussions.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
