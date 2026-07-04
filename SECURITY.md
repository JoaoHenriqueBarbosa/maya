# Security Policy

## Supported Versions

Maya is an exploratory prototype and is **not under active development**. There are no tagged
releases yet; only the latest `main` receives any attention.

| Version | Supported          |
|---------|--------------------|
| `main`  | :white_check_mark: |
| older   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it **privately**. Do not open a public
issue.

- **Email:** joaohenriquebarbosa21@gmail.com
- Please include enough detail to reproduce the issue (affected code path, a proof of concept,
  and the impact you believe it has).
- We aim to **acknowledge your report within 72 hours**.

## Process

1. You report the vulnerability privately via the email above.
2. We acknowledge receipt within 72 hours.
3. We investigate, confirm the issue, and assess its impact and scope.
4. We prepare and test a fix (or, given the prototype status, document a mitigation).
5. We disclose the issue and credit you, unless you prefer to remain anonymous.

## Scope

Because Maya is a Go library that compiles to WebAssembly and runs in the browser, the kinds of
issues most relevant here include:

- **Injection or unsafe input handling** — anything that lets untrusted data escape into the DOM
  or the WASM/JS bridge in an unintended way.
- **Data or credential exposure** — sensitive information leaking through the code, logs, or the
  compiled artifact.
- **Dependency vulnerabilities** — Maya intentionally ships with no third-party dependencies
  (standard library only), so this surface is minimal, but reports about the Go toolchain
  interaction or the bundled `wasm_exec.js` shim are welcome.

Please note that the debug logging currently present in the code (see the README's
[Honest status](README.md#honest-status)) is a known issue and is being tracked separately from
security reports.
