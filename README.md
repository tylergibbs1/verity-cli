# Verity CLI

Command-line interface for the [Verity API](https://verity.backworkai.com) - Medicare coverage policies, prior authorization requirements, and medical code lookups.

## Installation

### Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/tylerbryy/verity-cli/releases).

```bash
# macOS (Intel)
curl -L https://github.com/tylerbryy/verity-cli/releases/latest/download/verity-darwin-amd64 -o verity
chmod +x verity
sudo mv verity /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/tylerbryy/verity-cli/releases/latest/download/verity-darwin-arm64 -o verity
chmod +x verity
sudo mv verity /usr/local/bin/

# Linux
curl -L https://github.com/tylerbryy/verity-cli/releases/latest/download/verity-linux-amd64 -o verity
chmod +x verity
sudo mv verity /usr/local/bin/

# Windows
# Download verity-windows-amd64.exe from releases
```

### From Source

```bash
go install github.com/tylerbryy/verity-cli@latest
```

## Quick Start

```bash
# Set your API key
export VERITY_API_KEY=vrt_live_YOUR_API_KEY

# Look up a medical code
verity check 76942

# Search policies
verity policies list --query "ultrasound guidance"

# Check prior authorization
verity prior-auth 76942 --state TX --diagnosis M54.5

# Get policy details
verity policies get L33831 --include criteria,codes
```

## Configuration

The CLI looks for configuration in the following order:
1. Command-line flags
2. Environment variables (prefixed with `VERITY_`)
3. Config file (`~/.verity.yaml`)

### Config File Example

```yaml
# ~/.verity.yaml
api_key: vrt_live_YOUR_API_KEY
base_url: https://verity.backworkai.com/api/v1
output: table
```

### Environment Variables

```bash
export VERITY_API_KEY=vrt_live_YOUR_API_KEY
export VERITY_BASE_URL=https://verity.backworkai.com/api/v1
export VERITY_OUTPUT=json
```

## Commands

### `verity check [code]`

Look up a medical code (CPT, HCPCS, ICD-10, NDC).

```bash
# Basic lookup
verity check 76942

# Include RVU data
verity check 76942 --include rvu

# Include policies
verity check 76942 --include rvu,policies

# Filter by jurisdiction
verity check 76942 --jurisdiction JM

# JSON output
verity check 76942 --output json
```

**Flags:**
- `-i, --include`: Include additional data (rvu, policies)
- `-j, --jurisdiction`: Filter by MAC jurisdiction
- `-f, --fuzzy`: Enable fuzzy matching (default: true)
- `-o, --output`: Output format (table, json, yaml)

### `verity policies list`

Search and list policies.

```bash
# Search policies
verity policies list --query "ultrasound guidance"

# Filter by type
verity policies list --type LCD

# Filter by jurisdiction
verity policies list --jurisdiction JM

# Semantic search
verity policies list --query "imaging procedures" --mode semantic

# Include retired policies
verity policies list --status all
```

**Flags:**
- `-q, --query`: Search query
- `-m, --mode`: Search mode (keyword, semantic)
- `-t, --type`: Policy type (LCD, Article, NCD)
- `-j, --jurisdiction`: MAC jurisdiction
- `-s, --status`: Status (active, retired, all)

### `verity policies get [policy-id]`

Get detailed information about a specific policy.

```bash
# Basic policy info
verity policies get L33831

# Include criteria
verity policies get L33831 --include criteria

# Include codes and attachments
verity policies get L33831 --include codes,attachments,criteria
```

**Flags:**
- `-i, --include`: Include additional data (criteria, codes, attachments, versions)

### `verity prior-auth [procedure-codes...]`

Check prior authorization requirements.

```bash
# Check single procedure
verity prior-auth 76942 --state TX

# Check multiple procedures
verity prior-auth 76942 76937 --state TX

# Include diagnosis codes
verity prior-auth 76942 --diagnosis M54.5,G89.29 --state TX

# Check for different payer
verity prior-auth 76942 --state TX --payer uhc
```

**Flags:**
- `-d, --diagnosis`: Diagnosis codes (ICD-10), comma-separated
- `-s, --state`: Two-letter state code
- `-p, --payer`: Payer (medicare, aetna, uhc, all)

## Global Flags

These flags work with all commands:

- `--api-key`: Verity API key
- `--base-url`: API base URL
- `--config`: Config file path
- `-o, --output`: Output format (table, json, yaml)

## Examples

### Check if a procedure needs prior auth in Texas

```bash
verity prior-auth 76942 --state TX --diagnosis M54.5
```

### Find all LCD policies about ultrasound

```bash
verity policies list --query "ultrasound" --type LCD --output json
```

### Look up a code and get pricing info

```bash
verity check 76942 --include rvu
```

### Get full details of a specific policy

```bash
verity policies get L33831 --include criteria,codes --output json
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/tylerbryy/verity-cli.git
cd verity-cli

# Build for your platform
go build -o verity .

# Or build for all platforms
make build-all
```

## License

MIT License - see LICENSE file for details.

## Support

- Documentation: https://verity.backworkai.com/docs
- Issues: https://github.com/tylerbryy/verity-cli/issues
- Email: support@verity.backworkai.com
