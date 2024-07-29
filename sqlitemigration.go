package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type SqliteMigration struct{}

func (s *SqliteMigration) InputToMigration(input Input) (migration string) {
	migration += input.BusinessFacingMeta.String()
	migration += "\n\n"

	if input.Migration.CreateTable {
		migration += s.createTable(input)
	} else {
		migration += s.updateTable(input)
	}

	return migration
}

func (s *SqliteMigration) ToFile(migration, table string) error {
	now := time.Now()
	t := now.Format(layout)
	fileName := fmt.Sprintf("%s_%d_%s.sql", t, rand.Int(), table)

	err := os.WriteFile(fileName, []byte(migration), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteMigration) updateTable(input Input) (query string) {
	query += s.updateFields(input.Migration.Fields, input.Migration.Table)
	query += s.dropNewFields(input.Migration.Fields, input.Migration.Table)
	return query
}

func (s *SqliteMigration) updateFields(fields []Field, table string) (query string) {
	for _, field := range fields {
		query += fmt.Sprintf("Alter Table %s Add Column %s;\n", table, field.String())
	}
	query += "\n\n"
	return query
}

func (s *SqliteMigration) dropNewFields(fields []Field, table string) (query string) {
	query += "/*\n"
	for _, field := range fields {
		query += fmt.Sprintf("Alter Table %s Drop Column %s;\n", table, field.Field)
	}
	query += "*/\n"
	return query
}

func (s *SqliteMigration) createTable(input Input) (query string) {
	query += fmt.Sprintf("Create Table if not exists %s\n", input.Migration.Table)
	query += s.createFields(input.Migration.Fields)
	query += fmt.Sprintf("/*\nDrop Table if exists %s;\n*/\n", input.Migration.Table)
	return query
}

func (s *SqliteMigration) createFields(fields []Field) (query string) {
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
