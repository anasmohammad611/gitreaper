# GitReaper ☠️

Clean up your Git repository like a pro.  
**GitReaper** safely reaps stale, merged, or orphaned branches and cleans your workspace with intelligent, configurable rules.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/anasmohammad611/gitreaper/ci.yml?branch=main)](https://github.com/anasmohammad611/gitreaper/actions)
[![Latest Release](https://img.shields.io/github/v/release/anasmohammad611/gitreaper)](https://github.com/yourusername/gitreaper/releases)

---

## ⚡️ Features

- 🔍 **Dry run mode** with color-coded output
- 🧹 Delete local branches merged into `main`, `master`, or `develop`
- 🚫 Automatically protect important branches (via `.gitreaperrc`)
- 📅 Age-based cleanup (`--older-than`, `--keep-recent`)
- 🧠 Interactive mode with commit info and author context
- 🔥 `--force`, `--aggressive` cleanup for CI or cron
- 📊 Post-cleanup reports: what was deleted and how much space was saved
- 🪄 Git hook integration (e.g., run before `git push`)
- 🛠 Configurable per-repo or globally

---

## 📦 Installation

### Using `go install`:

```bash
go install github.com/yourusername/gitreaper@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

### Or build from source:

```bash
git clone https://github.com/yourusername/gitreaper.git
cd gitreaper
make build
./build/gitreaper --help
```

---

## 🚀 Usage

```bash
gitreaper --dry-run
```

Example:

```bash
gitreaper --older-than 30d --exclude "release/*" --interactive
```

Available flags:

| Flag             | Description                                  |
|------------------|----------------------------------------------|
| `--dry-run`      | Preview deletions without making changes     |
| `--older-than`   | Clean branches older than X days             |
| `--keep-recent`  | Keep branches modified in last X days        |
| `--force`        | Skip all confirmations                       |
| `--interactive`  | Show branch-by-branch prompt                 |
| `--config`       | Use custom config file                       |

---

## 🛡 Configuration (`.gitreaperrc`)

Create a `.gitreaperrc` file in your repo root or `~/.gitreaperrc`:

```toml
protect = ["main", "master", "develop", "release/*"]
dry_run = true
interactive = true
older_than_days = 30
```

---

## 📂 Project Structure

```
gitreaper/
├── cmd/gitreaper         # CLI entrypoint
├── internal/             # App logic modules
│   ├── git/              # Git command wrappers
│   ├── reaper/           # Branch cleanup logic
│   ├── config/           # Config parsing
│   ├── ui/               # Interactive prompts, output
│   └── errors/           # Custom error handling
├── pkg/                  # Exportable packages
├── docs/                 # Project documentation
├── scripts/              # Dev tools and helpers
├── testdata/             # Test fixtures
├── .github/workflows/    # CI/CD configs
├── .gitignore
├── LICENSE
├── Makefile
└── README.md
```

---

## 🧪 Testing

Run all tests:

```bash
make test
```

---

## 🧰 Development

Install dependencies:

```bash
make deps
```

Run linter:

```bash
make lint
```

---

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

## ✨ Future Roadmap

- [ ] GitHub remote branch cleanup
- [ ] GUI or TUI interface (fzf-style)
- [ ] VS Code extension
- [ ] Auto-scheduling or daemon mode
