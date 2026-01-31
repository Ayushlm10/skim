# AGENTS.md Distribution System - Specification

## What We're Building

An automated system that generates high-quality AGENTS.md files for every eligible repository, delivered via merge requests that repo owners review and merge.

---

## Why

AI coding tools (Claude Code, opencode, Cursor) work significantly better when they understand the codebase they're operating in. Without context, they make assumptions—wrong package manager, wrong file locations, wrong patterns, wrong commands. These mistakes waste developer time and erode trust in AI tooling.

**AGENTS.md** is a context file that AI tools read before working on code. It tells them how this specific repo works: what commands to run, where files go, what patterns to follow, what to avoid.

**If we do this well:**
- AI tools produce correct code on the first try more often
- Developers spend less time correcting AI mistakes
- AI tooling adoption increases because it actually work

**If we do this poorly:**
- AI tools are no better than before
- Developers lose trust and stop using AI tooling
- We've created maintenance burden with no value

---

## Success Criteria

### Primary Metric: MR Merge Rate

| Outcome | Interpretation |
|---------|----------------|
| >70% merge rate | Files are accurate and valuable |
| 40-70% merge rate | Files need improvement |
| <40% merge rate | Approach is flawed |

### Secondary Metrics

- **Time-to-merge**: Lower = higher confidence
- **Modification rate**: High = we're close but not perfect
- **Developer feedback**: Qualitative insights
- **Migration merge rate**: Track separately from new generation (may differ)

---

## User Journey

### Repo Maintainer Receives MR

1. Gets MR titled "Add AGENTS.md for AI Coding Tools"
2. Reviews generated content - is it accurate?
3. Edits if needed (they know their repo best)
4. Merges if useful, closes with feedback if not

### Developer Using AI Tools

1. Opens repo in AI coding tool
2. Tool reads AGENTS.md automatically
3. Tool understands how to build, test, and work in this repo
4. Developer gets correct suggestions from first interaction

---

## Eligibility Criteria

A repository is eligible if ALL of the following are true:

| Criterion | Requirement | Rationale |
|-----------|-------------|-----------|
| Activity | Commit within last 6 months | Active repos where AI tooling matters |
| No AGENTS.md | File doesn't exist | Don't overwrite existing work |
| Has code | Not empty/placeholder | Skip trivial repos |
| Not archived | Active repository | No point in archived repos |
| Minimum size | >5 files | Skip trivial repos |

**Review candidates:** Repos with activity 6-12 months ago are flagged for manual review.

### Repository Categories

Based on eligibility checks, repos fall into three categories:

| Category | Condition | Action |
|----------|-----------|--------|
| **Generate** | No AGENTS.md, no CLAUDE.md | Generate new AGENTS.md |
| **Migrate** | No AGENTS.md, has CLAUDE.md | Migrate CLAUDE.md → AGENTS.md |
| **Skip** | Has AGENTS.md | Already done |

---

## What Makes a Good AGENTS.md

Based on research across 60,000+ repositories:

### Constraints
- **Max 100 lines / 13KB** - Instruction-following degrades beyond this
- **Verified commands only** - Never hallucinate; omit if unsure
- **Specific over generic** - "use pnpm" not "use the package manager"

### Required Sections

| Section | Purpose |
|---------|---------|
| Project Overview | One sentence + tech stack |
| Commands | Exact build/test/lint commands |
| Boundaries | Always/Ask First/Never rules |

### Anti-Patterns

- Generic advice ("write clean code")
- Linting rules (use lint command instead)
- Detailed file paths (they change)
- Comprehensive docs (keep it focused)

---

## Rollout Plan

### Wave 1: Pilot (devtools group)
- Validate approach with friendly team
- Gather feedback, iterate
- Target: >50% merge rate

### Wave 2: Expansion
- Roll out to additional groups based on pilot learnings

### Wave 3: Full Rollout
- All eligible repositories

---

## Maintenance Strategy

After initial merge, AGENTS.md files need to stay current. We use a phased approach that evolves from passive suggestions to active agentic updates.

### Phase 1: AI Reviewer Suggestions (MVP)

**What**: OpenCode AI reviewer detects relevant changes in MRs and suggests AGENTS.md updates as comments.

**Triggers for suggestions:**
- New scripts added to package.json / Makefile
- Changed build/test commands
- New directories that affect project structure
- Dependency changes that affect tooling

**Output**: MR comment suggesting what to add/update in AGENTS.md.

| Pros | Cons |
|------|------|
| Low friction - just suggestions | Suggestions can be ignored |
| No new infra needed | Manual work to apply changes |
| Works with existing workflow | Doesn't scale to everyone |

**Good for**: Initial rollout, proving the concept works.

### Phase 2: @agent Bot Trigger (Future)

**What**: Users @ mention a bot in MR comments to trigger agentic actions.

**Example usage:**
```
@opencode-agent update AGENTS.md
@opencode-agent add command "pnpm test:e2e" to AGENTS.md
```

The bot sees the MR context, understands what changed, and makes inline updates.

**Why this matters:**
- **Habit formation** - Developers learn to work with AI agents
- **Pull model** - User requests > push model (less noise)
- **Context-aware** - Bot sees the MR, understands what changed
- **Democratized** - Anyone can trigger, not just maintainers
- **Extensible** - Same pattern works for other agentic tasks

| Pros | Cons |
|------|------|
| User-initiated = high intent | Requires new infra (bot service) |
| Teaches agentic thinking | Adoption curve |
| Scales to everyone | Needs permissions setup |
| Low friction for updates | Scope creep risk |

**Guardrails**: Initially, bot only modifies AGENTS.md. Broader capabilities added later.

### Phase 3: Proactive Auto-Updates (Future)

**What**: System detects significant changes and auto-creates MRs to update AGENTS.md.

**Only for**: Major changes (new build system, new test framework, etc.) - not every small change.

### Progression

```
Phase 1: "Hey, you should update AGENTS.md" (passive suggestion)
    ↓
Phase 2: "@opencode-agent update AGENTS.md" (user-triggered action)
    ↓
Phase 3: Auto-PR for significant changes (proactive, opt-in)
```

---

## CLAUDE.md Migration Strategy

AGENTS.md is the emerging industry standard for AI coding tool context files. Tools like Windsurf, Cursor, and others are converging on this format. To future-proof our repos while maintaining backward compatibility with Claude Code, we migrate CLAUDE.md → AGENTS.md.

### Migration Process

For repos with existing CLAUDE.md (but no AGENTS.md):

1. **Copy**: Create AGENTS.md with CLAUDE.md content as base
2. **Enhance** (strictly limited): Only add verified commands/info that CLAUDE.md missed
3. **Symlink**: Replace CLAUDE.md with symlink → AGENTS.md
4. **MR**: Create merge request with clear explanation

### Enhancement Rules (Strict)

| Allowed | Not Allowed |
|---------|-------------|
| ✅ Add verified commands CLAUDE.md doesn't mention | ❌ Rewrite existing good content |
| ✅ Add tech stack detection if CLAUDE.md is vague | ❌ Change tone/style significantly |
| ✅ Format improvements for readability | ❌ Remove anything from CLAUDE.md |

### Why Symlink?

- **Backward compatible**: Claude Code follows symlinks, continues working
- **Single source of truth**: One file to maintain
- **Industry alignment**: AGENTS.md is the standard going forward
- **Zero breakage**: Existing workflows unchanged

### MR Template for Migration

```markdown
## 🔄 Migrate CLAUDE.md → AGENTS.md

### Why?
AGENTS.md is the emerging industry standard for AI coding tool context files 
(Windsurf, Cursor, Claude Code, etc.). This migration ensures better tool 
compatibility and future-proofing.

### What changed?
- ✅ Created AGENTS.md (content preserved from CLAUDE.md)
- ✅ CLAUDE.md → symlink to AGENTS.md (backward compatible)
- ✅ [List any enhancements, if applicable]

### No functionality lost
Claude Code and other tools that read CLAUDE.md will continue working.

---
Feel free to close this MR if you prefer to keep CLAUDE.md standalone.
```

---

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Low merge rate | Start small, iterate based on feedback |
| Inaccurate generation | Human review required, validation layer |
| MR fatigue | Clear messaging, demonstrate value |
| Files go stale | AI reviewer suggests updates |
| Symlink breakage | Test on CI systems first; most tools handle symlinks correctly |
| Migration perceived as churn | Clear MR explanation + user choice to close |
| Over-enhancement of CLAUDE.md | Strict rules: only add verified, never remove |

---

## For Technical Implementation

See `plan.md` for technology stack, architecture, and implementation tasks.
