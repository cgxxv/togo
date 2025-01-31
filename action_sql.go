package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"

	"github.com/cgxxv/togo/parser"
	"github.com/cgxxv/togo/template"
)

type sqlParams struct {
	Package    string
	Dialect    string
	Tags       string
	Statements []*parser.Statement
}

var sqlCommand = cli.Command{
	Name:   "sql",
	Usage:  "embed sql statements",
	Action: sqlAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "package",
			Value: "sql",
		},
		cli.StringFlag{
			Name:  "dialect",
			Value: "sqlite3",
		},
		cli.StringFlag{
			Name:  "input",
			Value: "files/*.sql",
		},
		cli.StringFlag{
			Name:  "output",
			Value: "sql_gen.go",
		},
		cli.StringSliceFlag{
			Name: "tags",
		},
	},
}

func sqlAction(c *cli.Context) error {
	pattern := c.Args().First()
	if pattern == "" {
		pattern = c.String("input")
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	params := sqlParams{
		Package: c.String("package"),
		Dialect: c.String("dialect"),
		Tags:    strings.Join(c.StringSlice("tags"), " "),
	}

	parse := parser.New()
	for _, match := range matches {
		statements, perr := parse.ParseFile(match)
		if perr != nil {
			return perr
		}
		params.Statements = append(params.Statements, statements...)
	}

	wr := os.Stdout
	if output := c.String("output"); output != "" {
		wr, err = os.Create(output)
		if err != nil {
			return err
		}
		defer wr.Close()
	}

	return template.Execute(wr, "sql.tmpl", params)
}
