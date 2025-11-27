Postgres and Redis deployment notes
===================================

Managed services:
- Prefer managed Postgres (RDS/CloudSQL) with pgBouncer for pooling.
- Managed Redis (Elasticache/Memorystore) with replica for HA.

In-cluster (development or on-prem):
- Deploy Postgres as a StatefulSet with a PersistentVolume; enable WAL archiving and logical replication if read replicas needed.
- Deploy Redis as a StatefulSet with AOF persistence and a small replica count; use PodDisruptionBudgets.

Backups:
- Postgres: daily base backup + WAL shipping; test PITR.
- Redis: snapshot every 15–30m; back up to blob storage.
