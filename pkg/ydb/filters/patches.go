package filters

import (
	"embed"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xstring"
)

//go:embed patches/*.sql
var patchesFS embed.FS

var (
	_ core.Filter         = (*Patches)(nil)
	_ core.FilterWithArgs = (*Patches)(nil)
)

type patchEntry struct {
	name  string
	sql   string
	Count atomic.Int64
}

type Patches struct {
	patches map[string]*patchEntry

	fallbacks []core.Filter
}

func newPatches(fallbacks ...core.Filter) *Patches {
	patches := make(map[string]*patchEntry)
	entries, err := patchesFS.ReadDir("patches")
	if err != nil {
		panic(fmt.Sprintf("patches: read dir: %v", err))
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".in.sql") {
			continue
		}
		num := strings.TrimSuffix(name, ".in.sql")
		outName := num + ".out.sql"
		inBytes, err := patchesFS.ReadFile("patches/" + name)
		if err != nil {
			panic(fmt.Sprintf("patches: read %s: %v", name, err))
		}
		outBytes, err := patchesFS.ReadFile("patches/" + outName)
		if err != nil {
			panic(fmt.Sprintf("patches: read %s: %v", outName, err))
		}
		key := minify(xstring.FromBytes(inBytes))
		patches[key] = &patchEntry{name: num, sql: xstring.FromBytes(outBytes)}
	}
	return &Patches{
		patches:   patches,
		fallbacks: fallbacks,
	}
}

func minify(s string) string {
	ss := strings.Split(s, "\n")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return strings.TrimSpace(strings.Join(ss, " "))
}

func (p *Patches) patch(sql string) (string, bool) {
	if e, ok := p.patches[minify(sql)]; ok {
		if e.Count.Add(1) == 1 {
			fmt.Printf("[YDB] patch %q excluded from %v\n", e.name, p.UnusedPatches())
		}

		return e.sql, true
	}

	return sql, false
}

func (p *Patches) Do(sql string, d core.Dialect, t *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (p *Patches) DoWithArgs(sql string, dialect core.Dialect, table *core.Table, args ...any) (string, []any) {
	sql, ok := p.patch(sql)

	for _, f := range p.fallbacks {
		if ff, has := f.(core.FilterWithArgs); has {
			sql, args = ff.DoWithArgs(sql, dialect, table, args...)
		} else if !ok {
			sql = f.Do(sql, dialect, table)
		}
	}

	return sql, args
}

func (p *Patches) UnusedPatches() []string {
	var unused []string
	for _, e := range p.patches {
		if e.Count.Load() == 0 {
			unused = append(unused, e.name)
		}
	}
	return unused
}
