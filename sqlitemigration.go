package main

import (
	"fmt"
	"os"
	"time"
)

type SqliteMigration struct{}

func (s *SqliteMigration) InputToMigration(input Input) (migration string) {
	migration += applyMetadata(input.BusinessFacingMeta)
	migration += "\n\n"

	migration += fmt.Sprintf("Create Table if not exists %s\n", input.Migration.Table)
	migration += handleFields(input.Migration.Fields)

	migration += fmt.Sprintf("/*\nDrop Table if exists %s;\n*/\n", input.Migration.Table)

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

func handleFields(fields []Field) (query string) {
	query += "(\n"
	for i, field := range fields {
		query += field.String()
		if i < len(fields)-1 {
			query += ",\n"
		}
	}
	query += "\n);\n\n"
	return query
}

func applyMetadata(meta BusinessFacingMeta) (metadata string) {
	metadata = fmt.Sprintf("/*\nCreation Date: %s\nVersion: %s\nJira Ticket: %s\nBusiness Purpose: %s\n*/", meta.CreationDate, meta.Version, meta.JiraTicket, meta.BusinessPurpose)
	return metadata
}

func applyConstraints(constraints []string) (field string) {
	for _, constraint := range constraints {
		field += fmt.Sprintf(" %s", constraint)
	}
	return field
}
