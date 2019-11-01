<!--- GENERATED BY gomplate from scripts/docs/templates/monitor-page.md.tmpl --->

# kubernetes-proxy

Monitor Type: `kubernetes-proxy` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/kubernetes/proxy))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

Exports Prometheus metrics from the [kube-proxy](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-proxy)
metrics in Prometheus format. The monitor queries path `/metrics` by default when no path is configured. The monitor converts
the Prometheus metric types to SignalFx metric types as described [here](prometheus-exporter.md)

Example YAML Configuration

```yaml
monitors:
- type: kubernetes-proxy
  discoveryRule: kubernetes_pod_name =~ "kube-proxy" && target == "pod"
  configEndpointMappings:
    host: '"127.0.0.1"'
  port: 10249
  extraDimensions:
    metric_source: kubernetes-proxy
```


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: kubernetes-proxy
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | **yes** | `string` | Host of the exporter |
| `port` | **yes** | `integer` | Port of the exporter |
| `username` | no | `string` | Basic Auth username to use on each request, if any. |
| `password` | no | `string` | Basic Auth password to use on each request, if any. |
| `useHTTPS` | no | `bool` | If true, the agent will connect to the exporter using HTTPS instead of plain HTTP. (**default:** `false`) |
| `skipVerify` | no | `bool` | If useHTTPS is true and this option is also true, the exporter's TLS cert will not be verified. (**default:** `false`) |
| `caCertPath` | no | `string` | Path to the CA cert that has signed the TLS cert, unnecessary if `skipVerify` is set to false. |
| `clientCertPath` | no | `string` | Path to the client TLS cert to use for TLS required connections |
| `clientKeyPath` | no | `string` | Path to the client TLS key to use for TLS required connections |
| `useServiceAccount` | no | `bool` | Use pod service account to authenticate. (**default:** `false`) |
| `metricPath` | no | `string` | Path to the metrics endpoint on the exporter server, usually `/metrics` (the default). (**default:** `/metrics`) |
| `sendAllMetrics` | no | `bool` | Send all the metrics that come out of the Prometheus exporter without any filtering.  This option has no effect when using the prometheus exporter monitor directly since there is no built-in filtering, only when embedding it in other monitors. (**default:** `false`) |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - `apiserver_audit_event_total` (*cumulative*)<br>    Counter of audit events generated and sent to the audit backend.
 - `apiserver_audit_requests_rejected_total` (*cumulative*)<br>    Counter of apiserver requests rejected due to an error in audit logging backend.
 - `go_gc_duration_seconds` (*cumulative*)<br>    A summary of the GC invocation durations.
 - `go_gc_duration_seconds_count` (*cumulative*)<br>    A summary of the GC invocation durations. (count)
 - `go_gc_duration_seconds_sum` (*cumulative*)<br>    A summary of the GC invocation durations. (sum)
 - `go_goroutines` (*gauge*)<br>    Number of goroutines that currently exist.
 - `go_info` (*gauge*)<br>    Information about the Go environment.
 - `go_memstats_alloc_bytes` (*cumulative*)<br>    Total number of bytes allocated, even if freed.
 - `go_memstats_alloc_bytes_total` (*cumulative*)<br>    Total number of bytes allocated, even if freed. (total)
 - `go_memstats_buck_hash_sys_bytes` (*gauge*)<br>    Number of bytes used by the profiling bucket hash table.
 - `go_memstats_frees_total` (*cumulative*)<br>    Total number of frees.
 - `go_memstats_gc_cpu_fraction` (*gauge*)<br>    The fraction of this program's available CPU time used by the GC since the program started.
 - `go_memstats_gc_sys_bytes` (*gauge*)<br>    Number of bytes used for garbage collection system metadata.
 - `go_memstats_heap_alloc_bytes` (*gauge*)<br>    Number of heap bytes allocated and still in use.
 - `go_memstats_heap_idle_bytes` (*gauge*)<br>    Number of heap bytes waiting to be used.
 - `go_memstats_heap_inuse_bytes` (*gauge*)<br>    Number of heap bytes that are in use.
 - `go_memstats_heap_objects` (*gauge*)<br>    Number of allocated objects.
 - `go_memstats_heap_released_bytes` (*gauge*)<br>    Number of heap bytes released to OS.
 - `go_memstats_heap_sys_bytes` (*gauge*)<br>    Number of heap bytes obtained from system.
 - `go_memstats_last_gc_time_seconds` (*gauge*)<br>    Number of seconds since 1970 of last garbage collection.
 - `go_memstats_lookups_total` (*cumulative*)<br>    Total number of pointer lookups.
 - `go_memstats_mallocs_total` (*cumulative*)<br>    Total number of mallocs.
 - `go_memstats_mcache_inuse_bytes` (*gauge*)<br>    Number of bytes in use by mcache structures.
 - `go_memstats_mcache_sys_bytes` (*gauge*)<br>    Number of bytes used for mcache structures obtained from system.
 - `go_memstats_mspan_inuse_bytes` (*gauge*)<br>    Number of bytes in use by mspan structures.
 - `go_memstats_mspan_sys_bytes` (*gauge*)<br>    Number of bytes used for mspan structures obtained from system.
 - `go_memstats_next_gc_bytes` (*gauge*)<br>    Number of heap bytes when next garbage collection will take place.
 - `go_memstats_other_sys_bytes` (*gauge*)<br>    Number of bytes used for other system allocations.
 - `go_memstats_stack_inuse_bytes` (*gauge*)<br>    Number of bytes in use by the stack allocator.
 - `go_memstats_stack_sys_bytes` (*gauge*)<br>    Number of bytes obtained from system for stack allocator.
 - `go_memstats_sys_bytes` (*gauge*)<br>    Number of bytes obtained from system.
 - `go_threads` (*gauge*)<br>    Number of OS threads created.
 - `http_request_duration_microseconds` (*cumulative*)<br>    The HTTP request latencies in microseconds.
 - `http_request_duration_microseconds_count` (*cumulative*)<br>    The HTTP request latencies in microseconds. (count)
 - `http_request_duration_microseconds_sum` (*cumulative*)<br>    The HTTP request latencies in microseconds. (sum)
 - `http_request_size_bytes` (*cumulative*)<br>    The HTTP request sizes in bytes.
 - `http_request_size_bytes_count` (*cumulative*)<br>    The HTTP request sizes in bytes. (count)
 - `http_request_size_bytes_sum` (*cumulative*)<br>    The HTTP request sizes in bytes. (sum)
 - `http_requests_total` (*cumulative*)<br>    Total number of HTTP requests made.
 - `http_response_size_bytes` (*cumulative*)<br>    The HTTP response sizes in bytes.
 - `http_response_size_bytes_count` (*cumulative*)<br>    The HTTP response sizes in bytes. (count)
 - `http_response_size_bytes_sum` (*cumulative*)<br>    The HTTP response sizes in bytes. (sum)
 - `kubeproxy_network_programming_duration_seconds_bucket` (*cumulative*)<br>    In Cluster Network Programming Latency in seconds (bucket)
 - `kubeproxy_network_programming_duration_seconds_count` (*cumulative*)<br>    In Cluster Network Programming Latency in seconds (count)
 - `kubeproxy_network_programming_duration_seconds_sum` (*cumulative*)<br>    In Cluster Network Programming Latency in seconds (sum)
 - `kubeproxy_sync_proxy_rules_duration_seconds_bucket` (*cumulative*)<br>    SyncProxyRules latency in seconds (bucket)
 - `kubeproxy_sync_proxy_rules_duration_seconds_count` (*cumulative*)<br>    SyncProxyRules latency in seconds (count)
 - `kubeproxy_sync_proxy_rules_duration_seconds_sum` (*cumulative*)<br>    SyncProxyRules latency in seconds (sum)
 - `kubeproxy_sync_proxy_rules_endpoint_changes_pending` (*gauge*)<br>    Number of pending endpoint changes that have not yet been synced to the proxy
 - `kubeproxy_sync_proxy_rules_endpoint_changes_total` (*gauge*)<br>    Number of total endpoint changes that have not yet been synced to the proxy
 - `kubeproxy_sync_proxy_rules_last_timestamp_seconds` (*gauge*)<br>
 - `kubeproxy_sync_proxy_rules_latency_microseconds_bucket` (*cumulative*)<br>    (Deprecated) SyncProxyRules latency in microseconds (bucket)
 - `kubeproxy_sync_proxy_rules_latency_microseconds_count` (*cumulative*)<br>    (Deprecated) SyncProxyRules latency in microseconds (count)
 - `kubeproxy_sync_proxy_rules_latency_microseconds_sum` (*cumulative*)<br>    (Deprecated) SyncProxyRules latency in microseconds (sum)
 - `kubeproxy_sync_proxy_rules_service_changes_pending` (*gauge*)<br>    Number of pending service changes that have not yet been synced to the proxy.
 - `kubeproxy_sync_proxy_rules_service_changes_total` (*gauge*)<br>    Number of total service changes that have not yet been synced to the proxy.
 - `kubernetes_build_info` (*gauge*)<br>    A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running.
 - `process_cpu_seconds_total` (*cumulative*)<br>    Total user and system CPU time spent, in seconds.
 - `process_max_fds` (*gauge*)<br>    Maximum number of open file descriptors.
 - `process_open_fds` (*gauge*)<br>    Number of open file descriptors.
 - `process_resident_memory_bytes` (*gauge*)<br>    Resident memory size in bytes.
 - `process_start_time_seconds` (*gauge*)<br>    Start time of the process since unix epoch in seconds.
 - `process_virtual_memory_bytes` (*gauge*)<br>    Virtual memory size in bytes.
 - `process_virtual_memory_max_bytes` (*gauge*)<br>    Maximum amount of virtual memory available in bytes.
 - `rest_client_request_duration_seconds_bucket` (*cumulative*)<br>    Request latency in seconds. Broken down by verb and URL. (bucket)
 - `rest_client_request_duration_seconds_sum` (*cumulative*)<br>    (Deprecated) Request latency in seconds. Broken down by verb and URL. (sum)
 - `rest_client_request_latency_seconds_bucket` (*cumulative*)<br>    (Deprecated) Request latency in seconds. Broken down by verb and URL. (bucket)
 - `rest_client_request_latency_seconds_count` (*cumulative*)<br>    (Deprecated) Request latency in seconds. Broken down by verb and URL. (count)
 - `rest_client_request_latency_seconds_sum` (*cumulative*)<br>    (Deprecated) Request latency in seconds. Broken down by verb and URL. (sum)
 - `rest_client_requests_total` (*cumulative*)<br>    Number of HTTP requests, partitioned by status code, method, and host.

### Enabling Additional Metrics
To emit metrics that are not _default_, you can add those metrics in the
generic monitor-level `extraMetrics` config option.  Metrics that are derived
from specific configuration options that do not appear in the above list of
metrics do not need to be added to `extraMetrics`.

To see a list of metrics that will be emitted you can run `agent-status
monitors` after configuring this monitor in a running agent instance.

