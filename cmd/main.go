package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed app/*
var appFS embed.FS

type Options struct {
	Path string
	Mod string
	Meth string
}

func ParseOptions() (opts *Options) {
	opts = new(Options)
	flag.StringVar(&opts.Mod, "mod", "", "module name")
	flag.StringVar(&opts.Meth, "meth", "", "overwrite gometh location")

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

func RenderTemplates(opts *Options) error {
	return fs.WalkDir(appFS, "app", func(path string, d fs.DirEntry, err error) error {
		if path == "app" {
			return nil
		}
		projectPath := filepath.Join(opts.Path, filepath.Clean(strings.TrimSuffix(strings.TrimPrefix(path, "app/"), ".tmpl")))
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
		projectFile, err := os.Create(projectPath)
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

func PrepareProject(opts *Options) error {
	commands := [][]string{
		{"go", "mod", "init", opts.Mod},
	}
	if opts.Meth != "" {
		commands = append(commands, []string{"go", "mod", "edit", fmt.Sprintf("-replace=github.com/theaino/gometh=%s", opts.Meth)})
	}
	commands = append(commands, [][]string{
		{"go", "mod", "tidy"},
		{"yarn", "init", "-y"},
		{"yarn", "add", "sass"},
	}...)
	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = opts.Path
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("%s:\n%s", strings.Join(command, " "), string(out))
			return err
		}
	}
	return nil
}

func main() {
	opts := ParseOptions()
	opts.PrepareDirectory()
	if err := RenderTemplates(opts); err != nil {
		log.Fatal(err)
	}
	if err := PrepareProject(opts); err != nil {
		log.Fatal(err)
	}
}
