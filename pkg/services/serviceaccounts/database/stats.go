package database

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/serviceaccounts"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

func (s *ServiceAccountsStoreImpl) GetUsageMetrics(ctx context.Context) (*serviceaccounts.Stats, error) {
	dialect := s.sqlStore.GetDialect()

	// YDB does not support multiple scalar subqueries in one SELECT (error "missing '::' at 'COUNT'").
	// Run three separate queries and combine the results.
	if dialect.DriverName() == migrator.YDB {
		return s.getUsageMetricsYDB(ctx, dialect)
	}

	sb := &db.SQLBuilder{}
	sb.Write("SELECT ")
	sb.Write(`(SELECT COUNT(*) FROM ` + dialect.Quote("user") + ` ` +
		`WHERE is_service_account = ` + dialect.BooleanStr(true) + `) AS serviceaccounts,`)
	sb.Write(`(SELECT COUNT(*) FROM ` + dialect.Quote("api_key") + ` ` +
		`WHERE service_account_id IS NOT NULL) AS serviceaccount_tokens,`)
	sb.Write(`(SELECT COUNT(*) FROM ` + dialect.Quote("org_user") + ` AS ou ` +
		`JOIN ` + dialect.Quote("user") + ` AS u ON u.id = ou.user_id ` +
		`WHERE u.is_disabled = ` + dialect.BooleanStr(false) + ` ` +
		`AND u.is_service_account = ` + dialect.BooleanStr(true) + ` ` +
		`AND ou.role=?) AS serviceaccounts_with_no_role`)
	sb.AddParams("None")

	var sqlStats serviceaccounts.Stats
	if err := s.sqlStore.WithDbSession(ctx, func(sess *db.Session) error {
		_, err := sess.SQL(sb.GetSQLString(), sb.GetParams()...).Get(&sqlStats)
		return err
	}); err != nil {
		return nil, err
	}

	sqlStats.ForcedExpiryEnabled = s.cfg.SATokenExpirationDayLimit != 0

	return &sqlStats, nil
}

func (s *ServiceAccountsStoreImpl) getUsageMetricsYDB(ctx context.Context, dialect migrator.Dialect) (*serviceaccounts.Stats, error) {
	var nServiceaccounts, nTokens, nWithNoRole int64

	err := s.sqlStore.WithDbSession(ctx, func(sess *db.Session) error {
		// 1) count of service accounts
		q1 := `SELECT COUNT(*) AS c FROM ` + dialect.Quote("user") + ` WHERE is_service_account = ` + dialect.BooleanStr(true)
		if _, err := sess.SQL(q1).Get(&nServiceaccounts); err != nil {
			return err
		}

		// 2) count of tokens (api_key with service_account_id IS NOT NULL)
		q2 := `SELECT COUNT(*) AS c FROM ` + dialect.Quote("api_key") + ` WHERE service_account_id IS NOT NULL`
		if _, err := sess.SQL(q2).Get(&nTokens); err != nil {
			return err
		}

		// 3) service accounts with no role (role = 'None')
		q3 := `SELECT COUNT(*) AS c FROM ` + dialect.Quote("org_user") + ` AS ou ` +
			`JOIN ` + dialect.Quote("user") + ` AS u ON u.id = ou.user_id ` +
			`WHERE u.is_disabled = ` + dialect.BooleanStr(false) + ` AND u.is_service_account = ` + dialect.BooleanStr(true) + ` AND ou.role = ?`
		if _, err := sess.SQL(q3, "None").Get(&nWithNoRole); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &serviceaccounts.Stats{
		ServiceAccounts:           nServiceaccounts,
		Tokens:                    nTokens,
		ServiceAccountsWithNoRole: nWithNoRole,
		ForcedExpiryEnabled:       s.cfg.SATokenExpirationDayLimit != 0,
	}, nil
}
