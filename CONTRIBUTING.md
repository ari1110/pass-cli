# Contributing to Pass-CLI

Thank you for your interest in contributing to Pass-CLI! This document provides guidelines for contributing to the project.

## Documentation Governance

Pass-CLI maintains a [Documentation Lifecycle Policy](docs/DOCUMENTATION_LIFECYCLE.md) that defines retention periods, archival triggers, and decision workflows for all repository documentation. Contributors should consult the policy when adding new documentation or proposing changes to existing docs. The policy ensures documentation remains current and maintainable while preserving historical design context.

## How to Contribute

Contributions are welcome in the following areas:
- Bug reports and feature requests via GitHub Issues
- Code contributions via Pull Requests
- Documentation improvements
- Testing and quality assurance

## Pull Request Process

1. Fork the repository and create a feature branch
2. Make your changes following the project's coding standards
3. Add tests for any new functionality
4. Ensure all tests pass before submitting
5. Submit a pull request with a clear description of your changes

## Code Standards

- Follow Go best practices and idioms
- Write clear, concise commit messages
- Add documentation for new features
- Maintain backward compatibility where possible

## Testing

All code contributions should include appropriate tests. Run the test suite before submitting:

```bash
go test ./...
```

## Questions?

If you have questions about contributing, please open a GitHub Issue for discussion.
