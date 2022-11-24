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
} oldnamemap SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(key_size, sizeof(u32));
  __uint(value_size, sizeof(char *));
  __uint(max_entries, 1 << 13);
} newnamemap SEC(".maps");

struct enter_rename_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;

  const char *oldname;
  const char *newname;
};

struct enter_renameat_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;

  u64 olddfd;
  const char *oldname;
  u64 newdfd;
  const char *newname;
};

struct enter_link_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;

  const char *oldname;
  const char *newname;
};

struct enter_linkat_info {
  /* not allowed to read */
  unsigned long pad;

  int __syscall_nr;

  u64 olddfd;
  const char *oldname;
  u64 newdfd;
  const char *newname;
  int flags;
};

struct exit_rename_info {
  /* not allowed to read */
  u64 pad;

  int __syscall_nr;
  long ret;
};

struct event_t {
  u32 netns;
  u32 tid;
  u64 ts;
  int syscall;
  char comm[16];
  char oldname[200];
  char newname[200];
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

int renamelink_common(const char *oldname, const char *newname) {
  if (!is_correct_namespace()) {
    return 0;
  }
  if (oldname == NULL || newname == NULL) {
    return 0;
  }
  u32 tid = bpf_get_current_pid_tgid();

  bpf_map_update_elem(&oldnamemap, &tid, &oldname, BPF_ANY);
  bpf_map_update_elem(&newnamemap, &tid, &newname, BPF_ANY);

  return 0;
}

SEC("tracepoint/syscalls/sys_enter_rename")
int sys_enter_rename(struct enter_rename_info *info) {
  return renamelink_common(info->oldname, info->newname);
}

SEC("tracepoint/syscalls/sys_enter_renameat")
int sys_enter_renameat(struct enter_renameat_info *info) {
  return renamelink_common(info->oldname, info->newname);
}

SEC("tracepoint/syscalls/sys_enter_link")
int sys_enter_link(struct enter_link_info *info) {
  return renamelink_common(info->oldname, info->newname);
}

SEC("tracepoint/syscalls/sys_enter_linkat")
int sys_enter_linkat(struct enter_linkat_info *info) {
  return renamelink_common(info->oldname, info->newname);
}

SEC("tracepoint/syscalls/sys_exit_rename")
int sys_exit_rename(struct exit_rename_info *info) {

  u32 tid = bpf_get_current_pid_tgid();
  char **oldname = bpf_map_lookup_elem(&oldnamemap, &tid);
  if (oldname == NULL) {
    return 0;
  }
  if (*oldname == NULL) {
    return 0;
  }

  char **newname = bpf_map_lookup_elem(&newnamemap, &tid);
  if (newname == NULL) {
    return 0;
  }
  if (*newname == NULL) {
    return 0;
  }
  int zero = 0;
  struct event_t *event = bpf_map_lookup_elem(&heap, &zero);
  if (!event) {
    bpf_map_delete_elem(&oldnamemap, &tid);
    bpf_map_delete_elem(&newnamemap, &tid);
    return 0;
  }

  event->netns = get_current_ns();
  event->tid = tid;
  event->ret = info->ret;
  event->ts = bpf_ktime_get_ns();
  event->syscall = info->__syscall_nr;

  bpf_get_current_comm(&event->comm, sizeof(event->comm));
  bpf_probe_read_user_str(event->oldname, sizeof(event->oldname), *oldname);
  bpf_probe_read_user_str(event->newname, sizeof(event->newname), *newname);

  bpf_map_delete_elem(&oldnamemap, &tid);
  bpf_map_delete_elem(&newnamemap, &tid);

  bpf_perf_event_output(info, &pb, BPF_F_CURRENT_CPU, event, sizeof(*event));
  return 0;
}
