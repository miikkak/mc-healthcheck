# mc-healthcheck

A minimal health-check CLI for Minecraft servers, meant for use in container
`HEALTHCHECK` directives.

This is an independent implementation, not a fork of
[itzg/mc-monitor](https://github.com/itzg/mc-monitor) — inspired by its use
case (checking whether a Minecraft server is alive) but built from scratch to
do only that. `mc-monitor` also ships Prometheus, OpenTelemetry, and Telegraf
exporters; this tool has exactly two subcommands and no export backends.

## About this project

This was built with heavy Claude Code assistance — most of the implementation
is AI-generated, with the design and review driven by me. It has unit test
coverage (see Testing below) and runs as the container `HEALTHCHECK` for my
own production Minecraft server stack ([mc-server-container](https://github.com/miikkak/mc-server-container)),
so it sees real day-to-day use, not just its own test suite. Read the source
and file issues if something looks off.

## Usage

```shell
mc-healthcheck status --host localhost --port 25565 --timeout 5s
mc-healthcheck status-bedrock --host localhost --port 19132 --timeout 5s
```

Both subcommands exit `0` if a well-formed response comes back within the
timeout, non-zero otherwise. Nothing is printed on success; a message is
printed to stderr on failure.

### `status`

Performs a Java Edition Server List Ping (handshake + status request) over
TCP and confirms the response is well-formed JSON. Works against both Paper
servers and Velocity proxies, since both respond to the vanilla SLP protocol.

With `--json`, prints the decoded status response (player count, MOTD,
version, etc.) to stdout on success instead of staying silent:

```shell
mc-healthcheck status --host localhost --port 25565 --json
# {"players":{"max":20,"online":3,...},"version":{"name":"1.21",...},...}
```

| Flag        | Default     | Description                                         |
| ----------- | ----------- | --------------------------------------------------- |
| `--host`    | `localhost` | Server hostname                                     |
| `--port`    | `25565`     | Server port                                         |
| `--timeout` | `5s`        | Connection and read timeout                         |
| `--json`    | `false`     | Print the decoded status response as JSON to stdout |

### `status-bedrock`

Performs a RakNet unconnected ping over UDP, for a Bedrock listener (e.g. a
Geyser plugin bridging Bedrock clients into a Velocity proxy).

| Flag        | Default     | Description     |
| ----------- | ----------- | --------------- |
| `--host`    | `localhost` | Server hostname |
| `--port`    | `19132`     | Server port     |
| `--timeout` | `5s`        | Read timeout    |

## Building

```shell
make build
```

Produces an `mc-healthcheck` binary in the repo root. `make install` installs
it to `/usr/local/lib/helpers` (override with `PREFIX`/`DESTDIR`).

## Testing

```shell
go test ./... -race -cover
```
