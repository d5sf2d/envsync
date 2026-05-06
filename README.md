# envsync

> Utility to diff and sync `.env` files across environments with secret masking

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git && cd envsync && go build ./...
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.staging .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync --source .env.staging --target .env.production
```

**Mask secrets in output:**

```bash
envsync diff .env.staging .env.production --mask
```

Example output:

```
~ API_URL        staging.example.com  →  production.example.com
+ NEW_FEATURE_FLAG  true
- DEPRECATED_KEY    ********  (masked)
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--mask` | Redact secret values in diff output |
| `--dry-run` | Preview sync changes without writing |
| `--output` | Output format: `text`, `json`, `yaml` |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername