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

#define MAX_ENTRIES 8192

#define TASK_COMM_LEN 16

char __license[] SEC("license") = "Dual MIT/GPL";

struct event_t {
  u32 netns;
  u32 tid;
  u64 ts;
  char comm[TASK_COMM_LEN];
  union {
    u32 saddr_v4;
    u8 saddr_v6[16];
  };
  union {
    u32 daddr_v4;
    u8 daddr_v6[16];
  };
  u32 last_state;
  u32 new_state;
  u32 af;
  u32 lport;
  u32 dport;
};

#define AF_INET 2
#define AF_INET6 10

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(max_entries, MAX_ENTRIES);
  __type(key, struct sock *);
  __type(value, u32);
  __uint(map_flags, BPF_F_NO_PREALLOC);
} senders SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(key_size, sizeof(u32));
  __uint(value_size, sizeof(u32));
} pb SEC(".maps");

struct trace_event_raw_inet_sock_set_state___protocol2 {
  __u16 protocol;
} __attribute__((preserve_access_index));

struct trace_event_raw_inet_sock_set_state___protocol1 {
  u64 __pad; // First 8 bytes are not accessible by bpf code.
  const void *skaddr;
  int oldstate;
  int newstate;
  __u16 sport;
  __u16 dport;
  __u16 family;
  __u8 protocol;
  __u8 saddr[4];
  __u8 daddr[4];
  __u8 saddr_v6[16];
  __u8 daddr_v6[16];
} __attribute__((preserve_access_index));

SEC("tracepoint/sock/inet_sock_set_state")
int trace_inet_sock_set_state(
    struct trace_event_raw_inet_sock_set_state___protocol1 *args) {
  u32 tid = bpf_get_current_pid_tgid();
  if (tid != 0) {
    if (!is_correct_namespace()) {
      return 0;
    }
  }

  if (bpf_core_field_size(args->protocol) == 1) {
    if (args->protocol != IPPROTO_TCP) {
      return 0;
    }
  } else {
    struct trace_event_raw_inet_sock_set_state___protocol2 *case2 =
        (void *)args;
    if (case2->protocol != IPPROTO_TCP) {
      return 0;
    }
  }

  // sk is mostly used as a UUID, and for two tcp stats:
  struct sock *sk = (struct sock *)args->skaddr;

  // lport is either used in a filter here, or later
  u16 lport = args->sport;
  // FILTER_LPORT
  // if (lport != 8000 && lport != 9000) { birth.delete(&sk); return 0; }

  // dport is either used in a filter here, or later
  u16 dport = args->dport;
  // FILTER_DPORT
  // if (dport != 8000 && dport != 9000) { birth.delete(&sk); return 0; }
  int oldstate = args->oldstate;
  int newstate = args->newstate;

  /*
   * This tool includes PID and comm context. It's best effort, and may
   * be wrong in some situations. It currently works like this:
   * - record timestamp on any state < TCP_FIN_WAIT1
   * - cache task context on:
   *       TCP_SYN_SENT: tracing from client
   *       TCP_LAST_ACK: client-closed from server
   * - do output on TCP_CLOSE:
   *       fetch task context if cached, or use current task
   */

  // on enter listen send an event
  if ((newstate == TCP_LISTEN) ||
      ((newstate == TCP_CLOSE) && (oldstate == TCP_SYN_SENT)) ||
      ((oldstate == TCP_CLOSE) && (newstate == TCP_SYN_SENT))) {
    struct event_t e = {};
    e.last_state = oldstate;
    e.new_state = newstate;
    e.netns = get_current_net_ns();
    e.tid = tid;
    e.ts = bpf_ktime_get_ns();
    bpf_get_current_comm(&e.comm, sizeof(e.comm));
    e.af = args->family;
    e.lport = lport;
    e.dport = dport;
    if (args->family == AF_INET) {
      BPF_CORE_READ_INTO(&e.saddr_v4, sk, __sk_common.skc_rcv_saddr);
      BPF_CORE_READ_INTO(&e.daddr_v4, sk, __sk_common.skc_daddr);
    } else {
      BPF_CORE_READ_INTO(&e.saddr_v6, sk,
                         __sk_common.skc_v6_rcv_saddr.in6_u.u6_addr32);
      BPF_CORE_READ_INTO(&e.daddr_v6, sk,
                         __sk_common.skc_v6_daddr.in6_u.u6_addr32);
    }
    bpf_perf_event_output(args, &pb, BPF_F_CURRENT_CPU, &e, sizeof(e));
  }

  return 0;
}
