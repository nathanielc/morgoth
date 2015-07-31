package influxdb

import (
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/influxdb/influxdb/influxql"
	"time"
)

type QueryBuilder struct {
	statement *influxql.SelectStatement
	startTL   *influxql.TimeLiteral
	stopTL    *influxql.TimeLiteral
}

func NewQueryBuilder(quertStr string, groupByInterval time.Duration) (*QueryBuilder, error) {

	s, err := influxql.ParseStatement(quertStr)
	if err != nil {
		return nil, err
	}
	stmt, ok := s.(*influxql.SelectStatement)
	if !ok {
		return nil, fmt.Errorf("Query must be a select statement '%s'", quertStr)
	}

	if dur, err := stmt.GroupByInterval(); err == nil && dur > 0 {
		return nil, fmt.Errorf("Must not specify the 'GROUP BY time(x)' in query. Use the groupByInterval property.")
	}

	//Add New BinaryExpr for time clause
	startTL := &influxql.TimeLiteral{}
	startExpr := &influxql.BinaryExpr{
		Op:  influxql.GT,
		LHS: &influxql.VarRef{Val: "time"},
		RHS: startTL,
	}

	stopTL := &influxql.TimeLiteral{}
	stopExpr := &influxql.BinaryExpr{
		Op:  influxql.LT,
		LHS: &influxql.VarRef{Val: "time"},
		RHS: stopTL,
	}

	if groupByInterval > 0 {
		stmt.Dimensions = append(stmt.Dimensions,
			&influxql.Dimension{
				Expr: &influxql.Call{
					Name: "time",
					Args: []influxql.Expr{
						&influxql.DurationLiteral{groupByInterval},
					},
				},
			},
		)
	}

	if stmt.Condition != nil {
		stmt.Condition = &influxql.BinaryExpr{
			Op:  influxql.AND,
			LHS: stmt.Condition,
			RHS: &influxql.BinaryExpr{
				Op:  influxql.AND,
				LHS: startExpr,
				RHS: stopExpr,
			},
		}
	} else {
		stmt.Condition = &influxql.BinaryExpr{
			Op:  influxql.AND,
			LHS: startExpr,
			RHS: stopExpr,
		}
	}

	return &QueryBuilder{
		statement: stmt,
		startTL:   startTL,
		stopTL:    stopTL,
	}, nil
}

func (self *QueryBuilder) GetForTimeRange(start, stop time.Time) morgoth.Query {

	self.startTL.Val = start
	self.stopTL.Val = stop

	return morgoth.Query{
		Command: self.statement.String(),
		Start:   start,
		Stop:    stop,
	}
}
