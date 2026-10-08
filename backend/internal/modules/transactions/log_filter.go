package transactions

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// LogFilter is one row of the logs query builder. Filters are ANDed.
type LogFilter struct {
	Target string `json:"target"` // message | level | field
	Key    string `json:"key"`    // dot path into fields, e.g. "request.lpn"; empty = whole fields JSON
	Op     string `json:"op"`     // contains | not_contains | eq | neq | regex | not_regex
	Value  string `json:"value"`
}

var logFilterOps = map[string]string{
	"contains":     "strpos(lower(%s), lower(%s)) > 0",
	"not_contains": "strpos(lower(%s), lower(%s)) = 0",
	"eq":           "%s = %s",
	"neq":          "%s <> %s",
	"regex":        "%s ~* %s",
	"not_regex":    "%s !~* %s",
}

// buildLogWhere returns a WHERE clause and its args. Args start at $1;
// callers append further args after len(args).
func buildLogWhere(filters []LogFilter, args []any) (string, []any, error) {
	conds := []string{"created_at >= $1", "created_at <= $2"}
	for _, f := range filters {
		if f.Value == "" {
			continue
		}
		tmpl, ok := logFilterOps[f.Op]
		if !ok {
			return "", nil, fmt.Errorf("invalid op %q", f.Op)
		}
		var col string
		switch f.Target {
		case "message":
			col = "coalesce(message, '')"
		case "level":
			col = "coalesce(level, '')"
		case "field":
			if f.Key == "" {
				col = "coalesce(fields::text, '')"
			} else {
				args = append(args, pq.Array(strings.Split(f.Key, ".")))
				col = fmt.Sprintf("coalesce(fields #>> $%d::text[], '')", len(args))
			}
		default:
			return "", nil, fmt.Errorf("invalid target %q", f.Target)
		}
		args = append(args, f.Value)
		conds = append(conds, fmt.Sprintf(tmpl, col, fmt.Sprintf("$%d", len(args))))
	}
	return strings.Join(conds, " and "), args, nil
}
