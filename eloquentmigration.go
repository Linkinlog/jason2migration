package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"
)

const templateLayout = `<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

{{ .Meta }}

class {{.ActionTitle}}{{.TableTitle }}Table extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up(): void
    {
        Schema::{{.Action}}('{{.Table}}', function (Blueprint $table) {
            {{ .Fields }}
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down(): void
    {
        Schema::dropIfExists('{{.Table}}');
    }
}
`

type EloquentMigration struct{}

func (e *EloquentMigration) InputToMigration(input Input) (migration string) {
	action := "create"
	actionTitle := "Create"
	if !input.Migration.CreateTable {
		action = "table"
		actionTitle = "Update"
	}

	data := struct {
		Table       string
		TableTitle  string
		Action      string
		ActionTitle string
		Fields      string
		Meta        string
	}{
		TableTitle:  fmt.Sprintf("%s%s", strings.ToUpper(string(input.Migration.Table[0])), string(input.Migration.Table[1:])),
		Table:       input.Migration.Table,
		Action:      action,
		ActionTitle: actionTitle,
		Fields:      fieldsToEloquent(input.Migration.Fields),
		Meta:        input.BusinessFacingMeta.String(),
	}

	t := template.Must(template.New("migration").Parse(templateLayout))
	writer := &strings.Builder{}
	err := t.Execute(writer, data)
	if err != nil {
		slog.Error("error executing template", "error", err)
	}
	migration = writer.String()

	return migration
}

func (e *EloquentMigration) ToFile(migration, table string) error {
	now := time.Now()
	t := now.Format(layout)
	fileName := fmt.Sprintf("%s_%s.php", t, table)

	err := os.WriteFile(fileName, []byte(migration), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func fieldsToEloquent(fields []Field) string {
	var b strings.Builder
	for i, field := range fields {
		b.WriteString(fmt.Sprintf("$table->%s('%s')", sqlTypeToEloquent(field.DataType), field.Field))
		if len(field.Constraints) > 0 {
			buildConstraints(&b, field.Constraints)
		} else {
			b.WriteString("->nullable()")
		}
		b.WriteString(";")
		if i < len(fields)-1 {
			b.WriteString("\n            ") // tabs were being weird
		}
	}
	return b.String()
}

func buildConstraints(b *strings.Builder, constraints []string) {
	if !slices.Contains(constraints, "not null") {
		b.WriteString("->nullable()")
	}
	for i, c := range constraints {
		switch c {
		case "not null":
			b.WriteString("->nullable(false)")
		case "auto increment":
			b.WriteString("->autoIncrement()")
		default:
			slog.Warn("constraint not supported", "constraint", c)
			panic("constraint not supported")
		}
		if i < len(constraints)-1 {
			b.WriteString(", ")
		}
	}
}

func sqlTypeToEloquent(t string) string {
	t = strings.ToLower(t)

	var tToE map[string]string = map[string]string{
		"varchar": "string",
	}

	if _, ok := tToE[t]; !ok {
		return strings.ToLower(t)
	}
	return tToE[t]
}
