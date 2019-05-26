package querygen

import (
	pp "github.com/vilterp/go-pretty-print"
	"github.com/vilterp/treesql-to-sql/parse"
)

func Generate(stmt *parse.Select) (string, error) {
	return generate(stmt).String(), nil
}

func generate(stmt *parse.Select) pp.Doc {
	var sels []pp.Doc
	for _, selOrStar := range stmt.Selections {
		if selOrStar.Star {
			panic("don't know how to do *s yet")
		}
		sel := selOrStar.Selection
		if sel.SubSelect != nil {
			// TODO(vilterp): find foreign key; add where clause
			sels = append(sels, pp.Seq([]pp.Doc{
				pp.Textf("'%s', (", sel.Name), pp.Newline,
				pp.Indent(2, generate(sel.SubSelect)),
				pp.Newline,
				pp.Text(")"),
			}))
		} else {
			colName := sel.Name
			colExpr := colName
			//if sel.Type == TTimestamp {
			//	// RFC3339, which is what jsonpb wants timestamps to be in for it to unmarshal
			//	// them into protobuf timestamps.
			//	colExpr = "experimental_strftime(" + colName + ", '%Y-%m-%dT%H:%M:%S.%fZ')"
			//} else {
			//	colExpr = colName
			//}
			sels = append(sels, pp.Textf("'%s', %s", colName, colExpr))
		}
	}

	selsObj := pp.Seq([]pp.Doc{
		pp.Text("json_build_object("), pp.Newline,
		pp.Indent(2, pp.Join(sels, pp.CommaNewline)), pp.Newline,
		pp.Text(")"),
	})

	if !stmt.One {
		selsObj = pp.Surround("json_agg(", selsObj, ")")
	}

	docs := []pp.Doc{
		pp.Text("SELECT "), selsObj,
	}
	if stmt.Table != "" {
		docs = append(docs, pp.Newline, pp.Text("FROM "), pp.Text(stmt.Table))
	}
	if stmt.Where != nil {
		docs = append(
			docs, pp.Newline, pp.Text("WHERE "),
			pp.Text(stmt.Where.ColumnName), pp.Text(" = "), pp.Text(stmt.Where.Value),
		)
	}
	//if stmt.OrderByClause != "" {
	//	docs = append(docs, pp.Newline, pp.Text("ORDER BY "), pp.Text(stmt.OrderByClause))
	//}

	return pp.Seq(docs)
}
