#!/bin/bash

# Welcome Message Generator
# Generates a welcome message for new contributors

cat << 'EOF'
# Welcome to OAuth 2.0 SSO! 👋

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Development Setup

```bash
# Clone the repo
git clone https://github.com/your-username/OAuth_2.0.git
cd OAuth_2.0

# Copy environment file
cp .env.example .env

# Edit .env with your Auth0 credentials
# Run the app
go run main.go
```

## Need Help?

- Check the [README](../README.md) for documentation
- Open an [issue](https://github.com/qqwozz/OAuth_2.0/issues) for bugs
- Start a [discussion](https://github.com/qqwozz/OAuth_2.0/discussions) for questions

## Code of Conduct

Please be respectful and inclusive in all interactions.

Happy coding! 🚀
EOF
