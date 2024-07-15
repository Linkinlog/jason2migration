package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/template"
	"time"
)

const templateLayout = `<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

{{ .Meta }}

class Create{{.TableTitle }}Table extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
	public function up(): void
    {
        Schema::create('{{.Table}}', function (Blueprint $table) {
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
	t := template.Must(template.New("migration").Parse(templateLayout))
	data := struct {
		Table      string
		TableTitle string
		Fields     string
		Meta       string
	}{
		TableTitle: fmt.Sprintf("%s%s", strings.ToUpper(string(input.Migration.Table[0])), string(input.Migration.Table[1:])),
		Table:      input.Migration.Table,
		Fields:     fieldsToEloquent(input.Migration.Fields),
		Meta:       applyMetadata(input.BusinessFacingMeta),
	}

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
			for _, constraint := range field.Constraints {
				b.WriteString(fmt.Sprintf("->%s()", constraintToEloquent(constraint)))
			}
		}
		b.WriteString(";")
		if i < len(fields)-1 {
			b.WriteString("\n            ") // tabs were being weird
		}
	}
	return b.String()
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

func constraintToEloquent(c string) string {
	c = strings.ToLower(c)

	var cToE map[string]string = map[string]string{
		"not null": "nullable",
	}

	if _, ok := cToE[c]; !ok {
		return c
	}
	return cToE[c]
}
