package main

/*
#cgo pkg-config: libgit2

#include<git2.h>

// C macros seem not to be accessible in cgo
git_clone_options clone_opts_init = GIT_CLONE_OPTIONS_INIT;

int just_return_origin(git_remote **out, git_repository *repo, const char *name, const char *url, void *payload)
{
	return git_remote_lookup(out, repo, name);
}

int just_return_repo(git_repository **out, const char *path, int bare, void *payload)
{
	return git_submodule_open(out, (git_submodule*)payload);
}
*/
import "C"
import (
	"errors"
	"path/filepath"
	"unsafe"
)

func init() {
	C.git_libgit2_init()
}

func addSubmodule(url, repoDir, importPath string) error {
	var repo, r *C.git_repository
	var module *C.git_submodule

	curDir := C.CString(".")
	defer C.free(unsafe.Pointer(curDir))
	defer C.git_repository_free(repo)
	if C.git_repository_open(&repo, curDir) < 0 {
		return convertErr(C.giterr_last())
	}

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	submoduleSubpath := filepath.Join("vendor", importPath)
	cSubmoduleSubpath := C.CString(submoduleSubpath)
	defer C.free(unsafe.Pointer(cSubmoduleSubpath))
	defer C.git_submodule_free(module)
	if rc := C.git_submodule_add_setup(&module, repo, cURL, cSubmoduleSubpath, 1); rc < 0 {
		if rc == -4 {
			// submodule exists; job is done
			return nil
		}
		return convertErr(C.giterr_last())
	}

	cloneOpts := C.clone_opts_init
	cloneOpts.repository_cb = (C.git_repository_create_cb)(unsafe.Pointer(C.just_return_repo))
	cloneOpts.remote_cb = (C.git_remote_create_cb)(unsafe.Pointer(C.just_return_origin))
	cloneOpts.repository_cb_payload = unsafe.Pointer(module)
	//cloneOpts.remote_cb_payload = unsafe.Pointer(module);

	cSubmoduleRepoPath := C.CString(filepath.Join(repoDir, submoduleSubpath))
	defer C.free(unsafe.Pointer(cSubmoduleRepoPath))
	if C.git_clone(&r, cURL, cSubmoduleRepoPath, &cloneOpts) < 0 {
		return convertErr(C.giterr_last())
	}
	C.git_repository_free(r)

	if C.git_submodule_add_finalize(module) < 0 {
		return convertErr(C.giterr_last())
	}
	return nil
}

func findGitRoot(path string) (string, error) {
	buf := &C.git_buf{}
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	if C.git_repository_discover(buf, p, 0, nil) < 0 {
		return "", convertErr(C.giterr_last())
	}
	var repo *C.git_repository
	defer C.git_repository_free(repo)
	if C.git_repository_open(&repo, buf.ptr) < 0 {
		return "", convertErr(C.giterr_last())
	}
	return C.GoString(C.git_repository_workdir(repo)), nil
}

func convertErr(err *C.git_error) error {
	if err != nil {
		return errors.New(C.GoString(err.message))
	}
	return errors.New("unknown error")
}
