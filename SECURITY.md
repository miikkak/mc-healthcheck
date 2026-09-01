# Security Policy

## Supported Versions

Only the latest released version is supported. Security fixes are not backported to older tags.

| Version  | Supported          |
| -------- | ------------------ |
| latest   | :white_check_mark: |
| < latest | :x:                |

## Reporting a Vulnerability

Please follow these steps if you discover a security vulnerability in this project:

### Do Not

- **Do not** open a public GitHub issue for security vulnerabilities
- **Do not** disclose the vulnerability publicly until it has been addressed

### Do

1. **Report privately** via [GitHub Security Advisories](https://github.com/miikkak/mc-healthcheck/security/advisories/new) <!-- markdownlint-disable-line MD013 -->
2. **Include in your report:**
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact
   - Suggested fix (if you have one)

3. **Response timeline:**
   - You should receive an acknowledgment within 48 hours
   - We'll provide a detailed response within 7 days
   - We'll work with you to understand and fix the issue
   - We'll release a fix as soon as possible

## Scope

This connects outbound to a host/port you specify and parses whatever response comes back (a
Java Edition Server List Ping response, or a Bedrock RakNet unconnected-pong) as untrusted input

- it doesn't listen on any network port itself, doesn't execute anything derived from the
  response, and its only output is a pass/fail exit code plus, optionally, the parsed status as
  JSON on stdout. The most security-relevant code path is parsing that untrusted response,
  including length-prefixed fields (VarInts) that must be bounds-checked against integer overflow
  before use. If you find a malformed response that causes a crash, hang, excessive memory use, or
  any other misbehavior beyond a clean non-zero exit, that's exactly the kind of thing to report.

## Security Best Practices

When using this tool:

- Always use a specific released version, not a locally built development binary, in production
- Point it only at Minecraft servers/proxies you control or trust - it's designed to parse a
  well-known protocol, not to be hardened against an actively hostile server
- Keep the tool updated - check releases periodically or watch the repository

## Security Scanning

This project uses automated security scanning:

- **Trivy** (filesystem scan against `go.sum`) for dependency vulnerability scanning, on a
  weekly schedule and on demand
- **golangci-lint** (including security-relevant linters) on every PR
- **Renovate** for automated dependency updates

## Other Automated Review

Every pull request also gets an AI code review. This is a general correctness/quality review,
not a vulnerability scanner - don't rely on it as a substitute for the security scanning above.

## Disclosure Policy

- Security issues are fixed in private before public disclosure
- After a fix is released, we publish a security advisory
- We credit reporters in the advisory (unless they prefer anonymity)

## Past Security Advisories

No security advisories have been published yet.

## Contact

For security-related questions or concerns, please use the reporting method above rather than
public channels.
