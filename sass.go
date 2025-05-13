package meth

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func sassPlugin() api.Plugin {
	return api.Plugin{
		Name: "sass",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: `\.scss$`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				nodeBin, err := filepath.Abs("node_modules/.bin")
				if err != nil {
					return api.OnLoadResult{}, err
				}
				cmd := exec.Command(filepath.Join(nodeBin, "sass"), "--no-source-map", "--load-path=node_modules", args.Path)
				cmd.Stderr = os.Stdout
				out, err := cmd.Output()
				if err != nil {
					return api.OnLoadResult{}, err
				}
				css := string(out)

				return api.OnLoadResult{
					Contents: &css,
					Loader: api.LoaderCSS,
					ResolveDir: filepath.Dir(args.Path),
				}, nil
			})
		},
	}
}
