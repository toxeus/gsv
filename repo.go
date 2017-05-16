package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/libgit2/git2go"
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
	repoPath, err := git.Discover(pwd, false, nil)
	if err != nil {
		return nil, err
	}

	if err := os.Chdir(repoPath); err != nil {
		return nil, err
	}
	paths := []struct{ Root, Prefix string }{{filepath.Join(runtime.GOROOT(), "src"), ""}}
	return &repo{Path: strings.TrimSuffix(repoPath, ".git/"), DependencyPaths: paths}, nil
}

func (r *repo) vendorDep(path string) (string, error) {
	rr, err := vcs.RepoRootForImportPath(path, false)
	if err != nil {
		return "", err
	}
	importPath := filepath.FromSlash(path)
	subPath := filepath.Join("vendor", importPath)
	rep, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}
	submodule, err := rep.Submodules.Add(rr.Repo, subPath, true)
	if git.IsErrorCode(err, git.ErrExists) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if _, err = submodule.Clone(); err != nil {
		return "", err
	}
	if err := submodule.FinalizeAdd(); err != nil {
		return "", err
	}

	fullPath := filepath.Join(r.Path, subPath)
	r.DependencyPaths = append(r.DependencyPaths, struct{ Root, Prefix string }{fullPath, importPath})
	return fullPath, nil
}
