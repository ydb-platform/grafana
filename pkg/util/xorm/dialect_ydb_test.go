package xorm

import (
	"embed"
	"strconv"
	"strings"
	"testing"
)

//go:embed testdata/in_clause_rewrite/*.sql
var inClauseRewriteFS embed.FS

func TestRewriteQueryInClauses_FromFiles(t *testing.T) {
	input, err := inClauseRewriteFS.ReadFile("testdata/in_clause_rewrite/input.sql")
	if err != nil {
		t.Fatalf("read input.sql: %v", err)
	}
	want, err := inClauseRewriteFS.ReadFile("testdata/in_clause_rewrite/output.sql")
	if err != nil {
		t.Fatalf("read output.sql: %v", err)
	}
	query := strings.TrimSpace(string(input))
	expected := strings.TrimSpace(string(want))

	rewritten, m := rewriteQueryInClauses(query)
	if m == nil {
		t.Fatal("YDB rewrite did not apply: inArgMap is nil (expected IN clause to be collapsed)")
	}

	gotNorm := normalizeSQL(rewritten)
	expectedNorm := normalizeSQL(expected)
	if gotNorm != expectedNorm {
		t.Errorf("YDB rewrite result differs from expected output.sql\n--- got (rewritten)\n%s\n--- want (output.sql)\n%s", gotNorm, expectedNorm)
	}
}

// normalizeSQL reduces whitespace for comparison (single spaces, trim).
func normalizeSQL(s string) string {
	var b strings.Builder
	wasSpace := false
	for _, r := range strings.TrimSpace(s) {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !wasSpace {
				b.WriteRune(' ')
				wasSpace = true
			}
			continue
		}
		wasSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// TestRewriteQueryInClauses_FolderQuery checks the real folder list query (IN with 1397 uids)
// is rewritten to "IN $2" and placeholders after the list are renumbered.
func TestRewriteQueryInClauses_FolderQuery(t *testing.T) {
	// Build IN list $2,$3,...,$1398 (same as real GET /api/folders query)
	inList := make([]string, 1397)
	for i := 2; i <= 1398; i++ {
		inList[i-2] = "$" + strconv.Itoa(i)
	}
	inListStr := strings.Join(inList, ",")

	beforeIN := `SELECT dashboard.id, dashboard.org_id, dashboard.uid, dashboard.title, dashboard.slug, dashboard_tag.term, dashboard.is_folder, dashboard.folder_id, dashboard.deleted, folder.uid AS folder_uid, folder.title AS folder_slug, folder.title AS folder_title , dashboard.title AS title  FROM ( SELECT dashboard.id AS id, dashboard.title AS title FROM dashboard  WHERE dashboard.org_id=$1 AND dashboard.uid IN (`
	afterIN := `) AND dashboard.is_folder = true AND dashboard.folder_uid IS NULL  AND dashboard.uid != $1399 AND (dashboard.folder_uid != $1400 OR dashboard.folder_uid IS NULL) AND (((dashboard.org_id = $1401 AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  WHERE f1.org_id = $1402 AND f1.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN ( SELECT ur.role_id FROM user_role AS ur WHERE ur.user_id = $1403 AND (ur.org_id = $1404 OR ur.org_id = $1405) UNION SELECT br.role_id FROM builtin_role AS br WHERE br.role IN ($1406) AND (br.org_id = $1407 OR br.org_id = $1408) ) as all_role ON role.id = all_role.role_id)  AND action IN ($1409, $1410, $1411, $1412)))) AND dashboard.is_folder)) AND dashboard.deleted IS NULL ORDER BY title ASC LIMIT 50 OFFSET 0) AS ids
    INNER JOIN dashboard ON ids.id = dashboard.id
LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
  LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
 ORDER BY dashboard.title ASC`

	query := beforeIN + inListStr + afterIN
	rewritten, m := rewriteQueryInClauses(query)
	if m == nil {
		t.Fatal("YDB rewrite did not apply to folder query (inArgMap is nil)")
	}
	if !strings.Contains(rewritten, "dashboard.uid IN $2") {
		t.Errorf("expected rewritten query to contain 'dashboard.uid IN $2', got fragment: %s", firstFragment(rewritten, "dashboard.uid IN"))
	}
	if strings.Contains(rewritten, "$1399") || strings.Contains(rewritten, "$1400") {
		t.Error("expected placeholders after IN list to be renumbered (no $1399, $1400)")
	}
	// After collapse: $1399 -> $3, $1400 -> $4, etc. (delta -1396)
	if !strings.Contains(rewritten, "dashboard.uid != $3") {
		t.Errorf("expected renumbered $1399 to $3 in rewritten query")
	}
	if n := m.collapsedCount(); n != 1397 {
		t.Errorf("expected collapsedCount 1397, got %d", n)
	}
}

func firstFragment(s, sub string) string {
	if i := strings.Index(s, sub); i >= 0 {
		end := i + len(sub) + 40
		if end > len(s) {
			end = len(s)
		}
		return s[i:end]
	}
	return ""
}

func TestRewriteQueryInClausesNumeric_CollapsesLongIn(t *testing.T) {
	// Simulate folder list query: first " IN (" is a subquery (SELECT...), second is dashboard.uid IN ($2,...,$N)
	query := "SELECT id FROM dashboard WHERE org_id=$1 AND dashboard.uid IN ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60) AND deleted IS NULL"
	rewritten, m := rewriteQueryInClauses(query)
	if m == nil {
		t.Fatal("expected rewrite to apply (inArgMap non-nil)")
	}
	if !strings.Contains(rewritten, "IN $2") {
		t.Errorf("expected rewritten query to contain IN $2, got: %s", rewritten)
	}
	if strings.Contains(rewritten, "IN ($2,$3)") || strings.Contains(rewritten, "$60)") {
		t.Errorf("expected long IN list to be collapsed, got: %s", rewritten)
	}
	// Args: $1 single, $2..$60 -> one list, so 1 + 1 + 0 = 2 segments for 61 placeholders -> 2 args
	if len(m.segments) != 2 {
		t.Errorf("expected 2 segments (one single $1, one list), got %d", len(m.segments))
	}
}

func TestGetLastPartOfColumn(t *testing.T) {
	tests := []struct {
		column   string
		expected string
	}{
		{"table.column", "column"},
		{"column", "column"},
		{".", ""},
		{".x", "x"},
		{"xxxx.xxx.", ""},
	}

	for _, test := range tests {
		if getLastPartOfColumn(test.column) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, getLastPartOfColumn(test.column))
		}
	}
}

func TestRewriteYdbLowerFunctions(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no lower", "SELECT * FROM t WHERE id = $1", "SELECT * FROM t WHERE id = $1"},
		{"lower login", "LOWER(`user`.login)=LOWER($1)", "Unicode::ToLower(`user`.login)=Unicode::ToLower($1)"},
		{"lowercase", "where lower(email)=lower($1)", "where Unicode::ToLower(email)=Unicode::ToLower($1)"},
		// Already YQL: must not add nested Unicode::ToLower inside tolower
		{"idempotent unicode tolower", "Unicode::ToLower(x)=Unicode::ToLower($1)", "Unicode::ToLower(x)=Unicode::ToLower($1)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteYdbLowerFunctions(tt.in)
			if got != tt.want {
				t.Errorf("rewriteYdbLowerFunctions() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRewriteILIKEToLowerLike(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no ilike", "SELECT * FROM t WHERE id = $1", "SELECT * FROM t WHERE id = $1"},
		{"dashboard title ilike", "WHERE dashboard.org_id=$1 AND dashboard.title ILIKE $p2", "WHERE dashboard.org_id=$1 AND LOWER(dashboard.title) LIKE LOWER($p2)"},
		{"table_col ilike", "AND le.name ILIKE $p3", "AND LOWER(le.name) LIKE LOWER($p3)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteILIKEToLowerLike(tt.in)
			if got != tt.want {
				t.Errorf("rewriteILIKEToLowerLike() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldUseCostBasedOptimization(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"simple select", "SELECT * FROM cache_data WHERE key = $1", false},
		{"server_lock", "SELECT * FROM server_lock WHERE operation_uid = $1", false},
		{"dashboard list with permission", "SELECT id FROM dashboard WHERE org_id=$1 AND dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards') ORDER BY title LIMIT 50", true},
		{"folder and permission", "SELECT id FROM dashboard WHERE folder_id IN (SELECT id FROM folder WHERE uid IN (SELECT identifier FROM permission WHERE kind = 'folders')) ORDER BY title", true},
		{"two IN (SELECT", "SELECT a FROM t WHERE x IN (SELECT 1) AND y IN (SELECT 2 FROM u)", true},
		{"one IN (SELECT", "SELECT * FROM permission WHERE role_id IN (SELECT id FROM role)", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUseCostBasedOptimization(tt.query)
			if got != tt.want {
				t.Errorf("shouldUseCostBasedOptimization() = %v, want %v", got, tt.want)
			}
		})
	}
}
