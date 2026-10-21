# hook-email

[![Latest Release](https://img.shields.io/github/v/release/SemRels/hook-email?label=version&color=blue)](https://github.com/SemRels/hook-email/releases/latest)
[![CI](https://github.com/SemRels/hook-email/actions/workflows/ci.yml/badge.svg)](https://github.com/SemRels/hook-email/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/SemRels/hook-email)](LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/SemRels/hook-email/badge)](https://scorecard.dev/viewer/?uri=github.com/SemRels/hook-email)
[![REUSE status](https://api.reuse.software/badge/github.com/SemRels/hook-email)](https://api.reuse.software/info/github.com/SemRels/hook-email)

Sends a release notification email with the new version and changelog details.

This plugin is distributed as the standalone Go binary `semrel-plugin-hook-email`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/hook-email/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/hook-email:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/hook-email:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/hook-email/.github/workflows/container.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - uses: "@semrel/hook-email"
    phase: release
    args:
      smtp_host: "smtp.example.com"
      smtp_port: "587"
      smtp_user: "semrel-bot"
      smtp_pass: "${SMTP_PASSWORD}"
      from: "semrel@example.com"
      to: "team@example.com,ops@example.com"
      subject_template: "Release {{ .TagName }} is live"
      tls: "true"
```

## Options

Every option can be set as an entry under `args` in `.semrel.yaml`. semrel passes each entry to the plugin as the environment variable shown in the table, so you can also set that variable directly (for example in your CI environment) instead of using `args`.

| Argument | Environment variable | Required | Description | Default |
| --- | --- | --- | --- | --- |
| `smtp_host` | `SEMREL_PLUGIN_SMTP_HOST` | Required | SMTP server hostname. | None |
| `smtp_port` | `SEMREL_PLUGIN_SMTP_PORT` | Optional | SMTP server port. | 587 |
| `smtp_user` | `SEMREL_PLUGIN_SMTP_USER` | Required | SMTP username used for authentication. | None |
| `smtp_pass` | `SEMREL_PLUGIN_SMTP_PASS` | Required | SMTP password used for authentication. | None |
| `from` | `SEMREL_PLUGIN_FROM` | Required | From address for the release email. | None |
| `to` | `SEMREL_PLUGIN_TO` | Required | Comma-separated list of recipient email addresses. | None |
| `subject_template` | `SEMREL_PLUGIN_SUBJECT_TEMPLATE` | Optional | Custom subject template for the release email. | Built-in subject |
| `tls` | `SEMREL_PLUGIN_TLS` | Optional | Enable or disable TLS for the SMTP connection. | true |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` | Resolved release version for the current run. |
| `SEMREL_TAG_NAME` | Git tag name semrel will create or publish. |
| `SEMREL_NEXT_VERSION` | Next version computed by semrel for the release. |
| `SEMREL_CHANGELOG` | Generated changelog text for the release. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

## Example behavior

On a real release, the plugin sends an email that includes the new version, tag name, and changelog. In dry-run mode it reports what would be sent without sending the message.

## License

Apache-2.0
