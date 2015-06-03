package builder_test

import (
	"fmt"
	"github.com/influxdb/influxdb/influxql"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {

	s := `SELECT sum(value) FROM "kbps" WHERE time > now() - 120s AND deliveryservice='steam-dns' and cachegroup = 'total' GROUP BY time(60s)`
	stmt := &influxql.SelectStatement{
		IsRawQuery: false,
		Fields: []*influxql.Field{
			{Expr: &influxql.Call{Name: "sum", Args: []influxql.Expr{&influxql.VarRef{Val: "value"}}}},
		},
		Sources:    []influxql.Source{&influxql.Measurement{Name: "kbps"}},
		Dimensions: []*influxql.Dimension{{Expr: &influxql.Call{Name: "time", Args: []influxql.Expr{&influxql.DurationLiteral{Val: 60 * time.Second}}}}},
		Condition: &influxql.BinaryExpr{ // 1
			Op: influxql.AND,
			LHS: &influxql.BinaryExpr{ // 2
				Op: influxql.AND,
				LHS: &influxql.BinaryExpr{ //3
					Op:  influxql.GT,
					LHS: &influxql.VarRef{Val: "time"},
					RHS: &influxql.BinaryExpr{
						Op:  influxql.SUB,
						LHS: &influxql.Call{Name: "now"},
						RHS: &influxql.DurationLiteral{Val: mustParseDuration("120s")},
					},
				},
				RHS: &influxql.BinaryExpr{
					Op:  influxql.EQ,
					LHS: &influxql.VarRef{Val: "deliveryservice"},
					RHS: &influxql.StringLiteral{Val: "steam-dns"},
				},
			},
			RHS: &influxql.BinaryExpr{
				Op:  influxql.EQ,
				LHS: &influxql.VarRef{Val: "cachegroup"},
				RHS: &influxql.StringLiteral{Val: "total"},
			},
		},
	}

	fmt.Println(s)
	fmt.Println(stmt.String())
}

func mustParseDuration(s string) time.Duration {
	d, err := influxql.ParseDuration(s)
	panicIfErr(err)
	return d
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
