package bind

import (
	"database/sql/driver"
	"embed"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xstring"
)

//go:embed patches/*/*.sql
var patchesFS embed.FS

var _ Binder = (*Patches)(nil)

type patchesType string

const (
	PATCH_SELECT = patchesType("select")
	PATCH_UPDATE = patchesType("update")
	PATCH_INSERT = patchesType("insert")
	PATCH_DELETE = patchesType("delete")
)

type patchEntry struct {
	name  string
	sql   string
	Count atomic.Int64
}

type Patches struct {
	patches map[string]*patchEntry
}

func NewPatches(t patchesType) *Patches {
	patches := make(map[string]*patchEntry)
	entries, err := patchesFS.ReadDir("patches/" + string(t))
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
		inBytes, err := patchesFS.ReadFile("patches/" + string(t) + "/" + name)
		if err != nil {
			panic(fmt.Sprintf("patches: read %s: %v", name, err))
		}
		outBytes, err := patchesFS.ReadFile("patches/" + string(t) + "/" + outName)
		if err != nil {
			panic(fmt.Sprintf("patches: read %s: %v", outName, err))
		}
		key := minify(xstring.FromBytes(inBytes))
		patches[key] = &patchEntry{name: string(t) + "/" + num, sql: xstring.FromBytes(outBytes)}
	}
	return &Patches{
		patches: patches,
	}
}

func minify(s string) string {
	ss := strings.Split(s, "\n")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return strings.TrimSpace(strings.Join(ss, " "))
}

func (p *Patches) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	if e, ok := p.patches[minify(sql)]; ok {
		if e.Count.Add(1) == 1 {
			fmt.Printf("[YDB] patch %q excluded from %v\n", e.name, p.UnusedPatches())
		}
		return e.sql, args, nil
	}
	return sql, args, nil
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
