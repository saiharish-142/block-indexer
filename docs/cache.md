Redis cache strategy
====================

- Keys:
  - `latest:blocks`: list/hash of recent blocks (TTL 60s) for quick homepage fetches.
  - `address:{addr}:txs`: sorted set scored by block number for ordered recent transactions.
  - `tx:{hash}`: string/json cache for hot tx lookups (TTL 5m).
- Patterns:
  - Write-through for address recent txs during indexing using `cache.CacheRecentTx`.
  - Read-through for API handlers; if a miss occurs, hydrate from Postgres and set TTL.
- TTL guidance:
  - Heads / recent blocks: 30–60s.
  - Address tx list: 10–30m depending on churn.
  - Tx details: 5–15m.
- Invalidation:
  - Replace on write for heads/tx; for reorg handling, delete impacted keys for reorged ranges.
