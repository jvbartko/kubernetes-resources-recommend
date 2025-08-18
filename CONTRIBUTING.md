# 🤝 Contributing to Kubernetes Resources Recommend

Thank you for your interest in contributing to Kubernetes Resources Recommend! We welcome contributions from everyone, regardless of experience level.

## 🌍 Languages / 语言支持

- **English** (Current)
- **[中文贡献指南](CONTRIBUTING-zh.md)**

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to [project maintainers](mailto:your-email@example.com).

### Our Standards

- **Be Respectful**: Treat everyone with respect and kindness
- **Be Inclusive**: Welcome people of all backgrounds and identities
- **Be Collaborative**: Work together constructively
- **Be Professional**: Maintain professional behavior in all interactions

## 🚀 Getting Started

### Prerequisites

- Go 1.23.9 or higher
- Git
- Basic understanding of Kubernetes and Prometheus
- Familiarity with Go programming language

### First Time Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/kubernetes-resources-recommend.git
   cd kubernetes-resources-recommend
   ```

3. **Add the original repository as upstream**:
   ```bash
   git remote add upstream https://github.com/luozijian1990/kubernetes-resources-recommend.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Build the project**:
   ```bash
   make build
   # or
   go build -o bin/kubernetes-resources-recommend cmd/kubernetes-resources-recommend/main.go
   ```

## 🛠️ Development Setup

### Project Structure

```
kubernetes-resources-recommend/
├── cmd/                     # Main applications
├── internal/                # Private application code
│   ├── exporter/           # Export functionality
│   ├── prometheus/         # Prometheus client
│   ├── recommender/        # Core recommendation logic
│   └── types/              # Type definitions
├── pkg/                    # Public library code
└── docs/                   # Documentation
```

### Available Commands

```bash
make help          # View all available commands
make build         # Build the project
make test          # Run tests
make test-coverage # Run tests with coverage
make fmt           # Format code
make lint          # Lint code
make clean         # Clean build artifacts
```

### Environment Setup

Create a `.env` file for local development:
```bash
PROMETHEUS_URL=https://your-prometheus.example.com
CHECK_NAMESPACE=default
LIMITS=1.5
```

## 🔄 How to Contribute

### Types of Contributions

We welcome several types of contributions:

- 🐛 **Bug fixes**
- ✨ **New features**
- 📝 **Documentation improvements**
- 🧪 **Tests**
- 🌍 **Translations**
- 🎨 **UI/UX improvements**
- 📊 **Performance optimizations**

### Before You Start

1. **Check existing issues** to avoid duplicating work
2. **Create an issue** for new features or significant changes
3. **Discuss your approach** with maintainers if needed

### Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Write clean, readable code
   - Follow the existing code style
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test
   make test-coverage
   make lint
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** on GitHub

## 📥 Pull Request Process

### Before Submitting

- [ ] Code follows the project style guidelines
- [ ] Tests pass locally
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] No merge conflicts with main branch

### PR Title Format

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Performance improvement

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

1. **Automated checks** must pass
2. **Code review** by maintainers
3. **Testing** in different environments
4. **Approval** from project maintainers
5. **Merge** into main branch

## 🐛 Issue Guidelines

### Bug Reports

When reporting bugs, please include:

- **Clear title** and description
- **Steps to reproduce** the issue
- **Expected vs actual behavior**
- **Environment details**:
  - Go version
  - Kubernetes version
  - Prometheus version
  - Operating system
- **Logs and error messages**
- **Screenshots** if applicable

### Feature Requests

For feature requests, please provide:

- **Clear description** of the feature
- **Use case** and motivation
- **Proposed implementation** (if any)
- **Alternative solutions** considered
- **Additional context**

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature request
- `documentation`: Documentation improvements
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed
- `priority/high`: High priority issue
- `priority/low`: Low priority issue

## 🎨 Code Style

### Go Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused

### Naming Conventions

- **Packages**: lowercase, single word
- **Functions**: camelCase, start with uppercase if exported
- **Variables**: camelCase
- **Constants**: UPPER_CASE or camelCase
- **Files**: lowercase with underscores

### Code Organization

- Group related functionality
- Separate concerns into different packages
- Use interfaces for abstraction
- Handle errors appropriately
- Write self-documenting code

## 🧪 Testing

### Testing Strategy

- **Unit tests**: Test individual functions
- **Integration tests**: Test component interactions
- **End-to-end tests**: Test complete workflows

### Writing Tests

```go
func TestRecommenderFunction(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result := functionUnderTest(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/recommender -run TestSpecificFunction
```

### Test Coverage

- Aim for at least 80% test coverage
- Focus on critical business logic
- Test edge cases and error conditions

## 📝 Documentation

### Types of Documentation

- **Code comments**: Explain complex logic
- **README files**: Project overview and setup
- **API documentation**: Function and type documentation
- **User guides**: How-to guides and tutorials

### Documentation Standards

- Use clear, concise language
- Provide examples where helpful
- Keep documentation up-to-date
- Use proper markdown formatting

### Updating Documentation

When making changes:
- Update relevant README sections
- Add/update code comments
- Update API documentation
- Create/update user guides if needed

## 🏷️ Versioning

We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## 🎯 Roadmap

Current priorities:
- [ ] CPU resource recommendations
- [ ] GPU resource recommendations
- [ ] HPA integration
- [ ] Web UI interface
- [ ] API endpoints
- [ ] Helm chart

## 🙋‍♀️ Getting Help

If you need help:

1. **Check the documentation** first
2. **Search existing issues** on GitHub
3. **Create a new issue** with details
4. **Join discussions** in issues and PRs
5. **Contact maintainers** directly if needed

## 🏆 Recognition

Contributors will be recognized in:
- **README.md** contributor section
- **Release notes** for significant contributions
- **Special mentions** in project updates

Thank you for contributing to Kubernetes Resources Recommend! 🎉

---

## 📧 Contact

- **Maintainer**: [luozijian1990](https://github.com/luozijian1990)
- **Issues**: [GitHub Issues](https://github.com/luozijian1990/kubernetes-resources-recommend/issues)
- **Discussions**: [GitHub Discussions](https://github.com/luozijian1990/kubernetes-resources-recommend/discussions)
