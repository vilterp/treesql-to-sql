package querygen

import (
	"fmt"

	"github.com/vilterp/go-parserlib/examples/treesql"
	pp "github.com/vilterp/go-pretty-print"
	"github.com/vilterp/treesql-to-sql/schema"
)

func Generate(stmt *treesql.Select, schema schema.Schema) (string, error) {
	sqlDoc, err := generate(stmt, schema, nil)
	if err != nil {
		return "", err
	}
	return sqlDoc.String(), nil
}

func generate(stmt *treesql.Select, schema schema.Schema, outerTable *schema.TableWithFKs) (pp.Doc, error) {
	curTable, ok := schema[stmt.TableName.Text]
	if !ok {
		return nil, fmt.Errorf("no such table: %s", stmt.TableName)
	}
	var sels []pp.Doc
	for _, sel := range stmt.Selections {
		if sel.SubSelect != nil {
			subSelDoc, err := generate(sel.SubSelect, schema, curTable)
			if err != nil {
				return nil, err
			}
			sels = append(sels, pp.Seq([]pp.Doc{
				pp.Textf("'%s', (", sel.Name.Text), pp.Newline,
				pp.Indent(2, subSelDoc),
				pp.Newline,
				pp.Text(")"),
			}))
		} else {
			colName := sel.Name.Text
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

	if stmt.Many {
		selsObj = pp.Surround("json_agg(", selsObj, ")")
	}

	docs := []pp.Doc{
		pp.Text("SELECT "), selsObj,
	}
	if stmt.TableName.Text != "" {
		docs = append(docs, pp.Newline, pp.Text("FROM "), pp.Text(stmt.TableName.Text))
	}
	// TODO(vilterp): in parser: do where clauses
	//if stmt.Where != nil {
	//	docs = append(
	//		docs, pp.Newline, pp.Text("WHERE "),
	//		// TODO(vilterp): actually escape this
	//		pp.Text(stmt.Where.ColumnName), pp.Text(" = "), pp.Textf(`'%s'`, stmt.Where.Value),
	//	)
	//}
	// join
	// TODO(vilterp): actually get PK col name
	//   need to get this in schema package
	if outerTable != nil {
		//if stmt.Where != nil {
		//	return nil, fmt.Errorf("WHERE and outerTable both set")
		//}
		if !stmt.Many {
			// outer table pointing at this
			refCol, err := outerTable.FindColPointingAt(curTable.Table.TableName)
			if err != nil {
				return nil, err
			}
			docs = append(
				docs, pp.Newline, pp.Text("WHERE "),
				pp.Textf("%s.%s", outerTable.Table.TableName, refCol),
				pp.Text(" = "),
				pp.Textf("%s.%s", curTable.Table.TableName, "id"),
			)
		} else {
			// this pointing at outer table
			refCol, err := curTable.FindColPointingAt(outerTable.Table.TableName)
			if err != nil {
				return nil, err
			}
			docs = append(
				docs, pp.Newline, pp.Text("WHERE "),
				pp.Textf("%s.%s", curTable.Table.TableName, refCol),
				pp.Text(" = "),
				pp.Textf("%s.%s", outerTable.Table.TableName, "id"),
			)
		}
	}
	//if stmt.OrderByClause != "" {
	//	docs = append(docs, pp.Newline, pp.Text("ORDER BY "), pp.Text(stmt.OrderByClause))
	//}

	return pp.Seq(docs), nil
}
