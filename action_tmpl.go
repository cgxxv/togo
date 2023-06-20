package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/urfave/cli"

	"github.com/cgxxv/togo/template"
)

type (
	tmplParams struct {
		Encode  bool
		Package string
		Format  string
		Funcs   string
		Files   []*tmplFile
	}
	tmplFile struct {
		Base string
		Name string
		Path string
		Ext  string
		Data string
	}
)

var tmplCommand = cli.Command{
	Name:   "tmpl",
	Usage:  "embed template files",
	Action: tmplAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "package",
			Value: "template",
		},
		cli.StringFlag{
			Name: "func",
		},
		cli.StringFlag{
			Name:  "input",
			Value: "files/**",
		},
		cli.StringFlag{
			Name:  "output",
			Value: "template_gen.go",
		},
		cli.StringFlag{
			Name:  "format",
			Value: "text",
		},
		cli.StringFlag{
			Name:  "trim-prefix",
			Value: "files/",
		},
		cli.BoolFlag{
			Name: "encode",
		},
	},
}

func tmplAction(c *cli.Context) error {
	pattern := c.Args().First()
	if pattern == "" {
		pattern = c.String("input")
	}

	matches, err := doublestar.Glob(pattern)
	if err != nil {
		return err
	}

	params := tmplParams{
		Encode:  c.Bool("encode"),
		Package: c.String("package"),
		Format:  c.String("format"),
		Funcs:   c.String("func"),
	}

	prefix := c.String("trim-prefix")
	for _, match := range matches {
		stat, oserr := os.Stat(match)
		if oserr != nil {
			return oserr
		}
		if stat.IsDir() {
			continue
		}

		raw, ioerr := ioutil.ReadFile(match)
		if ioerr != nil {
			return ioerr
		}
		params.Files = append(params.Files, &tmplFile{
			Path: match,
			Name: strings.TrimPrefix(match, prefix),
			Base: strings.TrimSuffix(filepath.Base(match), filepath.Ext(match)),
			Ext:  filepath.Ext(match),
			Data: string(raw),
		})
	}

	wr := os.Stdout
	if output := c.String("output"); output != "" {
		wr, err = os.Create(output)
		if err != nil {
			return err
		}
		defer wr.Close()
	}

	return template.Execute(wr, "tmpl.tmpl", params)
}
