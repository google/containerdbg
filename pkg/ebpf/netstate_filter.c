#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>

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
  u16 last_state;
  u32 af;
  u16 lport;
  u16 dport;
};

#define AF_INET 2
#define AF_INET6 10

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ENTRIES);
    __type(key, struct sock *);
    __type(value, u64);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} listeners SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ENTRIES);
    __type(key, struct sock *);
    __type(value, u64);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} senders SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} pb SEC(".maps");

struct id_t {
    u32 pid;
    char task[TASK_COMM_LEN];
};

struct inet_sock_state_ctx {
    u64 __pad; // First 8 bytes are not accessible by bpf code.
    const void * skaddr;
    int oldstate;
    int newstate;
    __u16 sport;
    __u16 dport;
    __u16 family;
    __u16 protocol;
    __u8 saddr[4];
    __u8 daddr[4];
    __u8 saddr_v6[16];
    __u8 daddr_v6[16];
};

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ENTRIES);
    __type(key, struct sock *);
    __type(value, struct id_t);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} whoami SEC(".maps");

SEC("tracepoint/sock/inet_sock_set_state")
int trace_inet_sock_set_state(struct inet_sock_state_ctx *args)
{
    if (args->protocol != IPPROTO_TCP)
        return 0;

    u32 tid = bpf_get_current_pid_tgid();
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
    if (args->newstate == TCP_LISTEN) {
	struct event_t e = {};
	e.last_state = args->newstate;
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
	    BPF_CORE_READ_INTO(&e.saddr_v6, sk, __sk_common.skc_v6_rcv_saddr.in6_u.u6_addr32);
	    BPF_CORE_READ_INTO(&e.daddr_v6, sk, __sk_common.skc_v6_daddr.in6_u.u6_addr32);
	}
        bpf_perf_event_output(args, &pb, BPF_F_CURRENT_CPU, &e, sizeof(e));
    }
    return 0;
}

