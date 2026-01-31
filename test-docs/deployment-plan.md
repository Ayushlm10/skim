# Deployment Plan: Task 1.11 - Dashboard Deployment

## Overview

Deploy the AGENTS.md Dashboard to Cloudflare Workers with D1 database, serving both the API and React frontend from a single worker behind Cloudflare Access.

**Target URL:** `https://agents-md.devtools.cfdata.org`  
**Account:** DevTools: Dev Environments (`2469e0bba8bf4d732f65a093985146d6`)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  agents-md.devtools.cfdata.org               │
│                    (Cloudflare Access)                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   ┌─────────────────────────────────────────────────────┐   │
│   │              Cloudflare Worker                       │   │
│   │         agents-md-dashboard                          │   │
│   │                                                      │   │
│   │   /api/*  ────▶  API Handlers  ────▶  D1 Database   │   │
│   │                                                      │   │
│   │   /*      ────▶  Static Assets (React SPA)          │   │
│   └─────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Pre-Deployment Checklist

- [x] Wrangler installed (v4.59.2)
- [x] Authenticated to Cloudflare (athakur@cloudflare.com)
- [x] DevTools account selected
- [x] Frontend built (`dashboard/frontend/dist/`)
- [x] Worker modified for static asset serving
- [x] wrangler.jsonc updated with deployment config (account_id, routes, assets)
- [ ] D1 database created
- [ ] Schema applied to D1
- [ ] wrangler.jsonc updated with actual database_id
- [ ] Deployed to production
- [ ] Access configured (manual)
- [ ] Data synced to D1

---

## Step 1: Code Changes

### 1.1 Update `dashboard/worker/wrangler.jsonc`

**Before:**
```jsonc
{
  "name": "agents-md-dashboard",
  "main": "src/index.ts",
  "compatibility_date": "2025-01-19",
  "d1_databases": [
    {
      "binding": "DB",
      "database_name": "agents-md-dashboard",
      "database_id": "placeholder-replace-after-d1-create"
    }
  ],
  "dev": {
    "port": 8787
  }
}
```

**After:**
```jsonc
{
  "name": "agents-md-dashboard",
  "main": "src/index.ts",
  "compatibility_date": "2025-01-19",
  "account_id": "2469e0bba8bf4d732f65a093985146d6",
  
  // Disable public workers.dev URL - only allow access via custom domain (protected by Access)
  "workers_dev": false,
  
  // Custom domain
  "routes": [
    { "pattern": "agents-md.devtools.cfdata.org", "custom_domain": true }
  ],
  
  // Serve frontend static assets from worker
  "assets": {
    "directory": "../frontend/dist"
  },
  
  "d1_databases": [
    {
      "binding": "DB",
      "database_name": "agents-md-dashboard",
      "database_id": "<REPLACE_AFTER_D1_CREATE>"
    }
  ],
  
  "dev": {
    "port": 8787
  }
}
```

### 1.2 Update `dashboard/worker/src/types.ts`

Add ASSETS binding type:

```typescript
export interface Env {
  DB: D1Database;
  ASSETS: Fetcher;  // Add this line
}
```

### 1.3 Update `dashboard/worker/src/index.ts`

Modify the router to serve static assets for non-API routes:

```typescript
async function handleRequest(
  request: Request,
  env: Env
): Promise<Response> {
  const url = new URL(request.url);
  const path = url.pathname;

  // Handle CORS preflight
  if (request.method === "OPTIONS") {
    return new Response(null, {
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET, OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type",
      },
    });
  }

  // API routes - only allow GET
  if (path.startsWith("/api/")) {
    if (request.method !== "GET") {
      return errorResponse("method_not_allowed", "Only GET requests are allowed", 405);
    }

    // ... existing API route handlers ...
  }

  // Health check
  if (path === "/health") {
    return jsonResponse({ status: "ok" });
  }

  // For non-API routes, serve static assets (SPA)
  // This handles the React frontend
  return env.ASSETS.fetch(request);
}
```

---

## Step 2: Create D1 Database

```bash
cd dashboard/worker
npx wrangler d1 create agents-md-dashboard
```

Expected output:
```
✅ Successfully created DB 'agents-md-dashboard' in region WNAM
Created your new D1 database.

[[d1_databases]]
binding = "DB"
database_name = "agents-md-dashboard"
database_id = "<UUID>"
```

**Action:** Copy the `database_id` UUID and update `wrangler.jsonc`.

---

## Step 3: Apply Schema

```bash
cd dashboard/worker
npx wrangler d1 execute agents-md-dashboard --file=schema.sql
```

---

## Step 4: Build Frontend (if needed)

```bash
cd dashboard/frontend
bun run build
```

---

## Step 5: Deploy Worker

```bash
cd dashboard/worker
npx wrangler deploy
```

Expected output:
```
Total Upload: X KiB / gzip: X KiB
Worker Startup Time: X ms
Your Worker has access to the following bindings:
- D1 Database: DB (agents-md-dashboard)
- Assets: assets
Uploaded agents-md-dashboard
Published agents-md-dashboard
  https://agents-md-dashboard.<subdomain>.workers.dev
  agents-md.devtools.cfdata.org (custom domain)
```

---

## Step 6: Configure Cloudflare Access (Manual)

1. Go to Cloudflare Dashboard > Zero Trust > Access > Applications
2. Create new application for `agents-md.devtools.cfdata.org`
3. Configure authentication policy (e.g., Cloudflare email domain)
4. Test access

---

## Step 7: Sync Data to D1

After deployment, sync repository data:

```bash
# From project root
bun run cli discover --group cloudflare/devtools  # If not already done

# Export to SQL
bun run cli sync

# Import to D1
cd dashboard/worker
npx wrangler d1 execute agents-md-dashboard --file=../../data/sync.sql
```

---

## Verification

1. Visit `https://agents-md.devtools.cfdata.org` - should show dashboard
2. Check `https://agents-md.devtools.cfdata.org/api/stats` - should return JSON
3. Verify Access authentication is required

---

## Rollback

If issues occur:

```bash
# Delete the worker
cd dashboard/worker
npx wrangler delete agents-md-dashboard

# Delete the D1 database (destructive!)
npx wrangler d1 delete agents-md-dashboard
```

---

## Files Modified

| File | Change |
|------|--------|
| `dashboard/worker/wrangler.jsonc` | Add account_id, routes, assets, update database_id |
| `dashboard/worker/src/types.ts` | Add ASSETS binding type |
| `dashboard/worker/src/index.ts` | Add static asset serving for SPA |
