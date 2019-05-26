package parse

import (
	"testing"
)

type testCase struct {
	in  string
	out string
}

func TestParser(t *testing.T) {
	testCases := []testCase{
		//`CREATETABLE blog_posts (id STRING PRIMARYKEY, title STRING, author_id STRING REFERENCESTABLE blog_posts)`,

		{
			`MANY blog_posts { id, body, comments: MANY comments { id, body } }`,
			`MANY blog_posts {
  id,
  body,
  comments: MANY comments {
    id,
    body
  }
}`},
		{
			`ONE blog_posts WHERE id = "5" { id, title }`,
			`ONE blog_posts WHERE id = "5" {
  id,
  title
}`,
		},

		//`UPDATE blog_posts SET title = "bloop" WHERE id = "5"`,

		//`INSERT INTO blog_posts VALUES ("5", "bloop_doop")`,
	}

	for _, testCase := range testCases {
		statement, err := Parse(testCase.in)
		if err != nil {
			t.Fatal(err)
		}
		formatted := statement.Pretty().String()
		if formatted != testCase.out {
			t.Fatalf("parsed\n\n%s\n\nand it formatted back to\n\n%s", testCase.in, formatted)
		}
	}
}
