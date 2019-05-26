package parse

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

func (s *Statement) Pretty() pp.Doc {
	if s.Select != nil {
		return s.Select.Pretty()
	}
	panic("only know how to format selects")
}

func (s *Select) Pretty() pp.Doc {
	var initial pp.Doc
	if s.One {
		initial = pp.Text("ONE")
	} else {
		initial = pp.Text("MANY")
	}
	var selDocs []pp.Doc
	for _, selOrStar := range s.Selections {
		if selOrStar.Star {
			selDocs = append(selDocs, pp.Text("*"))
		} else {
			selDocs = append(selDocs, selOrStar.Selection.Pretty())
		}
	}
	maybeLive := ""
	if s.Live {
		maybeLive = " LIVE"
	}
	maybeWhere := ""
	if s.Where != nil {
		maybeWhere = fmt.Sprintf(` WHERE %s = "%s"`, s.Where.ColumnName, s.Where.Value)
	}
	return pp.Seq([]pp.Doc{
		initial,
		pp.Text(" "),
		pp.Text(s.TableName),
		pp.Text(maybeWhere),
		pp.Text(" {"), pp.Newline,
		pp.Indent(2, pp.Join(selDocs, pp.CommaNewline)),
		pp.Newline,
		pp.Text("}"),
		pp.Text(maybeLive),
	})
}

func (s *Selection) Pretty() pp.Doc {
	if s.SubSelect == nil {
		return pp.Text(s.Name)
	}
	return pp.Seq([]pp.Doc{
		pp.Text(s.Name),
		pp.Text(": "),
		s.SubSelect.Pretty(),
	})
}
