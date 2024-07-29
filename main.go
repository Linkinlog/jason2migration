package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

const (
	layout = "20060102150405"
)

func main() {
	sFlag := flag.String("s", "sqlite", "Strategy (sqlite, eloquent, mysql)")
	inputFile := flag.String("f", "input.json", "Input file")
	write := flag.Bool("w", false, "Write to file")
	flag.Parse()

	var s Strategy

	switch *sFlag {
	case "sqlite":
		s = &SqliteMigration{}
	case "eloquent":
		s = &EloquentMigration{}
	case "mysql":
		s = &MysqlMigration{}
	}

	inputs, err := createInputs(*inputFile)
	if err != nil {
		slog.Error("main error", "error", err)
		return
	}

	for _, input := range inputs {
		if *write {
			err := s.ToFile(s.InputToMigration(input), input.Migration.Table)
			if err != nil {
				slog.Error("main error", "error", err)
			}
		} else {
			fmt.Println(s.InputToMigration(input))
		}
	}
}

func createInputs(file string) ([]Input, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("createInputs: readFile error: %w", err)
	}

	var inputs []Input
	err = json.Unmarshal(fileBytes, &inputs)
	if err != nil {
		return nil, fmt.Errorf("createInputs: unmarshal error: %w", err)
	}
	return inputs, nil
}
