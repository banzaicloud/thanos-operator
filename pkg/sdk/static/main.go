package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shurcooL/vfsgen"
)

//go:generate go run main.go
func main() {
	crds := http.Dir(filepath.Join(getRepoRoot(), "config/crd/bases"))

	err := vfsgen.Generate(crds, vfsgen.Options{
		Filename:     filepath.Join(getRepoRoot(), "pkg/sdk/static/gen/crds/generated.go"),
		PackageName:  "crds",
		VariableName: "Root",
		FileModTime:  timePointer(time.Unix(0, 0)),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to generate crds vfs: %+v", err))
	}

	rbac := http.Dir(filepath.Join(getRepoRoot(), "config/rbac"))

	err = vfsgen.Generate(rbac, vfsgen.Options{
		Filename:     filepath.Join(getRepoRoot(), "pkg/sdk/static/gen/rbac/generated.go"),
		PackageName:  "rbac",
		VariableName: "Root",
		FileModTime:  timePointer(time.Unix(0, 0)),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to generate rbac vfs: %+v", err))
	}
}

// getRepoRoot returns the full path to the root of the repo
func getRepoRoot() string {
	_, filename, _, _ := runtime.Caller(0)

	dir := filepath.Dir(filename)

	return filepath.Dir(path.Join(dir, "../.."))
}

func timePointer(t time.Time) *time.Time {
	return &t
}
