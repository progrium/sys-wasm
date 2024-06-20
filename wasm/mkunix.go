//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// this list was created by running:
// go list -f '{{.ImportPath}}: {{.GoFiles}}' -tags linux,arm64 ./unix
// then eliminating obvious non-linux source filenames (_darwin, _bsd, etc)
var unixSrcFiles = []string{
	"affinity_linux.go",
	"aliases.go",
	"bluetooth_linux.go",
	"constants.go",
	"dev_linux.go",
	"dirent.go",
	"endian_little.go",
	"env_unix.go",
	"fcntl.go",
	"fdset.go",
	"ifreq_linux.go",
	"ioctl_linux.go",
	"ioctl_unsigned.go",
	"mremap.go",
	"pagesize_unix.go",
	"race0.go",
	"readdirent_getdents.go",
	"sockcmsg_linux.go",
	"sockcmsg_unix.go",
	"sockcmsg_unix_other.go",
	"syscall.go",
	"syscall_linux.go",
	"syscall_linux_arm64.go",
	"syscall_linux_gc.go",
	"syscall_unix.go",
	"syscall_unix_gc.go",
	"sysvshm_linux.go",
	"sysvshm_unix.go",
	"timestruct.go",
	"zerrors_linux.go",
	"zerrors_linux_arm64.go",
	"zptrace_armnn_linux.go",
	"zptrace_linux_arm64.go",
	"zsyscall_linux.go",
	"zsyscall_linux_arm64.go",
	"zsysnum_linux_arm64.go",
	"ztypes_linux.go",
	"ztypes_linux_arm64.go",
}

func main() {
	_, file, _, _ := runtime.Caller(0)
	os.Chdir(filepath.Dir(file))

	for _, filename := range unixSrcFiles {
		newname := strings.TrimSuffix(filename, filepath.Ext(filename))
		newname = strings.ReplaceAll(newname, "linux", "js")
		newname = strings.ReplaceAll(newname, "arm64", "wasm")
		if !strings.Contains(filename, "linux") {
			// double underscore to avoid conflicts with
			// files that already had a GOOS suffix
			newname = fmt.Sprintf("%s__js", newname)
		}
		data, err := os.ReadFile(filepath.Join("../unix", filename))
		if err != nil {
			log.Fatal(err)
		}
		buildTags, ok := findLineByPrefix(data, "//go:build")
		if ok {
			var tags []string
			if strings.Contains(buildTags, "linux") {
				tags = append(tags, "js")
			}
			if strings.Contains(buildTags, "arm64") {
				tags = append(tags, "wasm")
			}
			// some known special cases to include
			for _, special := range []string{"!race", "gc"} {
				if strings.Contains(buildTags, special) {
					tags = append(tags, special)
				}
			}
			newTags := fmt.Sprintf("//go:build %s", strings.Join(tags, " && "))
			data = bytes.Replace(data, []byte(buildTags), []byte(newTags), 1)
		}

		// replace syscall imports
		data = bytes.Replace(data, []byte(`"syscall"`), []byte(`"golang.org/x/sys/wasm/syscall"`), 1)

		if err := os.WriteFile(fmt.Sprintf("../unix/%s.go", newname), data, 0644); err != nil {
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
