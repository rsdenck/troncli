# Security Policy

## Supported Versions

Only the latest stable release of `troncli` is supported for security updates.

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**DO NOT** report security vulnerabilities through public GitHub issues.

If you believe you have found a security vulnerability in `troncli`, please report it directly to the maintainers via email at: `security@example.com` (Replace with actual email).

### Process

1.  **Report**: You report the vulnerability privately.
2.  **Ack**: We acknowledge receipt within 48 hours.
3.  **Investigate**: We investigate the issue and determine the impact.
4.  **Fix**: We develop a fix and test it.
5.  **Release**: We release a patched version.
6.  **Disclose**: We publicly disclose the vulnerability after the fix is available.

## Supply Chain Security

We take supply chain security seriously.
- All releases are signed with Cosign.
- SBOMs are generated for every release (CycloneDX).
- Dependencies are scanned nightly.
- CI/CD pipelines are pinned by SHA.
