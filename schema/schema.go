package schema

import (
	"database/sql"
	"fmt"

	"github.com/xo/xo/loaders"
	"github.com/xo/xo/models"
)

type TableWithFKs struct {
	Table *models.Table
	// TODO(vilterp): get columns
	FKs     []*models.ForeignKey
	Columns []*models.Column
}

func (t *TableWithFKs) FindColPointingAt(table string) (string, error) {
	for _, fk := range t.FKs {
		if fk.RefTableName == table {
			return fk.ColumnName, nil
		}
	}
	return "", fmt.Errorf("no col in %v found pointing at %v", t.Table.TableName, table)
}

type Schema map[string]*TableWithFKs

func LoadSchema(conn *sql.DB) (Schema, error) {
	out := Schema{}
	tables, err := loaders.PgTables(conn, "public", "r")
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		foreignKeys, err := models.PgTableForeignKeys(conn, "public", table.TableName)
		if err != nil {
			return nil, err
		}
		columns, err := models.PgTableColumns(conn, "public", table.TableName, false)
		if err != nil {
			return nil, err
		}
		out[table.TableName] = &TableWithFKs{
			Table:   table,
			FKs:     foreignKeys,
			Columns: columns,
		}
	}
	return out, nil
}
