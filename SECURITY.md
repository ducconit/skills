# Security Policy

## Authorized use only

These skills are for assessing systems you own or are explicitly authorized to test. Do not use them against third-party systems or without a documented scope. Keep findings, PoC, and evidence internal; never publish them.

## Reporting a vulnerability in this repository

Please report security issues privately to the repository owner through GitHub's private vulnerability reporting or a direct message. Do not open a public issue for an unfixed vulnerability, and do not include a working exploit in a public report.

Include, if available:

- affected file or skill and version/commit
- impact and an attacker precondition
- minimal reproduction with sanitized evidence (no real secrets or personal data)
- suggested fix

You will receive an acknowledgement once the report is triaged. Do not disclose the issue publicly until a fix is available and the owner agrees.

## Handling findings produced by the pentest skill

- Redact tokens, credentials, private keys, and personal data before storing or sharing evidence.
- Keep the generated `references/pentest.sqlite3`, reports, and evidence out of git.
- Use the unified gate (`gate-check`): no open Critical/High finding before a target moves forward. The standard is the same for every target.
