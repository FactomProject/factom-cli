# Factom Command Line Interface

[![CircleCI](https://circleci.com/gh/FactomProject/factom-cli/tree/develop.svg?style=shield)](https://circleci.com/gh/FactomProject/factom-cli/tree/develop)

`factom-cli` provides a convenient CLI interface for making and viewing entries
and transactions on the Factom blockchain by calling out to both
[`factomd`](https://github.com/FactomProject/factomd) and
[`factom-walletd`](https://github.com/FactomProject/factomd).

## Dependencies
### Build Dependencies
- Go 1.10 or higher

### External Dependencies
- Access to a `factomd` API endpoint
- Access to a `factom-walletd` API endpoint

### Optional Dependencies
- [CLI completion](https://github.com/AdamSLevy/complete-factom-cli) for Bash,
  Zsh or Fish

## Build and install
```
go get -u github.com/FactomProject/factom-cli
go install github.com/FactomProject/factom-cli
```
