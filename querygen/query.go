package querygen

import (
	"sort"

	pp "github.com/cockroachlabs/management-console/util/prettyprint"
)

type Node struct {
	TableName     string
	Columns       []ColumnDesc
	One           bool
	Children      map[string]*Node
	Where         string
	OrderByClause string
}

func (desc TableDesc) All() *Node {
	return &Node{
		TableName: desc.Name,
		Columns:   desc.Cols,
	}
}

func (desc TableDesc) One(where string) *Node {
	return &Node{
		TableName: desc.Name,
		Columns:   desc.Cols,
		One:       true,
		Where:     where,
	}
}

func (n *Node) WithChildren(children map[string]*Node) *Node {
	// TODO(vilterp): not mutate?
	n.Children = children
	return n
}

func (n *Node) OrderBy(orderBy string) *Node {
	n.OrderByClause = orderBy
	return n
}

func (n Node) ToSQL() pp.Doc {
	var sels []pp.Doc
	for _, col := range n.Columns {
		colName := col.Name
		var colExpr string
		if col.Type == TTimestamp {
			// RFC3339, which is what jsonpb wants timestamps to be in for it to unmarshal
			// them into protobuf timestamps.
			colExpr = "experimental_strftime(" + colName + ", '%Y-%m-%dT%H:%M:%S.%fZ')"
		} else {
			colExpr = colName
		}
		sels = append(sels, pp.Textf("'%s', %s", colName, colExpr))
	}

	sortedKeys := make([]string, len(n.Children))
	i := 0
	for key := range n.Children {
		sortedKeys[i] = key
		i++
	}
	sort.Strings(sortedKeys)

	for _, childKey := range sortedKeys {
		child := n.Children[childKey]
		sels = append(sels, pp.Seq([]pp.Doc{
			pp.Textf("'%s', (", childKey), pp.Newline,
			pp.Indent(2, child.ToSQL()),
			pp.Newline,
			pp.Text(")"),
		}))
	}

	selections := pp.Seq([]pp.Doc{
		pp.Text("json_build_object("), pp.Newline,
		pp.Indent(2, pp.Join(sels, pp.CommaNewline)), pp.Newline,
		pp.Text(")"),
	})

	if !n.One {
		selections = pp.Surround("json_agg(", selections, ")")
	}

	docs := []pp.Doc{
		pp.Text("SELECT "), selections,
	}
	if n.TableName != "" {
		docs = append(docs, pp.Newline, pp.Text("FROM "), pp.Text(n.TableName))
	}
	if n.Where != "" {
		docs = append(docs, pp.Newline, pp.Text("WHERE "), pp.Text(n.Where))
	}
	if n.OrderByClause != "" {
		docs = append(docs, pp.Newline, pp.Text("ORDER BY "), pp.Text(n.OrderByClause))
	}

	return pp.Seq(docs)
}
