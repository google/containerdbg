// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
  u32 tid = bpf_get_current_pid_tgid();
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
  u32 tid = bpf_get_current_pid_tgid();
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
  event->netns = get_current_net_ns();
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
