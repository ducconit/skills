# Agent Skills

Reusable skills for coding and security agents. Each directory under `skills/` is one self-contained skill with a `SKILL.md` entry point and optional focused references or scripts.

## Skills

| Skill | Purpose |
|---|---|
| [`pentest`](skills/pentest/SKILL.md) | Authorized, non-destructive penetration testing, local finding store, and evidence-based reporting |

The pentest skill applies one unified security baseline to every target. It has no "weaker environment" tier: an internal prototype and a production system are held to the same severity and release gate. A descriptive `context_note` may change test intensity (for example stricter rate limits or passive-only review), never the security standard.

### `skills/pentest` layout

```text
skills/pentest/
  SKILL.md                    # triggers, inputs/outputs, workflow, guardrails, HITL, reporting
  scripts/
    pentest_db.py             # stdlib-only SQLite CLI: targets, runs, endpoints, findings, patterns, gate-check
  references/
    schema.sql                # database schema
    patterns_seed.sql         # seed patterns (BOLA, BFLA, injection, SSRF, ...)
    known-patterns.md         # human-readable pattern catalogue
    report-template.md        # finding and report template
    cvss-guide.md             # severity and CVSS guidance
```

The SQLite file (`references/pentest.sqlite3`) is generated locally and gitignored. Only code and seed data are versioned.

## Quick start

```bash
python3 skills/pentest/scripts/pentest_db.py --db /tmp/pentest-demo.sqlite3 init
python3 skills/pentest/scripts/pentest_db.py --db /tmp/pentest-demo.sqlite3 pattern-search "BOLA"
python3 skills/pentest/scripts/pentest_db.py --db /tmp/pentest-demo.sqlite3 gate-check --target-id <id>
```

Use `--db` for an isolated local database. Never commit a database, report, or evidence file.

## Use

Copy the selected skill directory into the target agent's supported skills directory. Keep the skill's authorization and scope requirements intact. Provider-specific installation steps may differ.

## Safety

- Only test assets the user owns or is explicitly authorized to test; no third-party systems.
- Non-destructive by default: no data modification, no DoS, no auto-exploit, no reverse shell, no cloud metadata token retrieval.
- Human confirmation is required before switching from passive review to active probing.
- Evidence is redacted before storage; secrets, personal data, and PoC are never published.
- See [SECURITY.md](SECURITY.md) for private vulnerability reporting.

## Contributing

- Use lowercase kebab-case directory names.
- Put invocation metadata and core workflow in `SKILL.md`.
- Add references/scripts only when they materially improve repeatability.
- Keep procedures safe, scoped, evidence-based, and clear about authorization.
- Do not include credentials, target-specific secrets, or private engagement data.
- Keep generated databases, reports, and evidence out of git (see [.gitignore](.gitignore)).

## License

MIT with an authorized-use legal disclaimer — see [LICENSE](LICENSE).
