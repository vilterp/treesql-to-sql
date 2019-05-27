package querygen

import (
	"testing"

	"github.com/vilterp/treesql-to-sql/parse"
	"github.com/vilterp/treesql-to-sql/schema"
	"github.com/xo/xo/models"
)

func TestGenerate(t *testing.T) {
	blogSchema := schema.Schema{
		"posts": &schema.TableWithFKs{
			Table: &models.Table{
				TableName: "posts",
			},
		},
		"comments": &schema.TableWithFKs{
			Table: &models.Table{
				TableName: "comments",
			},
			FKs: []*models.ForeignKey{
				{
					ColumnName:   "post_id",
					RefTableName: "posts",
				},
			},
		},
	}

	query := `MANY posts {
  id,
  comments: MANY comments {
    id
  }
}`
	parsedQuery, err := parse.Parse(query)
	if err != nil {
		t.Fatal(err)
	}

	sqlQuery, err := Generate(parsedQuery.Select, blogSchema)
	if err != nil {
		t.Fatal(err)
	}

	expected := `SELECT json_agg(json_build_object(
  'id', id,
  'comments', (
    SELECT json_agg(json_build_object(
      'id', id
    ))
    FROM comments
    WHERE comments.post_id = posts.id
  )
))
FROM posts`

	if sqlQuery != expected {
		t.Fatalf("expected\n\n%s\n\ngot\n\n%s", expected, sqlQuery)
	}
}
