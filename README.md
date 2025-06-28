# sglobal

AWS Security Group Global Access Scanner - A command-line tool to identify AWS Security Groups with risky global access rules (0.0.0.0/0).

## Features

- Scan specific AWS regions or all regions
- Exclude security groups using a file-based exclusion list
- Support for multiple output formats (JSON, text, Markdown)
- Fast concurrent scanning across multiple regions

## Installation

### Using Homebrew (coming soon)

```bash
brew tap ktamamu/tap
brew install sglobal
```

### From Source

```bash
git clone https://github.com/ktamamu/sglobal.git
cd sglobal
go build -o sglobal
```

## Usage

### Basic Usage

```bash
# Scan current region
sglobal

# Scan specific region
sglobal --region us-east-1

# Scan all regions
sglobal --region all

# Use exclusion file
sglobal --exclude-file excludes.txt

# Text output
sglobal --output text

# Markdown output
sglobal --output markdown
```

### Exclusion File Format

Create a text file with one Security Group ID per line:

```
sg-12345678
sg-87654321
# This is a comment - lines starting with # are ignored
sg-abcdef12
```

### Output Formats

#### JSON Output (default)
```json
{"count":2,"results":[{"region":"us-east-1","security_group_id":"sg-12345678","group_name":"web-servers","description":"Security group for web servers","vpc_id":"vpc-12345678","risky_rules":[{"from_port":80,"to_port":80,"protocol":"tcp","cidr_blocks":["0.0.0.0/0"]}]}]}
```

#### Text Output
```
Found 2 security groups with global access:

Region: us-east-1
Security Group ID: sg-12345678
Group Name: web-servers
Description: Security group for web servers
VPC ID: vpc-12345678
Risky Inbound Rules:
  - Protocol: tcp, Port: 80, CIDR: 0.0.0.0/0
  - Protocol: tcp, Port: 443, CIDR: 0.0.0.0/0
```

#### Markdown Output
```markdown
# Security Groups with Global Access

**Found 2 security groups with global access**

## us-east-1

### sg-12345678 (web-servers)
- **Description:** Security group for web servers
- **VPC ID:** vpc-12345678
- **Risky Inbound Rules:**
  - Protocol: tcp, Port: 80, CIDR: 0.0.0.0/0
  - Protocol: tcp, Port: 443, CIDR: 0.0.0.0/0
```

## Configuration

The tool uses AWS SDK default configuration. Make sure you have:

1. AWS credentials configured (via AWS CLI, environment variables, or IAM roles)
2. Appropriate permissions to describe EC2 security groups and regions

### Required AWS Permissions

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeRegions"
            ],
            "Resource": "*"
        }
    ]
}
```

## Command Line Options

```
  -e, --exclude-file string   file containing security group IDs to exclude (one per line)
  -o, --output string         output format: json, text, markdown (default "json")
  -r, --region string         AWS region to scan (default: current profile region, 'all' for all regions)
      --config string         config file (default is $HOME/.sglobal.yaml)
  -h, --help                  help for sglobal
```

## Development

### Prerequisites

- Go 1.22 or later
- golangci-lint (for linting)

### Building from source

```bash
# Clone the repository
git clone https://github.com/ktamamu/sglobal.git
cd sglobal

# Install dependencies
make deps

# Install linting tools
make install-tools

# Run all checks (format, vet, lint, test)
make check

# Build
make build
```

### Running tests

```bash
# Run tests
make test

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
```

### Code quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all quality checks
make check
```

## License

MIT License