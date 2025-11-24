# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2025-11-24

### Added
- `--version` / `-v` flag to display version information
- Build information (version, commit hash, build date) embedded in binary
- Version info automatically injected by GoReleaser during releases

## [0.1.2] - 2025-11-24

### Changed
- Migrated to GoReleaser for automated release management
- Homebrew formula now uses pre-built binaries instead of building from source
- Simplified release workflow with GoReleaser automation
- Removed Go dependency from Homebrew installation (faster installation)

## [0.1.1] - 2025-11-23

### Changed
- Updated AWS SDK for Go v2 dependencies:
  - `github.com/aws/aws-sdk-go-v2`: 1.36.5 → 1.39.6
  - `github.com/aws/aws-sdk-go-v2/config`: 1.29.18 → 1.31.8
  - `github.com/aws/aws-sdk-go-v2/service/ec2`: 1.227.0 → 1.255.0
- Updated `github.com/spf13/cobra` from 1.9.1 to 1.10.1
- Updated `github.com/spf13/viper` from 1.20.1 to 1.21.0
- Updated GitHub Actions workflows:
  - `actions/checkout` from 4 to 5
  - `actions/setup-go` from 5 to 6
  - `actions/upload-artifact` from 4 to 5
  - `golangci/golangci-lint-action` from 6 to 9

### Fixed
- README documentation improvements
- Icon display fixes

## [0.1.0] - 2025-06-28

### Added
- Initial release of sglobal
- AWS Security Group public access scanner
- Multi-region scanning support (specific region or all regions)
- Security group exclusion list support via file
- Multiple output formats (JSON, text, Markdown)
- Concurrent scanning across multiple regions
- Detection of public IP ranges (0.0.0.0/0, ::/0, and other public CIDR blocks)
- Command-line interface with flags:
  - `--region`: Specify AWS region or scan all regions
  - `--exclude-file`: File containing security group IDs to exclude
  - `--output`: Choose output format (json, text, markdown)
- GitHub Actions integration examples
- Homebrew installation support
- MIT License

### Dependencies
- Go 1.23
- AWS SDK for Go v2
- Cobra v1.8.0 (CLI framework)
- Viper v1.18.2 (configuration)

[Unreleased]: https://github.com/ktamamu/sglobal/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/ktamamu/sglobal/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/ktamamu/sglobal/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/ktamamu/sglobal/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/ktamamu/sglobal/releases/tag/v0.1.0
