package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type MysqlMigration struct{}

func (m *MysqlMigration) InputToMigration(input Input) (migration string) {
	migration += input.BusinessFacingMeta.String()
	migration += "\n\n"

	if input.Migration.CreateTable {
		migration += m.createTable(input)
	} else {
		migration += m.updateTable(input)
	}

	return migration
}

func (m *MysqlMigration) ToFile(migration, table string) error {
	now := time.Now()
	t := now.Format(layout)
	fileName := fmt.Sprintf("%s_%d_%s.sql", t, rand.Int(), table)

	err := os.WriteFile(fileName, []byte(migration), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlMigration) updateTable(input Input) (query string) {
	query += fmt.Sprintf("Alter Table %s\n", input.Migration.Table)
	query += m.updateFields(input.Migration.Fields)
	query += m.dropNewFields(input.Migration.Fields, input.Migration.Table)
	return query
}

func (m *MysqlMigration) updateFields(fields []Field) (query string) {
	for i, field := range fields {
		query += fmt.Sprintf("Add Column %s", field.String())
		if i < len(fields)-1 {
			query += ",\n"
		}
	}
	query += ";\n\n"
	return query
}

func (m *MysqlMigration) dropNewFields(fields []Field, table string) (query string) {
	query += "/*\n"
	for _, field := range fields {
		query += fmt.Sprintf("Alter Table %s Drop Column %s;\n", table, field.Field)
	}
	query += "*/\n"
	return query
}

func (m *MysqlMigration) createTable(input Input) (query string) {
	query += fmt.Sprintf("Create Table if not exists %s\n", input.Migration.Table)
	query += m.createFields(input.Migration.Fields)
	query += fmt.Sprintf("/*\nDrop Table if exists %s;\n*/\n", input.Migration.Table)
	return query
}

func (m *MysqlMigration) createFields(fields []Field) (query string) {
	query += "(\n"
	for i, field := range fields {
		query += fmt.Sprintf("\t%s", field.String())
		if i < len(fields)-1 {
			query += ",\n"
		}
	}
	query += "\n);\n\n"
	return query
}
