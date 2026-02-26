# Contributing to TronCLI

First off, thanks for taking the time to contribute!

The following is a set of guidelines for contributing to TronCLI. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for TronCLI. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for TronCLI, including completely new features and minor improvements to existing functionality.

### Your First Code Contribution

Unsure where to begin contributing to TronCLI? You can start by looking through these `good-first-issue` and `help-wanted` issues:

*   [Good First Issues](https://github.com/rsdenck/troncli/labels/good%20first%20issue) - issues which should only require a few lines of code, and a test or two.
*   [Help Wanted Issues](https://github.com/rsdenck/troncli/labels/help%20wanted) - issues which should be a bit more involved than `good-first-issue` issues.

## Styleguides

### Git Commit Messages

*   Use the present tense ("Add feature" not "Added feature")
*   Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
*   Limit the first line to 72 characters or less
*   Reference issues and pull requests liberally after the first line
*   **IMPORTANT**: We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.
    *   `feat: add new feature`
    *   `fix: resolve bug`
    *   `docs: update documentation`
    *   `style: format code`
    *   `refactor: restructure code`
    *   `test: add tests`
    *   `chore: update build tasks`

### Go Styleguide

*   Follow [Effective Go](https://golang.org/doc/effective_go.html).
*   Run `make lint` before submitting a PR.
*   Run `make test` to ensure all tests pass.

## Development Setup

1.  Clone the repository
2.  Install dependencies: `go mod download`
3.  Install development tools:
    *   `golangci-lint`
    *   `commitlint` (Node.js)

## Pull Request Process

1.  Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2.  Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
3.  Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
4.  You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.
