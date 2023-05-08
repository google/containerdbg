containerdbg - Building
=======================

TL;DR
=====
In case you want to build your own version of containerdbg or want to build containerdbg images into a registry that is accesible by your cluster run the following:
```
export TARGET_REPO=<your container registry>
make all
```

Dependencies
=====================

## local build
The following libraries are required in order to succefully compile:

1. libbpf-dev
   1. on Debian/Ubuntu can be installed by running `sudo ./test/image/install_libbpf.sh`

1. clang
   1. On Debian/Ubuntu can be installed using `sudo apt-get install clang-13`

1. [ko](https://github.com/google/ko)
   1. can be installed using `go install github.com/google/ko@latest` (make sure `$HOME/go/bin` is in your `PATH`)

1. kpt
   1. `wget https://github.com/GoogleContainerTools/kpt/releases/download/v1.0.0-beta.23/kpt_linux_amd64`
   1. `sudo cp kpt_linux_amd64 /usr/local/bin/`

1. If compiling for local use, Docker is required
   1. https://docs.docker.com/engine/install/ubuntu/
   1. Allow sudoless docker by running `sudo adduser $USER docker`


## container build

Alternatively you could use our test image to build using the following command:

```bash	
docker run -v /var/run/docker.sock:/var/run/docker.sock -v $PWD:/build -w /build eu.gcr.io/modernize-prow/containerdbg-test:latest make all
```

## testing dependencies

The following tools are required to succefully run tests on the project:

1. Docker
1. [Kind](https://kind.sigs.k8s.io/)

Details
================

containerdbg is comprised of two main parts:

* The `containerdbg` binary - the main command line tool - you can build this tool standalone via running `make build-go`. For more information regarding how the binary is built see `./Makefile.go.mk`

* The `containerdbg` images - we use the [ko project](https://github.com/google/ko) to build the images for `containerdbg` - you can build the images by running `make images` - for more infomration see `Makefile.images.mk`

The eBPF filters are built using [bpf2go](https://github.com/cilium/ebpf/tree/master/cmd/bpf2go)

Tests
======
`containerdbg` - utilizes three testing techniques

1. Unit tests
1. Component tests - these can be found in the `test/component` and work by manually creating namespaces and injecting the `containerdbg` filters into them
1. E2E tests - these tests can be found in `test/e2e` and utilize [Kind](https://kind.sigs.k8s.io/) to deploy images to a local cluster and run through full flows

As using eBPF requires setting the `ulimit`, it is recommended to run the component tests in a docker container, there is a helper target to run the e2e test in a container `make test-in-container`

Technical details
=================

BPF CO-RE (Compile Once – Run Everywhere)
-----------------------------------------

Libbpf supports building BPF CO-RE-enabled applications, which, in contrast to
[BCC](https://github.com/iovisor/bcc/), do not require Clang/LLVM runtime
being deployed to target servers and doesn't rely on kernel-devel headers
being available.

It does rely on kernel to be built with [BTF type
information](https://www.kernel.org/doc/html/latest/bpf/btf.html), though.
Some major Linux distributions come with kernel BTF already built in:

  - Fedora 31+
  - RHEL 8.2+
  - OpenSUSE Tumbleweed (in the next release, as of 2020-06-04)
  - Arch Linux (from kernel 5.7.1.arch1-1)
  - Manjaro (from kernel 5.4 if compiled after 2021-06-18)
  - Ubuntu 20.10
  - Debian 11 (amd64/arm64)

If your kernel doesn't come with BTF built-in, you'll need to build custom
kernel. You'll need:
  - `pahole` 1.16+ tool (part of `dwarves` package), which performs DWARF to
    BTF conversion;
  - kernel built with `CONFIG_DEBUG_INFO_BTF=y` option;
  - you can check if your kernel has BTF built-in by looking for
    `/sys/kernel/btf/vmlinux` file:

vmlinux.h generation
-------------------

vmlinux.h contains all kernel types, both exported and internal-only. BPF
CO-RE-based applications are expected to include this file in their BPF
program C source code to avoid dependency on kernel headers package.

For more reproducible builds, vmlinux.h header file is pre-generated and
checked in along the other sources. This is done to avoid dependency on
specific user/build server's kernel configuration, because vmlinux.h
generation depends on having a kernel with BTF type information built-in
(which is enabled by `CONFIG_DEBUG_INFO_BTF=y` Kconfig option See below).

vmlinux.h is generated from upstream Linux version at particular minor
version tag. E.g., `vmlinux_505.h` is generated from v5.5 tag. Exact set of
types available in compiled kernel depends on configuration used to compile
it. To generate present vmlinux.h header, default configuration was used, with
only extra `CONFIG_DEBUG_INFO_BTF=y` option enabled.

The command used for generating the header was:
bpftool btf dump file /sys/kernel/btf/vmlinux format c > pkg/ebpf/headers/vmlinux_612.h

Given different kernel version can have incompatible type definitions, it
might be important to use vmlinux.h of a specific kernel version as a "base"
version of header. To that extent, all vmlinux.h headers are versioned by
appending <MAJOR><MINOR> suffix to a file name. There is always a symbolic
link vmlinux.h, that points to whichever version is deemed to be default
(usually, latest).
