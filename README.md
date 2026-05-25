# hook-email

Email hook plugin for Semantic Release.

Publishes Semantic Release notifications through email.

## Documentation

- Docs (coming soon): <https://github.com/SemRels/semrel/tree/main/docs/plugins/hook-email>
- Template source: <https://github.com/SemRels/plugin-template>

## Repository Layout

`	ext
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
`

## Development

`ash
go build ./cmd/plugin
go test ./...
`

## Configuration Example

`yaml
plugins:
  - name: hook-email
    type: hook
    config:
      smtp_host: smtp.example.com
      smtp_port: 587
      recipients:
        - releases@example.com
`

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.