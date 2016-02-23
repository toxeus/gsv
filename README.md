# gsv

## Table of Contents

* [Introduction](#introduction)
* [Dependencies](#dependencies)
* [Installation](#installation)
* [Usage](#usage)
* [FAQ](#faq)
* [TODO](#todo)
* [Credits](#credits)

## Introduction

gsv is the **G**o **S**ubmodule **V**endoring tool. It does native
Go vendoring using Git submodules. This approach makes configuration
files redundant and doesn't require additional tooling to build a
gsv-vendored Go project because Git (which you have installed anyway)
is used to track the revisions of your vendored dependencies.

Therefore, in order to fetch and install a Go package that has used
`gsv` to vendor its dependencies a simple

```sh
go get $PACKAGE_URL
```

will do the job. Go 1.5.X needs some additional steps (see FAQ).

Compared to a copy-based vendoring approach gsv preserves the dependencies'
histories and links them to the main project's history which facilitates
the usage of the Git tool suite. Therefore, it's possible to go through
each dependency's `git log` and analyze what changes have been done since
the last update. Also, `git bisect` can be used to find a commit in a
dependency's repository that for example caused the main project to break.
And so on.

## Dependencies

[libgit2](https://github.com/libgit2/libgit2) needs to be installed.
Packages exist for

* [Ubuntu](http://packages.ubuntu.com/source/wily/libgit2)
* [Fedora](https://admin.fedoraproject.org/pkgdb/package/rpms/libgit2/)
* [Arch Linux](https://www.archlinux.org/packages/extra/x86_64/libgit2/)
* [Homebrew](https://github.com/Homebrew/homebrew/blob/master/Library/Formula/libgit2.rb)

If your distro does not package it then you need to install it from
[source](https://github.com/libgit2/libgit2#building-libgit2---using-cmake).

## Installation

After installing the [dependencies](#dependencies) execute

```sh
go get github.com/toxeus/gsv
```

and make sure that `$GOPATH/bin` is in your `$PATH`.

## Usage

Let's assume we want to vendor `go-etcd` and its recursive dependencies
in our project

```sh
cd $GOPATH/$OUR_PROJECT
gsv github.com/coreos/go-etcd/etcd
git commit -m "vendor: added go-etcd and its dependencies"
```

Done.

## FAQ

### I have a problem. What next?

Please open an [issue](https://github.com/toxeus/gsv/issues/new)
and try to give reproducible examples ;)

### I want to contribute. What next?

If you want to contribute a bigger change then
please open an [issue](https://github.com/toxeus/gsv/issues/new)
such that we can discuss what the best way is to proceed.

If you want to fix something minor then feel free to open
a pull request straightaway.

### How do I update all my dependencies?

As of now using

```sh
git submodule foreach 'git fetch && git rebase master@{u}'
```

will do the trick. Some day it'll be integrated into `gsv`.

### Do you know about git2go?

Yes! Last time I checked the support for submodules in
[`git2go`](https://github.com/libgit2/git2go) was not sufficient
for this project's requirements. That's why the Git code
ended up being written in Cgo.

### How do I get and build a gsv project using Go 1.5?

Like this

```sh
export GO15VENDOREXPERIMENT=1
go get -d $PROJECT_URL
cd $GOPATH/$PROJECT_PATH
git submodule update --init
go install ./...
```

### How do I build gsv using Go 1.5?

See [here](#how-do-i-get-and-build-a-gsv-project-using-go-15)

### `go build ./...` or `go test ./...` fails

`gsv` only vendors dependencies which are needed to
satisfy the building and testing of **your** project.
Dependencies that are needed to build and test
the vendored dependencies are **not** pulled in.

As a consequence `go build ./...` and/or `go test ./...`
*might* break when run from the project's *root folder*.
To fix this the following commands should be used

```sh
go build $(go list ./... | grep -v vendor)
go test $(go list ./... | grep -v vendor)
```

This is not a joke but the
[official recommendation](https://github.com/golang/go/issues/11659#issuecomment-171678025)
by the golang team.

Note that pulling in the build- and test-dependencies of your
dependencies is unlikely to fix `go build ./...` and `go test ./...`
because the pulled in packages will then have unsatisfied
dependencies. And going all the way down in the recursion doesn't
seem to be the right solution for managing your build- and
test-dependencies.

## TODO

In order of priority. Might not be implemented exactly as
suggested here.

1. Add a `-purge` flag such that unused vendored submodules will
   be detected and removed.
1. Use code from Go's stdlib instead from gb-vendor
1. Add a `-updateall` flag such that all dependencies are
   updated to their current `origin/master`.
1. Take a look at `git2go` and see if there is progress in
   submodules support.
1. Look for alternatives for `libgit2`. This dependency reduces
   portability and makes installation a bit harder.

## Credits

This project could bootstrap on the work done for
[gb-vendor](https://github.com/constabulary/gb/tree/master/cmd/gb-vendor).
