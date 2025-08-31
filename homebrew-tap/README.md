# Homebrew Tap for GoKit CLI

This is the Homebrew tap for the [GoKit CLI](https://github.com/kumarabd/gokit) - a comprehensive toolkit for building Go microservices.

## Installation

```bash
# Add the tap
brew tap kumarabd/gokit

# Install gokit
brew install gokit
```

## Usage

```bash
# Create a new service
gokit service --name my-service --template http

# Add features to existing service
gokit add --feature monitoring --service my-service
gokit add --feature tracing --service my-service

# Check version
gokit version
```

## Updating

```bash
brew update
brew upgrade gokit
```

## Uninstalling

```bash
brew uninstall gokit
```

## Development

To install from source:

```bash
git clone https://github.com/kumarabd/gokit.git
cd gokit
make install
```

## License

MIT License - see the [main repository](https://github.com/kumarabd/gokit) for details.
