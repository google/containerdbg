//
// Copyright (c) 2022 Google LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// +build ignore

#include "vmlinux_protected.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>

#include "common.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(key_size, sizeof(u32));
  __uint(value_size, sizeof(char *));
  __uint(max_entries, 1 << 13);
} filename_map SEC(".maps");

struct enter_open_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;
  const char *filename;
  int flags;
  umode_t mode;
};

struct enter_openat_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;
  unsigned long dfd;
  const char *filename;
  int flags;
  umode_t mode;
};

struct enter_openat2_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;
  unsigned long dfd;
  const char *filename;
  void *open_how;
  size_t size;
};

struct exit_open_info {
  /* not allowed to read */
  u64 pad;

  int __syscall_nr;
  long ret;
};

struct event_t {
  u32 netns;
  u32 tid;
  u64 ts;
  char comm[16];
  char path[200];
  int ret;
};

/* BPF perfbuf map */
struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(key_size, sizeof(int));
  __uint(value_size, sizeof(int));
} pb SEC(".maps");

/* BPF perf map data store */
struct {
  __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
  __uint(max_entries, 1);
  __type(key, int);
  __type(value, struct event_t);
} heap SEC(".maps");

static void enter_open_common(const char *filename) {
  if (!is_correct_namespace()) {
    return;
  }
  if (filename == NULL) {
    return;
  }
  u32 tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
  bpf_map_update_elem(&filename_map, &tid, &filename, BPF_ANY);
}

SEC("tracepoint/syscalls/sys_enter_open")
int sys_enter_open(struct enter_open_info *info) {
  const char *filename = info->filename;
  enter_open_common(filename);

  return 0;
}

SEC("tracepoint/syscalls/sys_enter_openat")
int sys_enter_openat(struct enter_openat_info *info) {
  const char *filename = info->filename;

  enter_open_common(filename);

  return 0;
}

SEC("tracepoint/syscalls/sys_enter_openat2")
int sys_enter_openat2(struct enter_openat2_info *info) {
  const char *filename = info->filename;

  enter_open_common(filename);

  return 0;
}

static int exit_open_common(void *ctx, u64 ret) {
  u32 tid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
  char **filenameptr = bpf_map_lookup_elem(&filename_map, &tid);
  if (filenameptr == NULL) {
    return 0;
  }
  if (*filenameptr == NULL) {
    return 0;
  }
  int zero = 0;
  struct event_t *event = bpf_map_lookup_elem(&heap, &zero);
  if (!event) {
    bpf_map_delete_elem(&filename_map, &tid);
    return 0;
  }
  event->netns = get_current_ns();
  event->tid = tid;
  event->ret = ret;
  event->ts = bpf_ktime_get_ns();
  bpf_get_current_comm(&event->comm, sizeof(event->comm));
  bpf_probe_read_user_str(event->path, sizeof(event->path), *filenameptr);
  bpf_map_delete_elem(&filename_map, &tid);

  bpf_perf_event_output(ctx, &pb, BPF_F_CURRENT_CPU, event, sizeof(*event));

  return 0;
}

SEC("tracepoint/syscalls/sys_exit_open")
int sys_exit_open(struct exit_open_info *info) {
  return exit_open_common(info, info->ret);
}
