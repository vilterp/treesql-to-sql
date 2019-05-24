package querygen

import "github.com/vilterp/treesql-to-sql/parse"

func Generate(stmt *parse.Select) string {
	// TODO(vilterp): actually generate
	return "select json_agg(json_build_object('id', id, 'name', name, 'created_at', created_at)) FROM clusters"
}
