# Future Ideas

Ideas and features to potentially implement later.

---

## MR Status Refresh

**Problem:** Once MRs are created, we need a way to check if they've been merged/closed.

**Proposed Solution:** CLI-based refresh command

```bash
# Refresh a single repo's MR status
agents-md refresh --repo cloudflare/devtools/opencode

# Refresh all repos with open MRs
agents-md refresh --all-open

# Refresh and sync to D1
agents-md refresh --all-open && agents-md sync
```

**Implementation:**
1. Query repos with `status = 'mr_created'` from local SQLite
2. For each, call GitLab API to get MR state (`merged`, `closed`, `opened`)
3. Update local DB with new status (`mr_merged`, `mr_closed`, or keep `mr_created`)
4. User runs `sync` to push changes to D1 for dashboard

**Why CLI instead of Worker:**
- Keeps GitLab credentials local (not stored in Workers secrets)
- Dashboard/Worker is temporary - will be removed after rollout
- Simpler architecture

**Status:** Not implemented. Revisit when Phase 3 (Delivery) is complete and MRs are being created.

---
