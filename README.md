# speechnorm

[![CI](https://github.com/OliDow/speechnorm/actions/workflows/ci.yml/badge.svg)](https://github.com/OliDow/speechnorm/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/OliDow/speechnorm?label=release)](https://github.com/OliDow/speechnorm/releases)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/olidow/speechnorm.svg)](https://pkg.go.dev/github.com/olidow/speechnorm)
[![Codecov](https://codecov.io/gh/OliDow/speechnorm/branch/main/graph/badge.svg)](https://codecov.io/gh/OliDow/speechnorm)
[![Go Report Card](https://goreportcard.com/badge/github.com/olidow/speechnorm)](https://goreportcard.com/report/github.com/olidow/speechnorm)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/OliDow/speechnorm/badge)](https://scorecard.dev/viewer/?uri=github.com/OliDow/speechnorm)
[![License](https://img.shields.io/github/license/OliDow/speechnorm)](LICENSE)

Locale-aware number-to-words rewriter for TTS input. Converts digit
patterns in free-form text — currency, ordinals, comma-grouped integers,
plain integers — into spoken words for the target locale.

Zero non-stdlib dependencies.

## Install

Requires Go 1.24+.

```sh
go get github.com/olidow/speechnorm
```

## Usage

```go
import "github.com/olidow/speechnorm"

speechnorm.NormaliseNumbers("I paid $5 for 3 items", "en")
// → "I paid five dollars for three items"
```

## Supported locales

| Locale | Example input | Example output |
|---|---|---|
| `ar` | `لدي 3 قطط` | `لدي ثلاثة قطط` |
| `en` | `I have 3 cats` | `I have three cats` |
| `de` | `Ich habe 3 Katzen` | `Ich habe drei Katzen` |
| `es` | `Tengo 3 gatos` | `Tengo tres gatos` |
| `fr` | `J'ai 3 chats` | `J'ai trois chats` |
| `it` | `Ho 3 gatti` | `Ho tre gatti` |
| `pt` | `Tenho 3 gatos` | `Tenho três gatos` |

Unknown or empty locales return the input unchanged.

## Currency note

Currency words (`dollars`, `and fifty cents`) are always English regardless
of the locale. Only the *number words* respect the locale's converter.
Supported symbols: `$ € £ ¥ ₹`.

## Contributing

Issues and pull requests welcome at
https://github.com/OliDow/speechnorm/issues. Commit messages must follow
[Conventional Commits](https://www.conventionalcommits.org).

## License

MIT — see [LICENSE](LICENSE).
