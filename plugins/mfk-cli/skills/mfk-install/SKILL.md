---
name: mfk-install
description: Make the mfk command available, installing it only if it is missing. Use when another skill reports that `mfk` is not on PATH, or when the user asks to install, update or upgrade the ideamans Money Forward Kessai CLI. Prefers an already-installed binary, then the latest GitHub release, then a build from source with go install.
license: MIT
compatibility: Requires curl (or wget) and tar to install from a release, or a Go toolchain for the source fallback. Standalone — does not need mfk to be present already. Installs from the public repository github.com/ideamans/mfk-cli, so no GitHub authentication is needed.
allowed-tools: Bash(curl:*) Bash(wget:*) Bash(tar:*) Bash(unzip:*) Bash(go:*) Bash(uname:*) Bash(command:*) Bash(which:*) Bash(mkdir:*) Bash(mv:*) Bash(cp:*) Bash(rm:*) Bash(chmod:*) Bash(ls:*) Bash(test:*) Bash(echo:*) Read
---

# mfk-install

Make the `mfk` command usable, doing the least work that achieves it.

## Route 1 — an existing installation on PATH

```bash
command -v mfk && mfk --version
```

If that resolves, **use it and stop here.** Do not check for a newer release —
it costs an API call and the user did not ask for an upgrade.

Two checks before trusting the hit:

- **It is the right tool.** `mfk` is a short name. `mfk llm | head -1` must read
  `# mfk CLI — AIエージェント向けリファレンス`. If something else owns the name,
  tell the user and use an explicit path rather than shadowing theirs.
- **It is recent enough.** If `mfk llm` is not a known command, the binary
  predates the embedded reference. Say so and continue to route 2.

Continue past this section only when the command is missing, is the wrong tool,
is too old, or the user explicitly asked to update.

## Route 2 — the latest GitHub release

```bash
VERSION=$(curl -fsSL https://api.github.com/repos/ideamans/mfk-cli/releases/latest \
  | grep '"tag_name"' | head -1 | cut -d'"' -f4)   # e.g. v0.3.0

OS=$(uname -s | tr '[:upper:]' '[:lower:]')            # darwin | linux
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64  # amd64 | arm64
curl -fsSL -o /tmp/mfk.tar.gz \
  "https://github.com/ideamans/mfk-cli/releases/download/${VERSION}/mfk-cli_${VERSION#v}_${OS}_${ARCH}.tar.gz"
```

**The archive is named `mfk-cli_…` but the binary inside is `mfk`** — without the
`-cli` suffix. Windows ships a `.zip`.

If the download 404s, list the actual assets on the release page rather than
retrying variations.

### Install onto PATH

```bash
tar -xzf /tmp/mfk.tar.gz -C /tmp
mkdir -p ~/.local/bin && mv /tmp/mfk ~/.local/bin/ && chmod +x ~/.local/bin/mfk
```

Prefer the first writable directory already on PATH — `~/.local/bin`, then
`/usr/local/bin`. Two things not to do on your own initiative:

- If nothing on PATH is writable, leave the binary in `/tmp`, print the exact
  `sudo mv` command and let the user run it. Do not run `sudo` yourself.
- If `~/.local/bin` is not on PATH, give the user the line to add to their shell
  profile. Do not edit the profile for them.

## Route 3 — build from source

Needs a Go toolchain and compiles rather than downloads. The module root is the
main package here, so no path suffix is needed.

```bash
go install github.com/ideamans/mfk-cli@latest
```

The binary lands in `$(go env GOPATH)/bin` and is named **`mfk-cli`**, not `mfk`
— `go install` names it after the module. Either symlink it or tell the user
which name to invoke; do not leave them wondering why `mfk` is still missing.

## Verify

```bash
mfk --version
mfk llm | head -5
```

Report which route was taken, the version and the install path.

Then say what is still needed: `mfk` cannot call anything without `MFK_API_KEY`
(production) or `MFK_SANDBOX_API_KEY` with `--sandbox`. Mention that the key
decides which environment real billing data comes from, so it is worth setting
deliberately rather than copying whatever is at hand.
