<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/kong

 Monitors a Kong instance using [collectd-kong](https://github.com/signalfx/collectd-kong).

See the [integration documentation](https://github.com/signalfx/integrations/tree/master/collectd-kong)
for more information.

The `metrics` field below is populated with a set of metrics that are
described at https://github.com/signalfx/collectd-kong/blob/master/README.md.

Sample YAML configuration:

```yaml
monitors:
  - type: collectd/kong
    host: 127.0.0.1
    port: 8001
    metrics:
      - metric: request_latency
        report: true
      - metric: connections_accepted
        report: false
```

Sample YAML configuration with custom /signalfx route and white and blacklists

```yaml
monitors:
  - type: collectd/kong
    host: 127.0.0.1
    port: 8443
    url: https://127.0.0.1:8443/routed_signalfx
    authHeader:
      header: Authorization
      value: HeaderValue
    metrics:
      - metric: request_latency
        report: true
    reportStatusCodeGroups: true
    statusCodes:
      - 202
      - 403
      - 405
      - 419
      - "5*"
    serviceNamesBlacklist:
      - "*SomeService*"
```


Monitor Type: `collectd/kong`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/kong)

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | **yes** | `string` | Kong host to connect with (used for autodiscovery and URL) |
| `port` | **yes** | `integer` | Port for kong-plugin-signalfx hosting server (used for autodiscovery and URL) |
| `name` | no | `string` | Registration name when using multiple instances in Smart Agent |
| `url` | no | `string` | kong-plugin-signalfx metric plugin (**default:** `http://{{.Host}}:{{.Port}}/signalfx`) |
| `authHeader` | no | `object (see below)` | Header and its value to use for requests to SFx metric endpoint |
| `verifyCerts` | no | `bool` | Whether to verify certificates when using ssl/tls |
| `caBundle` | no | `string` | CA Bundle file or directory |
| `clientCert` | no | `string` | Client certificate file (with or without included key) |
| `clientCertKey` | no | `string` | Client cert key if not bundled with clientCert |
| `verbose` | no | `bool` | Whether to use debug logging for collectd-kong |
| `metrics` | no | `list of object (see below)` | List of metric names and report flags. See monitor description for more details. |
| `reportApiIds` | no | `bool` | Report metrics for distinct API IDs where applicable |
| `reportApiNames` | no | `bool` | Report metrics for distinct API names where applicable |
| `reportServiceIds` | no | `bool` | Report metrics for distinct Service IDs where applicable |
| `reportServiceNames` | no | `bool` | Report metrics for distinct Service names where applicable |
| `reportRouteIds` | no | `bool` | Report metrics for distinct Route IDs where applicable |
| `reportHttpMethods` | no | `bool` | Report metrics for distinct HTTP methods where applicable |
| `reportStatusCodeGroups` | no | `bool` | Report metrics for distinct HTTP status code groups (eg. "5xx") where applicable |
| `reportStatusCodes` | no | `bool` | Report metrics for distinct HTTP status codes where applicable (mutually exclusive with ReportStatusCodeGroups) |
| `apiIds` | no | `list of string` | List of API ID patterns to report distinct metrics for, if reportApiIds is false |
| `apiIdsBlacklist` | no | `list of string` | List of API ID patterns to not report distinct metrics for, if reportApiIds is true or apiIds are specified |
| `apiNames` | no | `list of string` | List of API name patterns to report distinct metrics for, if reportApiNames is false |
| `apiNamesBlacklist` | no | `list of string` | List of API name patterns to not report distinct metrics for, if reportApiNames is true or apiNames are specified |
| `serviceIds` | no | `list of string` | List of Service ID patterns to report distinct metrics for, if reportServiceIds is false |
| `serviceIdsBlacklist` | no | `list of string` | List of Service ID patterns to not report distinct metrics for, if reportServiceIds is true or serviceIds are specified |
| `serviceNames` | no | `list of string` | List of Service name patterns to report distinct metrics for, if reportServiceNames is false |
| `serviceNamesBlacklist` | no | `list of string` | List of Service name patterns to not report distinct metrics for, if reportServiceNames is true or serviceNames are specified |
| `routeIds` | no | `list of string` | List of Route ID patterns to report distinct metrics for, if reportRouteIds is false |
| `routeIdsBlacklist` | no | `list of string` | List of Route ID patterns to not report distinct metrics for, if reportRouteIds is true or routeIds are specified |
| `httpMethods` | no | `list of string` | List of HTTP method patterns to report distinct metrics for, if reportHttpMethods is false |
| `httpMethodsBlacklist` | no | `list of string` | List of HTTP method patterns to not report distinct metrics for, if reportHttpMethods is true or httpMethods are specified |
| `statusCodes` | no | `list of string` | List of HTTP status code patterns to report distinct metrics for, if reportStatusCodes is false |
| `statusCodesBlacklist` | no | `list of string` | List of HTTP status code patterns to not report distinct metrics for, if reportStatusCodes is true or statusCodes are specified |


The **nested** `authHeader` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `header` | **yes** | `string` | Name of header to include with GET |
| `value` | **yes** | `string` | Value of header |


The **nested** `metrics` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `metric` | **yes** | `string` | Name of metric, per collectd-kong |
| `report` | **yes** | `bool` | Whether to report this metric |




## Metrics

This monitor emits the following metrics.  Note that configuration options may
cause only a subset of metrics to be emitted.

| Name | Type | Description |
| ---  | ---  | ---         |
| `counter.kong.connections.accepted` | cumulative | Total number of all accepted connections. |
| `counter.kong.connections.handled` | cumulative | Total number of all handled connections (accounting for resource limits). |
| `counter.kong.kong.latency` | cumulative | Time spent in Kong request handling and balancer (ms). |
| `counter.kong.requests.count` | cumulative | Total number of all requests made to Kong API and proxy server. |
| `counter.kong.requests.latency` | cumulative | Time elapsed between the first bytes being read from each client request and the log writes after the last bytes were sent to the clients (ms). |
| `counter.kong.requests.size` | cumulative | Total bytes received/proxied from client requests. |
| `counter.kong.responses.count` | cumulative | Total number of responses provided to clients. |
| `counter.kong.responses.size` | cumulative | Total bytes sent/proxied to clients. |
| `counter.kong.upstream.latency` | cumulative | Time spent waiting for upstream response (ms). |
| `gauge.kong.connections.active` | gauge | The current number of active client connections (includes waiting). |
| `gauge.kong.connections.reading` | gauge | The current number of connections where nginx is reading the request header. |
| `gauge.kong.connections.waiting` | gauge | The current number of idle client connections waiting for a request. |
| `gauge.kong.connections.writing` | gauge | The current number of connections where nginx is writing the response back to the client. |
| `gauge.kong.database.reachable` | gauge | kong.dao:db.reachable() at time of metric query |


