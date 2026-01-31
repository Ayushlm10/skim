# agents-md-dsitribution

To install dependencies:

```bash
bun install
```

To run:

```bash
bun run index.ts
```

This project was created using `bun init` in bun v1.3.5. [Bun](https://bun.com) is a fast all-in-one JavaScript runtime.

# 1. Delete the local SQLite database entirely
rm -f data/agents.db
# 2. Delete the local D1 data
rm -rf dashboard/worker/.wrangler
# 3. Run fresh discovery
bun run cli discover --group cloudflare/devtools

# 2. Generate sync SQL from local SQLite
bun run cli sync

# 3. Create D1 schema (from worker directory)
cd dashboard/worker
pnpm dlx wrangler d1 execute agents-md-dashboard --local --file=schema.sql
# 4. Import the data
pnpm dlx wrangler d1 execute agents-md-dashboard --local --file=../../data/sync.sql
# 5. Start the worker API (keep this running)
bun run dev
Then in a second terminal:
# 6. Start the frontend
cd dashboard/frontend
bun run dev
