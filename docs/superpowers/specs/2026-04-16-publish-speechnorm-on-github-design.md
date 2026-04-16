# Publish `speechnorm` as a public Go module — design

**Date:** 2026-04-16
**Status:** approved (pending user spec review)
**Author:** brainstormed with Claude
**Module path (final):** `github.com/olidow/speechnorm`

## 1. Goal

Publish the `speechnorm` Go library so other Go projects can `go get` it instead of vendoring or copy-pasting the source. The end state:

- A public GitHub repo at `github.com/olidow/speechnorm`.
- Importable as `github.com/olidow/speechnorm` with no `internal/` indirection.
- Tagged releases automated by Conventional Commits via release-please.
- A baseline of supply-chain trust signals (Tier 2 from the research): pinned action SHAs, signed tags, branch protection, Scorecard, govulncheck, Dependabot — but no SLSA/cosign/SBOM (over-engineering for a library with no binaries).
- Documentation that renders well on `pkg.go.dev` (godoc, runnable examples, README, license).

Non-goals: shipping binaries, fuzzing infrastructure, multi-module monorepo layout, vanity import paths.

## 2. Decisions captured during brainstorming

| Decision | Choice | Rationale |
|---|---|---|
| Repo layout | Move source out of `internal/` to repo root | The current `go.mod` inside `internal/` works but is socially broken (every IDE/lint tool assumes `internal/` is unimportable). One-time cheap fix since there are no consumers yet. |
| Module path | `github.com/olidow/speechnorm` | Personal GitHub account, simplest ownership story, immutable from first publish. |
| `go` directive | Keep `go 1.24` | User is the only consumer right now; can lower later if external users hit version-floor problems. |
| Workflow rhythm | PR-based, even solo | Enables release-please auto-versioning + changelog; required CI checks; signals professionalism to consumers and Scorecard. |
| Test matrix | Linux × Go {1.23, 1.24} required + macOS/Windows × Go 1.24 non-blocking | Locale-heavy regex/Unicode code genuinely benefits from cross-platform coverage. |
| Coverage | Codecov via tokenless OIDC | Free, free badge, low setup. |
| Linter | golangci-lint v2 with library-friendly default set | `govet, staticcheck, errcheck, ineffassign, unused, gofumpt, goimports, revive, errorlint, misspell` |
| Release automation | `googleapis/release-please-action@v4`, `release-type: go` | Conventional-Commits-driven version bumps + changelog + GitHub Releases as one button-click. |
| Trust posture | Tier 2 (sweet spot) | Strong, legible signals to humans + scanners without Tier 3 maintenance burnout. |
| Signed commits/tags | SSH signing key registered on GitHub | Lower-friction than GPG; gets the "Verified" badge; satisfies branch-protection rule. |
| First tag | `v0.1.0` | Pre-1.0 freedom to break the API; stay in `v0.x` until the public API has survived contact with a real consumer. |
| Out of scope | SLSA, cosign, SBOM, fuzzing, PR/issue templates, FOSSA | Each justified individually below. |

## 3. Final repo layout

After restructure:

```
speechnorm/
├── .github/
│   ├── CODEOWNERS
│   ├── dependabot.yml
│   └── workflows/
│       ├── ci.yml
│       ├── release-please.yml
│       ├── scorecard.yml
│       └── dependency-review.yml
├── doc.go                    ← package-level godoc (NEW)
├── go.mod                    ← module github.com/olidow/speechnorm; go 1.24
├── converter.go
├── currency.go
├── normaliser.go
├── numwords_de.go
├── numwords_en.go
├── numwords_es.go
├── numwords_fr.go
├── numwords_it.go
├── numwords_pt.go
├── converter_test.go
├── currency_test.go
├── example_test.go           ← runnable godoc examples (NEW)
├── normaliser_test.go
├── numwords_*_test.go
├── numwords_*_scale_test.go
├── testhelpers_test.go
├── .golangci.yml             ← linter config (NEW)
├── CLAUDE.md                 ← updated (drop `cd internal`)
├── LICENSE                   ← existing MIT, unchanged
├── README.md                 ← NEW
└── SECURITY.md               ← NEW
```

The `internal/` directory is removed entirely.

## 4. New files in detail

### 4.1 `doc.go`

Single-file home for the package-level godoc that renders on `pkg.go.dev`:

```go
// Package speechnorm rewrites digit patterns in free-form text into
// locale-appropriate spoken words for TTS input. Supported locales:
// en, de, es, fr, it, pt. The single entry point is NormaliseNumbers.
//
// Currency words ("dollars", "and fifty cents") are always English
// regardless of locale; only the number words follow the locale's
// converter.
package speechnorm
```

### 4.2 `example_test.go`

One runnable example per locale, with `// Output:` blocks so they double as regression tests:

```go
func ExampleNormaliseNumbers_english() {
    fmt.Println(speechnorm.NormaliseNumbers("I have 3 cats", "en"))
    // Output: I have three cats
}

func ExampleNormaliseNumbers_german() {
    fmt.Println(speechnorm.NormaliseNumbers("Ich habe 3 Katzen", "de"))
    // Output: Ich habe drei Katzen
}
// … es, fr, it, pt
```

These render as runnable snippets on `pkg.go.dev` and run as part of `go test`.

### 4.3 `README.md`

Sections, in order:
1. Title + one-line description.
2. Badges row (CI status, latest release, pkg.go.dev, Codecov, Go Report Card, Scorecard, license).
3. One-paragraph what/why.
4. Install: `go get github.com/olidow/speechnorm`.
5. 5-line usage snippet (English example).
6. Supported locales table (en/de/es/fr/it/pt with one example each).
7. "Currency note" caveat (English-only currency words).
8. Contributing/issues link.
9. License line.

Kept short — long READMEs go stale faster than they help.

### 4.4 `SECURITY.md`

Three lines: link to GitHub Private Vulnerability Reporting, max-72h-acknowledgement promise, supported-versions note ("latest minor of the current major").

### 4.5 `CLAUDE.md` updates

Remove every `cd internal` reference; commands become `go test ./...` from the repo root. Update the architecture section's note about "module root in `internal/`" — that nuance disappears with the restructure.

### 4.6 `.golangci.yml`

```yaml
version: "2"
linters:
  enable:
    - errcheck
    - errorlint
    - gofumpt
    - goimports
    - govet
    - ineffassign
    - misspell
    - revive
    - staticcheck
    - unused
```

## 5. Workflows in detail

### 5.1 `.github/workflows/ci.yml`

Triggers: `push: branches: [main]`, `pull_request`.
Workflow-level: `permissions: read-all`.

Jobs (all start with `step-security/harden-runner@<sha>` in `audit` mode):

- **`test`** matrix:
  - Required: `{go: ['1.23', '1.24'], os: [ubuntu-latest]}`
  - Non-blocking: `{go: '1.24', os: [macos-latest, windows-latest]}` with `continue-on-error: true`
  - Steps: checkout → setup-go (with `go-version: ${{ matrix.go }}`, `cache-dependency-path: go.mod`) → `go build ./...` → `go test -race -shuffle=on -covermode=atomic -coverprofile=coverage.out ./...`
  - Coverage upload only on Linux/Go-1.24 job; needs `id-token: write`; `fail_ci_if_error: false`
- **`lint`**: golangci-lint-action v9 with `version:` pinned to a specific golangci-lint release (the implementation step picks the latest stable v2 at the time and records it; **never `latest`** — implicit "latest" breaks CI when golangci-lint adds rules)
- **`govulncheck`**: install + run on `./...`
- **`commitlint`** (only on `pull_request`): `wagoid/commitlint-github-action`

All third-party actions pinned to a full commit SHA with `# vX.Y.Z` trailing comment.

### 5.2 `.github/workflows/release-please.yml`

Triggers: `push: branches: [main]`.
Permissions: `contents: write`, `pull-requests: write`.
Single job: `googleapis/release-please-action@<sha> # v4` with `release-type: go`.

Behaviour: scans Conventional Commits since the last tag, opens/updates a release PR. Merging that PR cuts the tag and creates the GitHub Release.

Pre-1.0 default kept: `v0.x` breaking changes bump minor, not major.

### 5.3 `.github/workflows/scorecard.yml`

Triggers: `schedule: cron '0 6 * * 1'`, `branch_protection_rule`.
Permissions: `security-events: write`, `id-token: write`, `contents: read`.
`ossf/scorecard-action@<sha> # v2`. Publishes to GitHub code scanning, powers the Scorecard badge. Target: 8.5+.

### 5.4 `.github/workflows/dependency-review.yml`

Triggers: `pull_request`. `actions/dependency-review-action@<sha> # v4`. Marginal value at zero deps; activates the day a dep is added.

### 5.5 `.github/dependabot.yml`

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule: { interval: weekly }
  - package-ecosystem: github-actions
    directory: /
    schedule: { interval: weekly }
    groups:
      actions: { patterns: ["*"] }
```

Grouped GH Actions updates avoid drowning in PRs.

### 5.6 `.github/CODEOWNERS`

```
* @olidow
```

## 6. GitHub repo settings (manual, one-time)

- **Default branch:** `main`.
- **Merge button:** squash-merge only; "Default to PR title and description"; rebase + merge-commit disabled.
- **Branch protection on `main`:**
  - Require PR before merging — set "Required approving reviews" to **0** (solo workflow; you cannot self-approve on GitHub, but the rule still forces PRs and required status checks). Bump to 1 when external contributors arrive.
  - Required status checks: `test (1.23, ubuntu-latest)`, `test (1.24, ubuntu-latest)`, `lint`, `govulncheck`, `commitlint`.
  - Require branches up to date.
  - Require signed commits.
  - Require linear history.
  - Block force-push and deletion.
- **Repo Security tab:**
  - Enable Private Vulnerability Reporting.
  - Enable Dependabot alerts + security updates.
  - Enable secret scanning + push protection.
- **Account-level:**
  - 2FA enforced (preferably hardware key).

## 7. Local one-time setup

- Configure SSH-based git signing:
  ```
  git config --global gpg.format ssh
  git config --global user.signingkey ~/.ssh/id_ed25519.pub
  git config --global commit.gpgsign true
  git config --global tag.gpgSign true
  ```
- Register the same SSH key on GitHub under "SSH and GPG keys → New signing key" (yes — separately from the auth key entry).

## 8. First-publish ritual (do once, in order)

1. Restructure (one commit): from the repo root, run `git mv internal/* ./` (this moves `go.mod` and every `*.go` file up one level), then `rmdir internal`, then commit with message `chore: move package source to repo root`.
2. Add `doc.go`, `example_test.go`, `README.md`, `SECURITY.md`, `.golangci.yml` (separate commits, Conventional-Commit messages).
3. Wire the four workflows + Dependabot config + CODEOWNERS.
4. Push to `github.com/olidow/speechnorm` (create the repo first; choose "no template").
5. Configure repo settings (Section 6).
6. Open a PR with the changes; verify all CI checks run green; merge.
7. Wait for release-please to open the first release PR (`chore(main): release v0.1.0`); merge it.
8. Trigger module proxy fetch:
   ```
   GOPROXY=proxy.golang.org go list -m github.com/olidow/speechnorm@v0.1.0
   ```
9. Within ~10 minutes, `https://pkg.go.dev/github.com/olidow/speechnorm` should render with all docs + examples.

## 9. Things that will likely need a second attempt (first-time honest)

- **Branch-protection required-status-check names** — these have to exactly match the workflow job names *as they appear in the first successful CI run*. Set protection AFTER the first CI run completes so GitHub auto-suggests the names. Otherwise you'll typo one.
- **release-please's first PR** can look weird because it has no prior tag to anchor to. The fix is usually a `Release-As: 0.1.0` footer on one commit, or letting it default and accepting whatever it picks first.
- **Conventional Commit titles** — `commitlint` will reject things like "Update README" until you learn to type `docs: update readme`. Annoying for a day, then automatic.
- **Codecov OIDC** — sometimes works on first try, sometimes the upload step needs you to log into Codecov once to claim the repo. Keep `fail_ci_if_error: false` so this doesn't block merges.
- **pkg.go.dev indexing latency** — usually <10 minutes after the proxy fetch, but can be up to an hour. Do not panic.
- **Action SHA pinning** — Dependabot's first PR re-pins everything. Worth one careful review pass to confirm versions match comments before accepting auto-merge later.
- **Signing setup** — verify with `git commit --allow-empty -m "test signing"` and check the GitHub commit page shows "Verified" before turning on the branch-protection rule, or you'll lock yourself out of pushing.

## 10. Out of scope (justified)

- **SLSA build provenance.** Library has no maintainer-built artefact; consumers compile from source. `sum.golang.org` already provides equivalent tamper-evidence. Add the day a CLI binary ships.
- **cosign keyless signing of release artefacts.** Only artefact would be a `checksums.txt` consumers don't fetch. ROI too low.
- **SBOM generation.** Zero deps + stdlib = trivial SBOM. Add later if scope grows.
- **Fuzzing / OSS-Fuzz.** Meaningful for crypto/parser libraries, not for digit substitution with deterministic inputs.
- **PR/issue templates, FOSSA, semantic-release, multi-arch matrices, container builds, vanity import paths.** Not earning their complexity for this shape of repo.

## 11. Open questions / future revisits

- When to cut `v1.0.0` — needs at least one external consumer using the API in anger.
- Whether to add Arabic (`ar`) locale; `normaliser_test.go` already has a comment hinting at a separate spec.
- Whether to lower the `go` directive (currently 1.24) once external consumers appear; revisit when an issue is filed asking for an older Go.
- Whether to add SLSA/cosign/SBOM if the library ever ships a companion CLI tool.
