# TreeSQL to SQL

![screenshot](images/image.png)

_TreeSQL UI Console._
- _Schema on the left_
- _TreeSQL editor and completions in the middle_
- _generated SQL on the top right_
- _JSON results on the bottom right_

## What is in this repo?

A server which provides:

- An HTTP endpoint which executes TreeSQL queries by translating them to SQL queries. It returns
  results as JSON.
- A simple console UI (pictured) which allows authoring TreeSQL queries (with schema-aware
  autocomplete), executing them, and seeing results.

## What is TreeSQL?

A new query language which is halfway between SQL and GraphQL.

It's like SQL, except you write joins by starting with one table, and nesting the tables
you want to join to within curly braces, forming a tree of joins. Results come back as a tree
of JSON mirroring the query try.

It's like GraphQL, except that it has a direct mapping onto a DB schema (i.e. there's no arbitrary
code between the query and the data storage, as there is in GraphQL servers) and it has a more
SQL-like syntax (WHERE clauses, etc).

## Why another query language?

GraphQL gives developers the tree-structured queries and results they want, but is backend-agnostic.
Developers struggle to provide efficient GraphQL APIs on top of existing databases, whether they're 
relational or NoSQL.

SQL is implemented by many sophisticated database systems, but doesn't allow developers to directly
express the queries they need to fetch the data to render their UIs. ORMs attempt to bridge this gap,
but have clumsy support for efficiently loading graphs of objects across multiple tables.

TreeSQL attempts to bridge this gap, by offering tree-structured queries in a way that maps directly
to a conventional relational DB schema. This prototype executes queries by translating them to SQL,
but TreeSQL could potentially be executed directly by the database engine -- either incorporating it
into the SQL grammar, or as a separate syntax altogether.
  
## How is TreeSQL translated to SQL?

By translating it to a correlated subquery, which uses Postgres/Cockroach's builtin JSON
functions to combine all results into a single JSON datum -- a one-row, one-column result set. 

## What's missing?

A lot of things.
- `GROUP BY` (unclear what the syntax should be for this)
- `ORDER BY`
- Joining from one to one (only supports one to many)
- etc

## Related work

- [EdgeDB's EdgeQL](https://edgedb.com/) is another attempt to find a hybrid syntax between GraphQL and SQL
- [Hasura](https://hasura.io/) generates correlated subqueries from GraphQL queries, I believe
