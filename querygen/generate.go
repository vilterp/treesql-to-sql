package querygen

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
	"github.com/vilterp/treesql-to-sql/parse"
	"github.com/vilterp/treesql-to-sql/schema"
)

func Generate(stmt *parse.Select, schema schema.Schema) (string, error) {
	sqlDoc, err := generate(stmt, schema, nil)
	if err != nil {
		return "", err
	}
	return sqlDoc.String(), nil
}

func generate(stmt *parse.Select, schema schema.Schema, outerTable *schema.TableWithFKs) (pp.Doc, error) {
	curTable, ok := schema[stmt.TableName]
	if !ok {
		return nil, fmt.Errorf("no such table: %s", stmt.TableName)
	}
	var sels []pp.Doc
	for _, selOrStar := range stmt.Selections {
		if selOrStar.Star {
			panic("don't know how to do *s yet")
		}
		sel := selOrStar.Selection
		if sel.SubSelect != nil {
			subSelDoc, err := generate(sel.SubSelect, schema, curTable)
			if err != nil {
				return nil, err
			}
			sels = append(sels, pp.Seq([]pp.Doc{
				pp.Textf("'%s', (", sel.Name), pp.Newline,
				pp.Indent(2, subSelDoc),
				pp.Newline,
				pp.Text(")"),
			}))
		} else {
			colName := sel.Name
			colExpr := colName
			// TODO(vilterp): check that column exists
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
	if stmt.TableName != "" {
		docs = append(docs, pp.Newline, pp.Text("FROM "), pp.Text(stmt.TableName))
	}
	if stmt.Where != nil {
		docs = append(
			docs, pp.Newline, pp.Text("WHERE "),
			pp.Text(stmt.Where.ColumnName), pp.Text(" = "), pp.Text(stmt.Where.Value),
		)
	}
	if outerTable != nil {
		if stmt.Where != nil {
			return nil, fmt.Errorf("WHERE and outerTable both set")
		}
		refCol, err := curTable.FindColPointingAt(outerTable.Table.TableName)
		if err != nil {
			return nil, err
		}
		docs = append(
			docs, pp.Newline, pp.Text("WHERE "),
			pp.Textf("%s.%s", curTable.Table.TableName, refCol),
			pp.Text(" = "),
			// TODO(vilterp): actually get PK col name
			pp.Textf("%s.%s", outerTable.Table.TableName, "id"),
		)
	}
	//if stmt.OrderByClause != "" {
	//	docs = append(docs, pp.Newline, pp.Text("ORDER BY "), pp.Text(stmt.OrderByClause))
	//}

	return pp.Seq(docs), nil
}
