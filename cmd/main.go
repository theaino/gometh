package main

import (
	"embed"
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed app/*
var appFS embed.FS

type Options struct {
	Path string
	Mod string
}

func ParseOptions() (opts *Options) {
	opts = new(Options)
	flag.StringVar(&opts.Mod, "mod", "", "module name")

	flag.Parse()
	opts.Path = flag.Arg(0)
	return
}

func (opts *Options) PrepareDirectory() {
	os.MkdirAll(opts.Path, os.ModePerm)
	if opts.Mod == "" {
		opts.Mod = filepath.Base(opts.Path)
	}
}

func PrepareProject(opts *Options) {
}

func RenderTemplates(opts *Options) error {
	return fs.WalkDir(appFS, "app", func(path string, d fs.DirEntry, err error) error {
		projectPath := filepath.Join(opts.Path, filepath.Clean(strings.TrimPrefix(path, "app/")))
		if d.IsDir() {
			return os.MkdirAll(projectPath, os.ModePerm)
		}
		tmplFile, err := appFS.Open(path)
		if err != nil {
			return err
		}
		tmplContents, err := io.ReadAll(tmplFile)
		if err != nil {
			return err
		}
		projectFile, err := os.Open(projectPath)
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".tmpl" {
			tmpl, err := template.New("").Parse(string(tmplContents))
			if err != nil {
				return err
			}
			return tmpl.Execute(projectFile, *opts)
		}
		_, err = projectFile.Write(tmplContents)
		return err
	})
}

func main() {
	opts := ParseOptions()
	opts.PrepareDirectory()
	RenderTemplates(opts)
}
