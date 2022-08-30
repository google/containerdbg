// Copyright 2021 Google LLC All Rights Reserved.
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

  event->netns = get_current_net_ns();
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
