package main

import "fmt"

type Input struct {
	BusinessFacingMeta BusinessFacingMeta `json:"businessFacingMeta"`
	Migration          Migration          `json:"migration"`
}

type BusinessFacingMeta struct {
	CreationDate    string `json:"creationDate"`
	Version         string `json:"version"`
	JiraTicket      string `json:"jiraTicket"`
	BusinessPurpose string `json:"businessPurpose"`
}

type Migration struct {
	Table   string  `json:"table"`
	Fields  []Field `json:"fields"`
	Indexes []Index `json:"indexes"`
}

type Field struct {
	Field       string   `json:"field"`
	DataType    string   `json:"dataType"`
	Constraints []string `json:"constraints"`
}

func (f *Field) String() string {
	return fmt.Sprintf("\t%s %s%s", f.Field, f.DataType, applyConstraints(f.Constraints))
}

type Index struct {
	IndexName string   `json:"indexName"`
	Fields    []string `json:"fields"`
	Unique    bool     `json:"unique"`
}

type Strategy interface {
	InputToMigration(input Input) string
	ToFile(migration, table string) error
}
