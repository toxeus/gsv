package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// TODO: replace this dependency by stdlib
	"github.com/constabulary/gb/vendor"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s url|import path\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	log.SetFlags(0)

	switch len(flag.Args()) {
	case 0:
		log.Fatal("import path missing")
	case 1:
		if err := run(flag.Arg(0)); err != nil {
			log.Fatal(err.Error())
		}
	default:
		log.Fatal("more than one import path supplied")
	}
}

func run(path string) error {
	repo, err := NewRepo()
	if err != nil {
		return err
	}
	fullPath, err := repo.vendorDep(strings.TrimSuffix(stripscheme(path), ".git"))
	if err != nil {
		return err
	}

	for {
		dsm, err := vendor.LoadPaths(repo.DependencyPaths...)
		if err != nil {
			return err
		}

		is, ok := dsm[fullPath]
		if !ok {
			return fmt.Errorf("unable to locate depset for %q", path)
		}

		pkg := findMissing(pkgs(is.Pkgs), dsm)
		if pkg == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "analyzing recursive dependency %s\n", *pkg)
		if _, err := repo.vendorDep(*pkg); err != nil {
			return err
		}
	}

	return nil
}
