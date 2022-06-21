package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var base string
var Configs = make([]Config, 0)
var ConfigName = "config.json"

func init() {
	base = os.Getenv("DOCKER_CONFIG_PATH")
	if base != "" {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	base = filepath.Join(home, ".docker")
}

type Config struct {
	Name string
	Path string
}

func main() {
	var ConfigPath = filepath.Join(base, ConfigName)
	err := filepath.WalkDir(base, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() && path != base {
			return fs.SkipDir
		}
		filename := filepath.Base(path)
		if strings.HasPrefix(filename, ConfigName) {
			if filename != ConfigName {
				ext := filepath.Ext(filename)[1:]
				Configs = append(Configs, Config{ext, path})
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	index := 0
	names := make([]string, len(Configs))
	for i, m := range Configs {
		names[i] = m.Name
	}
	var qs = &survey.Select{
		Message: "Choose a Config:",
		Options: names,
		Description: func(value string, index int) string {
			return Configs[index].Path
		},
	}

	err = survey.AskOne(qs, &index)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(Configs[index])
	fmt.Println(ConfigPath)
	os.Remove(ConfigPath)
	err = os.Link(Configs[index].Path, ConfigPath)
	if err != nil {
		panic(err)
	}
}
