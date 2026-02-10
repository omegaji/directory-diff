# directory-diff
A tiny diff engine to detect any new files being added, deleted or modified
# directory-diff

**A tiny diff engine to detect added, modified, or deleted files in directories**  
This tool scans one or more directories and reports which files were added, changed, or removed — similar to how Git detects work‑tree changes.:contentReference[oaicite:1]{index=1}

---

## 🔧 Features

- Detects **added**, **modified**, and **removed** files
- Simple CLI tool written in Go
- Fast and tiny — perfect for scripting or automation

---

## 🚀 Install

### Option A — Install from source (Go)

If you have Go installed (1.18+), build and install the CLI directly:
```bash
go install github.com/omegaji/directory-diff/cmd/directory-diff@latest
```

### Option B - Download prebuilt binaries

```bash
We publish prebuilt binaries for common platforms under Releases.

Go to:
https://github.com/omegaji/directory-diff/releases

Download the binary for your OS/architecture (e.g., directory-diff-linux-amd64)

Make it executable:

chmod +x directory-diff-linux-amd64