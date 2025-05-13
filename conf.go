package meth

type ServerConf struct {
	Host string
	Port string
}

type BuildConf struct {
	Esbuild bool
	Sass bool
	Entrypoints []string
	DistDir string
}

type Conf struct {
	Server ServerConf
	Build BuildConf
}
