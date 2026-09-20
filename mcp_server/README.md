# SafeLine MCP Server

SafeLine MCP Server provides read-only, multi-instance access to SafeLine CE through the Model Context Protocol. This release intentionally exposes one Tool and one SafeLine API:

| MCP Tool | SafeLine request | Effect |
|---|---|---|
| `get_attack_events` | `GET /api/open/events` | Read only |

The original site-creation, blacklist, and whitelist Tools have been removed. The downstream client has no generic write method and cannot issue SafeLine `POST`, `PUT`, `PATCH`, or `DELETE` requests.

## Protocol

- MCP SDK: `mcp-go v1.0.0-beta.1`
- Go: `1.25.14` or newer
- Transport: Streamable HTTP
- Endpoint: `POST /mcp`
- Latest supported MCP revision: `2026-07-28`
- Removed legacy routes: `/sse` and `/message`

The server supports the 2026-07-28 stateless protocol and retains Streamable HTTP compatibility with older MCP clients on the same `/mcp` endpoint.

## Configuration

Instances are deployment configuration, not MCP Tools. Each Tool call must specify an `instance_id`; the server never guesses an instance or reuses another instance's credential.

```yaml
server:
  name: "SafeLine MCP Server"
  version: "2.0.0"
  host: "127.0.0.1"
  port: 5678

logger:
  level: "info"
  file_path: ""
  console: true
  caller: false
  development: false

instances:
  - id: "production-a"
    display_name: "Production A"
    base_url: "https://10.0.0.10:9443"
    token_file: "/run/secrets/production-a.token"
    timeout: 30
    debug: false
    ca_file: ""
    insecure_skip_verify: false

  - id: "production-b"
    display_name: "Production B"
    base_url: "https://10.0.0.11:9443"
    token_file: "/run/secrets/production-b.token"
    timeout: 30
    debug: false
    ca_file: ""
    insecure_skip_verify: false
```

`LISTEN_ADDRESS` and `LISTEN_PORT` override listener settings. There is deliberately no global `SAFELINE_API_TOKEN`: every instance uses a separate token file.

`ca_file` is the certificate authority that signed the certificate of the
instance. A SafeLine console serves the certificate it generated for itself,
which no public authority signed: point `ca_file` at that certificate (or at the
authority that issued it) and the connection stays verified. Relative paths are
resolved from the directory of the configuration file.

`insecure_skip_verify` disables TLS verification for the SafeLine API. The API
Token is an administrator credential that is sent on every request, so leave the
setting at `false` and use `ca_file` for an instance with a self-signed
certificate. Setting it to `true` logs a warning at startup, and it cannot be
combined with `ca_file`.

Each `display_name` must be unique, ignoring letter case, and must not match another instance's `id`. During server discovery (or legacy initialization), the server publishes only the `display_name` to `instance_id` mappings as Server Instructions. This lets a user refer to a friendly name while the AI still calls `get_attack_events` with the stable, explicit `instance_id`. Instance addresses and credentials are never included in those instructions.

SafeLine CE does not yet provide a dedicated read-only service credential for all required management APIs. This release therefore uses an administrator API Token per instance, while the MCP process enforces the fixed read-only API allowlist.

## AI client authentication

MCP client authentication is independent from the downstream SafeLine API Tokens.

Authentication is on by default: the server always installs the bearer gate
unless a deployment explicitly opts out, and it refuses to start without a token
on any non-loopback listener.

### Initialization

Generate the MCP client Bearer Token once during deployment initialization:

```bash
docker compose run --rm mcp_server auth init
```

The command uses `crypto/rand` to generate a 256-bit token. It displays the plaintext once and stores only its SHA-256 hash and creation time in `state/auth.json` with mode `0600`. A normal service restart only reads that state and never changes or prints the token.

Clients send the generated token on every MCP request:

```http
Authorization: Bearer slmcp_...
```

The built-in listener is plain HTTP. Docker Compose publishes it on
`127.0.0.1` by default so the reusable Bearer Token does not cross a network
in cleartext. For access from another machine, put the endpoint behind a
trusted HTTPS-terminating gateway or keep it on an equivalently protected
private transport. Do not set `MCP_PUBLISH_ADDRESS=0.0.0.0` on an untrusted
network without TLS termination.

The server refuses to start while this state is missing, so a deployment cannot
serve any request before the token exists. Explicit rotation is available
locally:

```bash
docker compose run --rm mcp_server auth rotate
```

After the service restarts, rotation invalidates the previous token. This
static Bearer gate is intended for self-hosted deployment; it is not presented
as a complete MCP OAuth 2.1 implementation.
Run `auth rotate` as a single serialized deployment operation; concurrent
rotation commands are not supported.

### Disabling authentication

`MCP_AUTH_DISABLED=true` (or `MCP_AUTH_ENABLED=false`) removes the gate. The
server then refuses to start unless the listener is loopback-only, so a disabled
gate cannot silently publish attack telemetry to the network:

```text
MCP_AUTH_DISABLED=true
```

## Docker deployment

1. Build the image. The included `.dockerignore` excludes `secrets/` and
   `state/` from the build context even if local fallback directories exist:

   ```bash
   docker compose build
   ```

2. Prefer protected deployment directories outside the source checkout, then
   create one SafeLine token file per instance:

   ```bash
   install -d -m 700 /etc/safeline-mcp/secrets /var/lib/safeline-mcp
   install -m 600 /path/to/production-a.token /etc/safeline-mcp/secrets/production-a.token
   install -m 600 /path/to/production-b.token /etc/safeline-mcp/secrets/production-b.token
   export SAFELINE_MCP_SECRETS_DIR=/etc/safeline-mcp/secrets
   export SAFELINE_MCP_STATE_DIR=/var/lib/safeline-mcp
   ```

   For local-only development, `./secrets` and `./state` are accepted and are
   excluded by both `.gitignore` and `.dockerignore`; create them with the same
   permissions before running Compose.

3. Edit `config.yaml` so each `token_file` matches its file under `/run/secrets`.

4. Run the one-time `auth init` command and save the displayed token. The
   service refuses to start until this state exists.

5. Start the service:

   ```bash
   docker compose up -d
   ```

Compose refuses to auto-create missing bind-mount paths. Its host port is
loopback-only unless `MCP_PUBLISH_ADDRESS` is explicitly changed.

## Running with Go

```bash
go test ./...
go build -o mcp-server .
./mcp-server auth init --state-file "$PWD/state/auth.json"
./mcp-server --config ./config.yaml
```

The default configuration listens on `127.0.0.1`, which is the only address the
server accepts while authentication is disabled:

```bash
MCP_AUTH_DISABLED=true ./mcp-server --config ./config.yaml
```

## Tool contract

`get_attack_events` accepts:

| Field | Type | Required | Constraint |
|---|---|---|---|
| `instance_id` | string | Yes | Must match a configured instance |
| `page` | integer | No | Default `1`, minimum `1` |
| `page_size` | integer | No | Default `10`, range `1..100` |
| `ip` | string | No | Source IP, CIDR, or partial-IP filter |
| `host` | string | No | Protected host filter |
| `port` | integer | No | Range `1..65535` |
| `start` | integer | No | Unix milliseconds, non-negative |
| `end` | integer | No | Unix milliseconds, non-negative and not before `start` |

Unknown input fields and fractional integers are rejected by server-side JSON Schema validation. Successful results include both text fallback content and `structuredContent` with:

- `source.instance_id`
- `source.display_name`
- `page` and `page_size`
- `nodes`
- `total`

SafeLine HTTP failures, HTTP-200 envelopes with a non-empty `err`, unknown instances, and input validation failures are returned as MCP Tool Execution Errors with `result.isError=true`. The MCP server does not interpret human-readable `msg` content when `err` is empty or absent.
