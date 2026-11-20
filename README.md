# Agent Skills

Reusable skills for coding and security agents. Each directory under `skills/` is one self-contained skill with a `SKILL.md` entry point and optional focused references or scripts.

## Skills

| Skill | Purpose |
|---|---|
| [`pentest`](skills/pentest/SKILL.md) | Authorized, scoped penetration testing and vulnerability reporting |

## Use

Copy the selected skill directory into the target agent's supported skills directory. Keep the skill's authorization and scope requirements intact. Provider-specific installation steps may differ.

## Contributing

- Use lowercase kebab-case directory names.
- Put invocation metadata and core workflow in `SKILL.md`.
- Add references/scripts only when they materially improve repeatability.
- Keep procedures safe, scoped, evidence-based, and clear about authorization.
- Do not include credentials, target-specific secrets, or private engagement data.
