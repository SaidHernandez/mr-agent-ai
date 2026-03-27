
# Skill Creator Agent

## When to Use

Create a new skill when:
- A pattern is used repeatedly and the AI needs explicit guidance
- Project conventions differ from generic best practices
- A complex workflow needs step-by-step instructions
- A decision tree would reduce ambiguity

**Don't create a skill when:**
- Documentation already exists — create a ``references/`` entry instead
- The pattern is trivial or self-explanatory
- It's a one-off task

---

## Skill Structure

```
skills/{skill-name}/
├── SKILL.md              # Required — main skill file
├── assets/               # Optional — templates, schemas, examples
│   ├── TEMPLATE.ts
│   └── schema.json
└── references/           # Optional — links to local docs
    └── docs.md
```

`assets/` = code templates, JSON schemas, example configs.
`references/` = pointers to LOCAL files (not web URLs).

---

## Decision Tree

```
Pattern reused 3+ times?            → Create a skill
Project-specific convention?        → Create a prowler-{name} or project-{name} skill
Need code templates?                → Add to assets/
Linking to existing docs?           → Use references/ with local path
Already exists in skills/?          → Update it — don't duplicate
```

---

## Naming Conventions

| Type | Pattern | Examples |
|------|---------|----------|
| Generic skill | `{technology}` | `typescript`, `playwriter` |
| Project-specific | `{project}-{component}` | `my-app-auth` |
| Workflow skill | `{action}-{target}` | `skill-creator`, `github-pr` |
| Testing skill | `test-{component}` | `test-services` |

---

## Required Frontmatter Fields

| Field | Required | Value |
|-------|----------|-------|
| `name` | Yes | Lowercase, hyphens |
| `description` | Yes | What + Trigger in one block |
| `license` | Yes | `MIT` |
| `allowed-tools` | Yes | Comma-separated tool list |
| `metadata.version` | Yes | Semantic version as string |

---

## Content Guidelines

**DO:**
- Start with the most critical patterns
- Use decision trees for ambiguous choices
- Include ✅ / ❌ code examples
- Add `Commands ` section with copy-paste bash

**DON'T:**
- Duplicate content from existing skills — reference instead
- Add lengthy prose explanations — bullet points and code
- Include troubleshooting sections
- Use web URLs in `references/` — use local paths

---

## After Creating a Skill

Add it to `AGENTS.md`:

```markdown
| `mr-agent-{skill-name}` | {trigger description} | [`skills/{skill-name}/SKILL.md`](skills/{skill-name}/SKILL.md) |
```

---

## Checklist Before Creating

- [ ] Skill doesn't already exist (checked `skills/` with Glob)
- [ ] Pattern is reusable — not a one-off
- [ ] Name follows conventions
- [ ] Frontmatter includes `allowed-tools`
- [ ] Decision tree covers the most common branching point
- [ ] ✅ / ❌ code examples are minimal and focused
- [ ] `Commands ` section is present
- [ ] `Definition of Done` checklist present
- [ ] Added to `AGENTS.md`

---

## Commands

```bash
# List existing skills
ls skills/

# Check if skill already exists
ls skills/ | grep <skill-name>
```

## Resources

- **Template**: See [assets/SKILL-TEMPLATE.md](assets/SKILL-TEMPLATE.md) for the full template
