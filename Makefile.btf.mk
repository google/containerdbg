
.PHONY: btf-headers
btf-headers: pkg/ebpf/headers/vmlinux.h

pkg/ebpf/headers:
	mkdir $@

pkg/ebpf/headers/vmlinux.h: pkg/ebpf/headers
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > $@

