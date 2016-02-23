package main

import (
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/tools/go/vcs"
)

type repo struct {
	Path            string
	DependencyPaths []struct{ Root, Prefix string }
}

func NewRepo() (*repo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	repoPath, err := findGitRoot(pwd)
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(repoPath); err != nil {
		return nil, err
	}
	paths := []struct{ Root, Prefix string }{{filepath.Join(runtime.GOROOT(), "src"), ""}}
	return &repo{Path: repoPath, DependencyPaths: paths}, nil
}

func (r *repo) vendorDep(path string) (string, error) {
	rr, err := vcs.RepoRootForImportPath(path, false)
	if err != nil {
		return "", err
	}
	if err := addSubmodule(rr.Repo, filepath.Base(r.Path), rr.Root); err != nil {
		return "", err
	}

	importPath := filepath.FromSlash(path)
	fullPath := filepath.Join(r.Path, "vendor", importPath)
	r.DependencyPaths = append(r.DependencyPaths, struct{ Root, Prefix string }{fullPath, importPath})
	return fullPath, nil
}
