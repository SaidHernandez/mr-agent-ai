
# Supply Chain Security Agent

## When to Use

- Adding a new dependency to the project
- Reviewing a Dependabot or Renovate update PR
- Auditing the project's current package security posture
- A third-party PR modifies lock files or adds packages

---

## Decision Tree

```
Adding a new dependency?      → Check shelf-time + reputation before installing
Reviewing a Dependabot PR?    → Check changelog scope + changed files
Third-party PR touches deps?  → Look for new install scripts or unexpected network calls
Dependency suddenly does HTTP?→ Red flag — review before merging
```

---

## Shelf Time

Wait at least **7 days** after a version is published before installing it.
Most supply chain attacks are distributed in the first hours after a package is compromised.

**npm** — add to `.npmrc`:
```ini
min-release-age=7d
ignore-scripts=true
```

**Python / uv** — add to `pyproject.toml`:
```toml
[tool.uv]
exclude-newer = "1 week"
```

**Python / uv** — or per-package exception for urgent CVE fixes:
```toml
[tool.uv]
exclude-newer = "1 week"
exclude-newer-package = { requests = "2026-03-25T16:00:00Z" }
```

**Go** — the Go checksum database (`sum.golang.org`) verifies integrity by default.
Add `go mod verify` to your CI pipeline to detect tampered modules.

**Java / Maven** — no native shelf-time. Use Dependabot + explicit version pinning in `pom.xml`.

**Rules:** Never install a dependency with `--ignore-scripts` overridden. Always use lock files.

---

## Reviewing Dependency Update PRs

Before approving a Dependabot or Renovate PR:

1. **Check the version bump scope** — patch bumps are usually safe; major bumps need careful review
2. **Read the changelog** — look for new features that involve network calls, file system access, or subprocess execution
3. **Diff the lock file** — unexpected transitive dependency changes are a red flag
4. **Check for new install scripts** — `postinstall`, `prepare`, `setup.py` changes deserve extra scrutiny

**Rules:** Never auto-merge a major version bump. Always review the diff, not just the title.

---

## Red Flags

Stop and investigate before merging when you see:

- **New maintainer** — package ownership transferred recently (check npm/PyPI history)
- **New install scripts** — `postinstall`, `prepare`, or `setup.py` added in an update
- **Typosquatting** — package name very similar to a popular one (`lodahs` vs `lodash`)
- **Sudden network calls** — a previously offline library now imports `http`, `urllib`, `net/http`
- **Obfuscated code** — minified or encoded strings in a non-frontend package
- **Version confusion attack** — an internal package name now exists on the public registry

---

## Secret Scanning

Prevent credentials from entering the git history.

**gitleaks** — scans commits for leaked secrets:
```bash
gitleaks detect --source . --verbose
```

**trufflehog** — finds secrets in git history:
```bash
trufflehog git file://. --only-verified
```

**Pre-commit hook** — blocks commits containing secrets:
```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.0
    hooks:
      - id: gitleaks
```

**Rules:** Configure at least one secret scanner in CI. Never commit `.env` files — add them to `.gitignore`.

---

## Lock Files

Lock files guarantee reproducible installs and prevent version drift.

| Ecosystem | Lock file |
|-----------|-----------|
| npm | `package-lock.json` / `yarn.lock` / `pnpm-lock.yaml` |
| Python | `uv.lock` / `poetry.lock` |
| Go | `go.sum` |
| Ruby | `Gemfile.lock` |

**Rules:** Always commit lock files. Never ignore lock file changes in PRs — they are part of the security review.

---

## Audit Commands

```bash
# Node.js
npm audit
npm audit fix

# Python
pip-audit
pip-audit --fix

# Go
go mod verify
govulncheck ./...

# Ruby
bundle audit
```

---

## Definition of Done

- [ ] Shelf-time configured for all detected ecosystems
- [ ] Lock file committed and reviewed in every dependency PR
- [ ] At least one secret scanner (gitleaks or trufflehog) active in CI
- [ ] Dependabot or Renovate configured for automated dependency updates
- [ ] No install scripts ignored (`ignore-scripts=true` for npm)
