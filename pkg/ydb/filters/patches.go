package filters

import (
	"embed"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

//go:embed patches/*.sql
var patchesFS embed.FS

var (
	_ core.Filter = (*Patches)(nil)
)

type patchEntry struct {
	name  string
	sql   string
	Count atomic.Int64
}

type Patches struct {
	m map[string]*patchEntry
}

func newPatches() *Patches {
	m := make(map[string]*patchEntry)
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
		key := strings.TrimSpace(string(inBytes))
		m[key] = &patchEntry{name: num, sql: string(outBytes)}
	}
	return &Patches{m: m}
}

func (p *Patches) Do(sql string, _ core.Dialect, _ *core.Table) string {
	if e, ok := p.m[strings.TrimSpace(sql)]; ok {
		e.Count.Add(1)
		fmt.Println(p.UnusedPatches())
		return e.sql
	}
	return sql
}

func (p *Patches) UnusedPatches() []string {
	var unused []string
	for _, e := range p.m {
		if e.Count.Load() == 0 {
			unused = append(unused, e.name)
		}
	}
	return unused
}
