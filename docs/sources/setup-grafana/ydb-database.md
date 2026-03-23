# YDB as Grafana database (experimental, fork)

This branch (`ydb-v10.1.5`) ports the YDB backend support from `ydb-v12.2.1` (diff vs tag `v12.2.1` / stable 12.2.1) onto Grafana **v10.1.5**.

## Scope

- In-tree XORM fork under `pkg/util/xorm` (with `pkg/util/xorm/core` and YDB dialect), used via imports `github.com/grafana/grafana/pkg/util/xorm` instead of `xorm.io/xorm`.
- `pkg/util/sqlite` helpers required by the updated SQLite dialect.
- SQL store migrator: `migrator.YDB`, `YDBDialect`, `DialectRecursiveCTE`, and migration helpers (e.g. `RawSQLMigration.YDB()`).
- Integration tests: `make test-go-integration-ydb` with `GRAFANA_TEST_DB=ydb` (start YDB separately).

## Dependencies

- `github.com/ydb-platform/ydb-go-sdk/v3` (version pinned in `go.mod`; resolves with `replace google.golang.org/grpc => v1.57.1` for compatibility with this tree).
- `github.com/ydb-platform/ydb-go-yc-metadata` for optional `yc_auth=metadata_credentials` in DSN (see YDB SDK docs).

## Not ported

Paths that exist only in Grafana 12.x (e.g. `pkg/storage/unified`, some authz SQL files) are omitted; their YDB-specific edits do not apply to v10.1.5.

`pkg/infra/serverlock` in 10.x uses a different implementation than 12.x; the YDB `UPSERT … RETURNING id` / `createLock` flow from `ydb-v12.2.1` was not mechanically ported. If you rely on HA server locks against YDB, validate lock behaviour under concurrency.
