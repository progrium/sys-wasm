# wasm support

This fork generates copies of the `golang.org/x/sys/unix` package source files that would be 
used for linux/arm64 and changes their build tags to js/wasm, allowing libraries using `golang.org/x/sys/unix` (often assuming Linux) to build for wasm. 

It also does something similar for the `syscall` standard library, and replaces
imports in the generated source added to `unix` with the generated `wasm/syscall` package here. 

After generating, both have been tuned by hand to sub out empty function bodies, replace internal imports, etc until things compile.

At this point, there is no way to hook into replace stubbed syscalls, but they are at least defined and programs will compile.
