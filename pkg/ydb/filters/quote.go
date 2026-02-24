package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var (
	_ core.Filter = (*Quote)(nil)
)

type Quote struct{}

func (s *Quote) Do(sql string, dialect core.Dialect, table *core.Table) string {
	prefix, suffix := byte('`'), byte('`')
	raw := []byte(sql)
	for i, cnt := 0, 0; i < len(raw); i = i + 1 {
		if raw[i] == '`' {
			if cnt%2 == 0 {
				raw[i] = prefix
			} else {
				raw[i] = suffix
			}
			cnt++
		}
	}
	return string(raw)
}
