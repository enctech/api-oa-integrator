package transactions

import "testing"

func TestBuildLogWhere(t *testing.T) {
	where, args, err := buildLogWhere([]LogFilter{
		{Target: "message", Op: "contains", Value: "exit"},
		{Target: "field", Key: "req.lpn", Op: "eq", Value: "ABC"},
		{Target: "level", Op: "eq", Value: ""}, // empty value skipped
	}, []any{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	want := "created_at >= $1 and created_at <= $2 and strpos(lower(coalesce(message, '')), lower($3)) > 0 and coalesce(fields #>> $4::text[], '') = $5"
	if where != want || len(args) != 5 {
		t.Fatalf("got %q (%d args)", where, len(args))
	}
	if _, _, err := buildLogWhere([]LogFilter{{Target: "message", Op: "; drop", Value: "x"}}, nil); err == nil {
		t.Fatal("expected invalid op error")
	}
}
