# Factom Command Line Interface

[![CircleCI](https://circleci.com/gh/FactomProject/factom-cli/tree/develop.svg?style=shield)](https://circleci.com/gh/FactomProject/factom-cli/tree/develop)

`factom-cli` provides a convenient CLI interface for making and viewing entries
and transactions on the Factom blockchain by calling out to both
[`factomd`](https://github.com/FactomProject/factomd) and
[`factom-walletd`](https://github.com/FactomProject/factomd).

## Dependencies
### Build Dependencies
- Go 1.13 or higher

### External Dependencies
- Access to a `factomd` API endpoint
- Access to a `factom-walletd` API endpoint

### Optional Dependencies
- [CLI completion](https://github.com/AdamSLevy/complete-factom-cli) for Bash,
  Zsh or Fish

## Package distribution

Binaries for your platform can be downloaded from the [GitHub release page](https://github.com/FactomProject/factom-cli/releases).

## Build and install

```
make install
```

To cross compile to all supported platforms:
```
make all
```
