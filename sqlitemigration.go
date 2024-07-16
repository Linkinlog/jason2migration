package main

import (
	"fmt"
	"os"
	"time"
)

type SqliteMigration struct{}

func (s *SqliteMigration) InputToMigration(input Input) (migration string) {
	migration += input.BusinessFacingMeta.String()
	migration += "\n\n"

	if input.Migration.CreateTable {
		migration += createTable(input)
	} else {
		migration += updateTable(input)
	}

	return migration
}

func (s *SqliteMigration) ToFile(migration, table string) error {
	now := time.Now()
	t := now.Format(layout)
	fileName := fmt.Sprintf("%s_%s.sql", t, table)

	err := os.WriteFile(fileName, []byte(migration), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func updateTable(input Input) (query string) {
	query += fmt.Sprintf("Alter Table %s\n", input.Migration.Table)
	query += updateFields(input.Migration.Fields)
	query += dropNewFields(input.Migration.Fields, input.Migration.Table)
	return query
}

func updateFields(fields []Field) (query string) {
	query += "(\n"
	for _, field := range fields {
		query += fmt.Sprintf("\tAdd Column %s;\n", field.String())
	}
	query += "\n);\n\n"
	return query
}

func dropNewFields(fields []Field, table string) (query string) {
	query += "/*\n"
	for _, field := range fields {
		query += fmt.Sprintf("Alter Table %s Drop Column %s;\n", table, field.Field)
	}
	query += "*/\n"
	return query
}

func createTable(input Input) (query string) {
	query += fmt.Sprintf("Create Table if not exists %s\n", input.Migration.Table)
	query += createFields(input.Migration.Fields)
	query += fmt.Sprintf("/*\nDrop Table if exists %s;\n*/\n", input.Migration.Table)
	return query
}

func createFields(fields []Field) (query string) {
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

func applyConstraints(constraints []string) (field string) {
	for _, constraint := range constraints {
		field += fmt.Sprintf(" %s", constraint)
	}
	return field
}
