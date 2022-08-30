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

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(key_size, sizeof(u32));
  __uint(value_size, sizeof(u8));
  __uint(max_entries, 1 << 10);
  __uint(pinning, LIBBPF_PIN_BY_NAME);
} net_ns SEC(".maps");

static u32 get_current_net_ns() {
  struct task_struct *curtask = (struct task_struct *)bpf_get_current_task();

  return BPF_CORE_READ(curtask, nsproxy, net_ns, ns.inum);
}

static inline bool is_correct_namespace() {
  u32 current_ns = get_current_net_ns();
  u8 *expected_ns = (u8 *)bpf_map_lookup_elem(&net_ns, &current_ns);
  if (expected_ns == NULL) {
    return false;
  }
  return *expected_ns == 1;
}
