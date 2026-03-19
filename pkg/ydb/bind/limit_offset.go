package bind

import "regexp"

var _ Binder = (*CastLimitOffsetToUint64)(nil)

var (
	limitParamRe  = regexp.MustCompile(`(?i:LIMIT)\s+\?`)
	offsetParamRe = regexp.MustCompile(`(?i:OFFSET)\s+\?`)
)

type CastLimitOffsetToUint64 struct{}

func (f *CastLimitOffsetToUint64) Rebind(sql string, args ...any) (string, []any, error) {
	sql = limitParamRe.ReplaceAllString(sql, "LIMIT CAST(? AS Uint64)")
	sql = offsetParamRe.ReplaceAllString(sql, "OFFSET CAST(? AS Uint64)")
	return sql, args, nil
}
