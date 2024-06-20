//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// this list was created by running:
// go list -f '{{.ImportPath}}: {{.GoFiles}}' -tags linux,arm64 syscall
// then eliminating obvious non-linux source filenames (_darwin, _bsd, etc)
var syscallSrcFiles = []string{
	"asan0.go",
	"dirent.go",
	"endian_little.go",
	"env_unix.go",
	"exec_linux.go",
	"exec_unix.go",
	"flock_linux.go",
	"forkpipe2.go",
	"lsf_linux.go",
	"msan0.go",
	"net.go",
	"netlink_linux.go",
	"rlimit.go",
	"rlimit_stub.go",
	"setuidgid_linux.go",
	"sockcmsg_linux.go",
	"sockcmsg_unix.go",
	"sockcmsg_unix_other.go",
	"syscall.go",
	"syscall_linux.go",
	"syscall_linux_accept4.go",
	"syscall_linux_arm64.go",
	"syscall_unix.go",
	"time_nofake.go",
	"timestruct.go",
	"zerrors_linux_arm64.go",
	"zsyscall_linux_arm64.go",
	"zsysnum_linux_arm64.go",
	"ztypes_linux_arm64.go",
}

func main() {
	_, file, _, _ := runtime.Caller(0)
	os.Chdir(filepath.Dir(file))

	cmd := exec.Command("go", "env", "GOROOT")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	syscallDir := fmt.Sprintf("%s/src/syscall", strings.TrimSpace(string(out)))

	for _, filename := range syscallSrcFiles {
		newname := strings.TrimSuffix(filename, filepath.Ext(filename))
		newname = strings.ReplaceAll(newname, "linux", "js")
		newname = strings.ReplaceAll(newname, "arm64", "wasm")
		if !strings.Contains(filename, "linux") {
			// double underscore to avoid conflicts with
			// files that already had a GOOS suffix
			newname = fmt.Sprintf("%s__js", newname)
		}
		data, err := os.ReadFile(filepath.Join(syscallDir, filename))
		if err != nil {
			log.Fatal(err)
		}
		buildTags, ok := findLineByPrefix(data, "//go:build")
		if ok {
			var tags []string
			if strings.Contains(buildTags, "linux") || strings.Contains(buildTags, "unix") {
				tags = append(tags, "js")
			}
			if strings.Contains(buildTags, "arm64") {
				tags = append(tags, "wasm")
			}
			// some known special cases to include
			for _, special := range []string{"!asan", "!msan", "!386", "!arm"} {
				if strings.Contains(buildTags, special) {
					tags = append(tags, special)
				}
			}
			if len(tags) == 0 {
				// for some reason, just say js
				tags = append(tags, "js")
			}
			newTags := fmt.Sprintf("//go:build %s", strings.Join(tags, " && "))
			data = bytes.Replace(data, []byte(buildTags), []byte(newTags), 1)
		}
		if err := os.WriteFile(fmt.Sprintf("./syscall/%s.go", newname), data, 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("wrote", newname)
	}
}

func findLineByPrefix(data []byte, prefix string) (string, bool) {
	prefixBytes := []byte(prefix)
	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if bytes.HasPrefix(line, prefixBytes) {
			return string(line), true
		}
	}

	return "", false
}
