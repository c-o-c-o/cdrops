package main

import (
	"cdrops/gcmz"
	"flag"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	params := flag.Args()

	data, err := gcmz.ReadGCMZDropsData()
	if err != nil {
		print(err.Error())
	}

	for _, p := range params {
		v := strings.Split(p, "*")

		layer, err := strconv.Atoi(v[0])
		if err != nil {
			print(err.Error())
			return
		}

		msAdv, err := strconv.Atoi(v[1])
		if err != nil {
			print(err.Error())
			return
		}

		paths := v[2:]

		for i, v := range paths {
			p, err := filepath.Abs(path.Clean(v))
			if err != nil {
				panic(err.Error())
			}
			paths[i] = strings.ReplaceAll(p, "\\", "\\\\")
		}

		gcmz.DropFiles(
			layer,
			msAdv,
			paths,
			data)
	}
}
