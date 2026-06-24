# hook-email

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
  --certificate-identity-regexp 'https://github.com/SemRels/hook-email/.github/workflows/release.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: hook-email
    path: ~/.semrel/plugins/semrel-plugin-hook-email
    env:
      SEMREL_PLUGIN_SMTP_HOST: "smtp.example.com"
      SEMREL_PLUGIN_SMTP_PORT: "587"
      SEMREL_PLUGIN_SMTP_USER: "semrel-bot"
      SEMREL_PLUGIN_SMTP_PASS: "${SMTP_PASSWORD}"
      SEMREL_PLUGIN_FROM: "semrel@example.com"
      SEMREL_PLUGIN_TO: "team@example.com,ops@example.com"
      SEMREL_PLUGIN_SUBJECT_TEMPLATE: "Release {{ .TagName }} is live"
      SEMREL_PLUGIN_TLS: "true"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_SMTP_HOST` | Required | SMTP server hostname. | None |
| `SEMREL_PLUGIN_SMTP_PORT` | Optional | SMTP server port. | 587 |
| `SEMREL_PLUGIN_SMTP_USER` | Required | SMTP username used for authentication. | None |
| `SEMREL_PLUGIN_SMTP_PASS` | Required | SMTP password used for authentication. | None |
| `SEMREL_PLUGIN_FROM` | Required | From address for the release email. | None |
| `SEMREL_PLUGIN_TO` | Required | Comma-separated list of recipient email addresses. | None |
| `SEMREL_PLUGIN_SUBJECT_TEMPLATE` | Optional | Custom subject template for the release email. | Built-in subject |
| `SEMREL_PLUGIN_TLS` | Optional | Enable or disable TLS for the SMTP connection. | true |

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
