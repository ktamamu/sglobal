# sglobal

AWS Security Group Public Access Scanner - A command-line tool to identify AWS Security Groups with risky public access rules.

## Features

- Scan specific AWS regions or all regions
- Exclude security groups using a file-based exclusion list
- Support for multiple output formats (JSON, text, Markdown)
- Fast concurrent scanning across multiple regions
- Detects public IP ranges including 0.0.0.0/0, ::/0, and other public CIDR blocks

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
| Region | Security Group ID | Group Name | Description | VPC ID | Protocol | Port(s) | CIDR |
|--------|------------------|------------|-------------|--------|----------|---------|------|
| us-east-1 | sg-12345678 | web-servers | Security group for web servers | vpc-12345678 | tcp | 80 | 0.0.0.0/0 |
| us-east-1 | sg-12345678 | web-servers | Security group for web servers | vpc-12345678 | tcp | 443 | 0.0.0.0/0 |
```

## GitHub Actions Integration

You can integrate sglobal into your GitHub Actions workflow for automated security monitoring:

```yaml
name: Check SG Global Rules

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  workflow_dispatch:    # Manual trigger

jobs:
  check_sg_rules:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build sglobal
        run: go build -o sglobal .

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Scan security groups
        run: ./sglobal --region all --output markdown > sg-report.md

      - name: Create issue if violations found
        if: success()
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('sg-report.md', 'utf8');
            if (report.includes('Found') && !report.includes('Found 0')) {
              github.rest.issues.create({
                owner: context.repo.owner,
                repo: context.repo.repo,
                title: 'Security Alert: Public Security Groups Detected',
                body: '## Security Group Scan Results\n\n' + report
              });
            }
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Command Line Options

```
  -e, --exclude-file string   file containing security group IDs to exclude (one per line)
  -o, --output string         output format: json, text, markdown (default "json")
  -r, --region string         AWS region to scan (default: current profile region, 'all' for all regions)
      --config string         config file (default is $HOME/.sglobal.yaml)
  -h, --help                  help for sglobal
```

## License

MIT License
