## Appendix: Database Internals

*Deep dive into how databases actually work. Understanding storage engines, query optimization, replication, and sharding transforms database usage from cargo-culting to informed decisions.*

### J.1 Storage Engine Fundamentals

**Goal:** Understand how databases store and retrieve data

**Learning objectives:**
- Understand B-tree vs LSM-tree trade-offs
- Know how storage engines impact performance
- Choose appropriate storage engines for workloads

**Tasks:**
- [ ] Create `experiments/scenarios/storage-engines/`
- [ ] Storage engine role:
  - [ ] Data organization on disk
  - [ ] Index structures
  - [ ] Read/write paths
  - [ ] Recovery mechanisms
- [ ] B-tree storage:
  - [ ] Tree structure and balancing
  - [ ] Page-based organization
  - [ ] In-place updates
  - [ ] Read-optimized characteristics
- [ ] B-tree operations:
  - [ ] Point lookups
  - [ ] Range scans
  - [ ] Insertions and splits
  - [ ] Deletions and merges
- [ ] LSM-tree storage:
  - [ ] Log-structured design
  - [ ] Memtable + SSTables
  - [ ] Compaction strategies
  - [ ] Write-optimized characteristics
- [ ] LSM-tree operations:
  - [ ] Write path (memtable → SSTable)
  - [ ] Read path (merge across levels)
  - [ ] Compaction (leveled, tiered, FIFO)
  - [ ] Bloom filters for efficiency
- [ ] Trade-offs:
  - [ ] Write amplification
  - [ ] Read amplification
  - [ ] Space amplification
  - [ ] Workload fit
- [ ] Storage engines in practice:
  - [ ] PostgreSQL (B-tree heap)
  - [ ] MySQL InnoDB (B-tree)
  - [ ] RocksDB (LSM-tree)
  - [ ] LevelDB (LSM-tree)
  - [ ] Cassandra (LSM-tree)
- [ ] Column stores:
  - [ ] Column-oriented storage
  - [ ] Compression benefits
  - [ ] Analytics workloads
  - [ ] ClickHouse, DuckDB
- [ ] **ADR:** Document storage engine selection criteria

---

### J.2 Indexing Deep Dive

**Goal:** Master database indexing for performance

**Learning objectives:**
- Understand index types and structures
- Design effective indexing strategies
- Troubleshoot index-related performance issues

**Tasks:**
- [ ] Create `experiments/scenarios/database-indexing/`
- [ ] Index fundamentals:
  - [ ] Index as lookup structure
  - [ ] Primary vs secondary indexes
  - [ ] Clustered vs non-clustered
  - [ ] Index overhead (writes, storage)
- [ ] B-tree indexes:
  - [ ] Structure and navigation
  - [ ] Prefix compression
  - [ ] Index-only scans
  - [ ] Range query efficiency
- [ ] Hash indexes:
  - [ ] O(1) point lookups
  - [ ] No range query support
  - [ ] Memory-resident typically
  - [ ] Use cases
- [ ] Composite indexes:
  - [ ] Multi-column indexes
  - [ ] Column ordering importance
  - [ ] Leftmost prefix rule
  - [ ] Covering indexes
- [ ] Partial indexes:
  - [ ] Conditional indexing
  - [ ] Reduced storage
  - [ ] Query matching requirements
  - [ ] PostgreSQL partial indexes
- [ ] Expression indexes:
  - [ ] Indexing computed values
  - [ ] Function-based indexes
  - [ ] JSON path indexes
- [ ] Full-text indexes:
  - [ ] Text search capabilities
  - [ ] Inverted indexes
  - [ ] Stemming and tokenization
  - [ ] PostgreSQL tsvector
- [ ] Spatial indexes:
  - [ ] R-tree structure
  - [ ] Geographic queries
  - [ ] PostGIS
- [ ] Index maintenance:
  - [ ] Index bloat
  - [ ] Reindexing strategies
  - [ ] Statistics updates
  - [ ] Index monitoring
- [ ] Anti-patterns:
  - [ ] Over-indexing
  - [ ] Unused indexes
  - [ ] Wrong column order
  - [ ] Indexing low-cardinality columns
- [ ] Design indexing strategy for sample schema
- [ ] **ADR:** Document indexing guidelines

---

### J.3 Query Planning & Optimization

**Goal:** Understand how databases execute queries

**Learning objectives:**
- Read and interpret query plans
- Optimize slow queries
- Understand query optimizer decisions

**Tasks:**
- [ ] Create `experiments/scenarios/query-optimization/`
- [ ] Query processing pipeline:
  - [ ] Parsing
  - [ ] Planning/optimization
  - [ ] Execution
- [ ] Query planner role:
  - [ ] Generate execution plans
  - [ ] Estimate costs
  - [ ] Choose optimal plan
- [ ] Plan operations:
  - [ ] Sequential scan
  - [ ] Index scan
  - [ ] Bitmap scan
  - [ ] Nested loop join
  - [ ] Hash join
  - [ ] Merge join
- [ ] EXPLAIN analysis:
  - [ ] Reading EXPLAIN output
  - [ ] EXPLAIN ANALYZE for actual times
  - [ ] Cost estimates vs actuals
  - [ ] Row estimate accuracy
- [ ] Statistics:
  - [ ] Table statistics
  - [ ] Column statistics (histograms)
  - [ ] Statistics freshness
  - [ ] ANALYZE command
- [ ] Join optimization:
  - [ ] Join order selection
  - [ ] Join algorithm selection
  - [ ] Join hints (when appropriate)
- [ ] Subquery optimization:
  - [ ] Correlated vs uncorrelated
  - [ ] Subquery flattening
  - [ ] Exists vs IN vs JOIN
- [ ] Common issues:
  - [ ] Missing indexes
  - [ ] Stale statistics
  - [ ] Type mismatches
  - [ ] Function calls preventing index use
- [ ] Query rewrites:
  - [ ] Equivalent transformations
  - [ ] Performance improvements
  - [ ] Maintaining correctness
- [ ] PostgreSQL-specific:
  - [ ] pg_stat_statements
  - [ ] auto_explain
  - [ ] Plan caching
- [ ] Optimize real-world queries
- [ ] **ADR:** Document query optimization process

---

### J.4 Transactions & Isolation

**Goal:** Understand transaction guarantees and isolation levels

**Learning objectives:**
- Understand ACID properties deeply
- Choose appropriate isolation levels
- Debug transaction-related issues

**Tasks:**
- [ ] Create `experiments/scenarios/transactions/`
- [ ] ACID properties:
  - [ ] Atomicity (all or nothing)
  - [ ] Consistency (valid state transitions)
  - [ ] Isolation (concurrent transaction behavior)
  - [ ] Durability (committed = persistent)
- [ ] Isolation anomalies:
  - [ ] Dirty reads
  - [ ] Non-repeatable reads
  - [ ] Phantom reads
  - [ ] Write skew
  - [ ] Lost updates
- [ ] Isolation levels:
  - [ ] Read Uncommitted
  - [ ] Read Committed
  - [ ] Repeatable Read
  - [ ] Serializable
- [ ] Implementation approaches:
  - [ ] Locking (2PL)
  - [ ] MVCC (Multi-Version Concurrency Control)
  - [ ] Optimistic concurrency
- [ ] MVCC deep dive:
  - [ ] Version chains
  - [ ] Snapshot isolation
  - [ ] Garbage collection
  - [ ] Write conflicts
- [ ] PostgreSQL specifics:
  - [ ] Default Read Committed
  - [ ] Serializable Snapshot Isolation (SSI)
  - [ ] Transaction IDs
  - [ ] VACUUM necessity
- [ ] MySQL/InnoDB specifics:
  - [ ] Gap locking
  - [ ] Next-key locking
  - [ ] Deadlock detection
- [ ] Deadlocks:
  - [ ] Causes
  - [ ] Detection
  - [ ] Prevention strategies
  - [ ] Application handling
- [ ] Long transactions:
  - [ ] Problems caused
  - [ ] MVCC bloat
  - [ ] Lock contention
  - [ ] Best practices
- [ ] Demonstrate isolation anomalies
- [ ] **ADR:** Document isolation level selection

---

### J.5 Replication Internals

**Goal:** Understand database replication mechanisms

**Learning objectives:**
- Understand replication protocols
- Configure replication appropriately
- Handle replication lag and failures

**Tasks:**
- [ ] Create `experiments/scenarios/db-replication/`
- [ ] Replication goals:
  - [ ] High availability
  - [ ] Read scaling
  - [ ] Geographic distribution
  - [ ] Disaster recovery
- [ ] Physical replication:
  - [ ] Byte-level log shipping
  - [ ] WAL (Write-Ahead Log) streaming
  - [ ] Block-level replication
  - [ ] Standby types (hot, warm)
- [ ] Logical replication:
  - [ ] Row-level changes
  - [ ] Schema flexibility
  - [ ] Selective replication
  - [ ] Cross-version replication
- [ ] PostgreSQL streaming replication:
  - [ ] WAL senders and receivers
  - [ ] Synchronous vs asynchronous
  - [ ] Replication slots
  - [ ] Cascading replication
- [ ] PostgreSQL logical replication:
  - [ ] Publications and subscriptions
  - [ ] Logical decoding
  - [ ] Conflict handling
- [ ] MySQL replication:
  - [ ] Binary log replication
  - [ ] GTID (Global Transaction ID)
  - [ ] Group Replication
  - [ ] Semi-synchronous replication
- [ ] Replication lag:
  - [ ] Causes
  - [ ] Monitoring
  - [ ] Impact on reads
  - [ ] Mitigation strategies
- [ ] Failover:
  - [ ] Automatic vs manual
  - [ ] Promotion process
  - [ ] Client reconnection
  - [ ] Split-brain prevention
- [ ] Replication tools:
  - [ ] Patroni (PostgreSQL HA)
  - [ ] Orchestrator (MySQL)
  - [ ] pg_basebackup
- [ ] Set up replication cluster
- [ ] **ADR:** Document replication architecture

---

### J.6 Sharding & Partitioning

**Goal:** Scale databases beyond single node

**Learning objectives:**
- Understand partitioning strategies
- Implement database sharding
- Handle cross-shard operations

**Tasks:**
- [ ] Create `experiments/scenarios/db-sharding/`
- [ ] Partitioning vs Sharding:
  - [ ] Partitioning (single database)
  - [ ] Sharding (multiple databases)
  - [ ] When each applies
- [ ] Table partitioning:
  - [ ] Range partitioning
  - [ ] List partitioning
  - [ ] Hash partitioning
  - [ ] Composite partitioning
- [ ] PostgreSQL partitioning:
  - [ ] Declarative partitioning
  - [ ] Partition pruning
  - [ ] Partition maintenance
  - [ ] Partition-wise joins
- [ ] Sharding strategies:
  - [ ] Key-based (hash) sharding
  - [ ] Range-based sharding
  - [ ] Directory-based sharding
  - [ ] Geographic sharding
- [ ] Shard key selection:
  - [ ] High cardinality
  - [ ] Query patterns alignment
  - [ ] Even distribution
  - [ ] Access locality
- [ ] Cross-shard challenges:
  - [ ] Distributed queries
  - [ ] Distributed transactions
  - [ ] Referential integrity
  - [ ] Global sequences
- [ ] Sharding solutions:
  - [ ] Vitess (MySQL)
  - [ ] Citus (PostgreSQL)
  - [ ] Application-level sharding
  - [ ] ProxySQL routing
- [ ] Resharding:
  - [ ] Adding shards
  - [ ] Rebalancing data
  - [ ] Online resharding
  - [ ] Minimizing downtime
- [ ] NewSQL alternatives:
  - [ ] CockroachDB
  - [ ] TiDB
  - [ ] YugabyteDB
  - [ ] Automatic sharding
- [ ] Implement sharded database
- [ ] **ADR:** Document sharding strategy

---

### J.7 Connection Management

**Goal:** Optimize database connection handling

**Learning objectives:**
- Understand connection overhead
- Implement connection pooling
- Tune connection settings

**Tasks:**
- [ ] Create `experiments/scenarios/connection-management/`
- [ ] Connection overhead:
  - [ ] TCP handshake
  - [ ] TLS negotiation
  - [ ] Authentication
  - [ ] Session initialization
  - [ ] Memory per connection
- [ ] Connection pooling:
  - [ ] Pool concepts
  - [ ] Pool sizing
  - [ ] Connection lifecycle
  - [ ] Pool exhaustion handling
- [ ] Application-level pooling:
  - [ ] HikariCP (Java)
  - [ ] pgx pool (Go)
  - [ ] SQLAlchemy pool (Python)
  - [ ] Configuration best practices
- [ ] External poolers:
  - [ ] PgBouncer
  - [ ] Pgpool-II
  - [ ] ProxySQL (MySQL)
- [ ] PgBouncer deep dive:
  - [ ] Pooling modes (session, transaction, statement)
  - [ ] Configuration
  - [ ] Limitations (prepared statements)
  - [ ] Monitoring
- [ ] Pool sizing:
  - [ ] Too small: queuing
  - [ ] Too large: resource exhaustion
  - [ ] Optimal sizing guidelines
  - [ ] Connections vs CPU cores
- [ ] Connection limits:
  - [ ] Database max_connections
  - [ ] Per-user limits
  - [ ] Per-database limits
  - [ ] Monitoring connection usage
- [ ] Kubernetes considerations:
  - [ ] Pod scaling impact
  - [ ] Sidecar poolers
  - [ ] Connection storms on restart
  - [ ] Graceful shutdown
- [ ] Troubleshooting:
  - [ ] Connection leaks
  - [ ] Pool exhaustion
  - [ ] Idle connection timeout
  - [ ] Connection validation
- [ ] Configure PgBouncer for Kubernetes
- [ ] **ADR:** Document connection pooling strategy

---

### J.8 Backup & Recovery

**Goal:** Implement reliable database backup strategies

**Learning objectives:**
- Understand backup types and trade-offs
- Implement point-in-time recovery
- Test recovery procedures

**Tasks:**
- [ ] Create `experiments/scenarios/db-backup-recovery/`
- [ ] Backup types:
  - [ ] Logical backups (pg_dump, mysqldump)
  - [ ] Physical backups (file-level)
  - [ ] Continuous archiving (WAL)
  - [ ] Snapshots (storage-level)
- [ ] Logical backup:
  - [ ] SQL dump format
  - [ ] Custom format (parallel restore)
  - [ ] Selective backup (tables, schemas)
  - [ ] Restore process
- [ ] Physical backup:
  - [ ] pg_basebackup
  - [ ] File system backup
  - [ ] Consistent snapshots
  - [ ] Faster for large databases
- [ ] Point-in-Time Recovery (PITR):
  - [ ] Continuous WAL archiving
  - [ ] Base backup + WAL replay
  - [ ] Recovery target specification
  - [ ] Recovery timeline
- [ ] WAL archiving:
  - [ ] archive_command
  - [ ] WAL-G, pgBackRest
  - [ ] S3/GCS storage
  - [ ] Retention policies
- [ ] Backup tools:
  - [ ] pgBackRest (PostgreSQL)
  - [ ] Barman (PostgreSQL)
  - [ ] WAL-G (PostgreSQL, MySQL)
  - [ ] Percona XtraBackup (MySQL)
- [ ] Kubernetes backup:
  - [ ] Velero integration
  - [ ] Operator-based backup
  - [ ] Storage snapshots
  - [ ] Cross-cluster backup
- [ ] Recovery testing:
  - [ ] Regular restore tests
  - [ ] Recovery time measurement
  - [ ] Data validation
  - [ ] Runbook maintenance
- [ ] RPO and RTO:
  - [ ] Recovery Point Objective
  - [ ] Recovery Time Objective
  - [ ] Backup frequency alignment
  - [ ] Architecture implications
- [ ] Implement PITR with pgBackRest
- [ ] **ADR:** Document backup strategy

---

### J.9 Database Observability

**Goal:** Monitor and troubleshoot database performance

**Learning objectives:**
- Instrument database monitoring
- Identify performance issues
- Build effective dashboards

**Tasks:**
- [ ] Create `experiments/scenarios/db-observability/`
- [ ] Key metrics:
  - [ ] Query throughput (QPS)
  - [ ] Query latency (p50, p99)
  - [ ] Connection count
  - [ ] Cache hit ratios
- [ ] PostgreSQL statistics:
  - [ ] pg_stat_user_tables
  - [ ] pg_stat_user_indexes
  - [ ] pg_stat_activity
  - [ ] pg_stat_statements
- [ ] pg_stat_statements:
  - [ ] Query fingerprinting
  - [ ] Execution statistics
  - [ ] Top queries by time
  - [ ] Query plan changes
- [ ] Wait events:
  - [ ] pg_stat_activity.wait_event
  - [ ] Lock waits
  - [ ] I/O waits
  - [ ] CPU waits
- [ ] Lock monitoring:
  - [ ] pg_locks view
  - [ ] Lock conflicts
  - [ ] Deadlock detection
  - [ ] Blocking queries
- [ ] Replication monitoring:
  - [ ] pg_stat_replication
  - [ ] Replication lag
  - [ ] Slot status
  - [ ] WAL generation rate
- [ ] Storage monitoring:
  - [ ] Table and index sizes
  - [ ] Bloat estimation
  - [ ] Disk usage trends
  - [ ] VACUUM monitoring
- [ ] Prometheus exporters:
  - [ ] postgres_exporter
  - [ ] mysqld_exporter
  - [ ] Custom queries
- [ ] Grafana dashboards:
  - [ ] Overview dashboard
  - [ ] Query performance dashboard
  - [ ] Replication dashboard
  - [ ] Alert configuration
- [ ] Log analysis:
  - [ ] Slow query log
  - [ ] Error log patterns
  - [ ] Log aggregation
- [ ] Build comprehensive monitoring
- [ ] **ADR:** Document database monitoring strategy

---
