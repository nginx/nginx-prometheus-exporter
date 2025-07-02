# NGINX Prometheus Exporter with Web Configuration for TLS

This example shows how to run NGINX Prometheus Exporter with web configuration. In this folder you will find an example
configuration `web-config.yml` that enables TLS and specifies the path to the TLS certificate and key files.
Additionally, there are two example TLS files `server.crt` and `server.key` that are used in the configuration.

See the full documentation of the
[web configuration](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md) for more details.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Prerequisites](#prerequisites)
- [Running NGINX Prometheus Exporter with Web Configuration in TLS mode](#running-nginx-prometheus-exporter-with-web-configuration-in-tls-mode)
- [Verification](#verification)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisites

- NGINX Prometheus Exporter binary. See the [main README](../../README.md) for installation instructions.
- NGINX or NGINX Plus running on the same machine.

## Running NGINX Prometheus Exporter with Web Configuration in TLS mode

You can run NGINX Prometheus Exporter with web configuration in TLS mode using the following command:

```console
nginx-prometheus-exporter --web.config.file=web-config.yml --nginx.scrape-uri="http://127.0.0.1:8080/stub_status"
```

you should see an output similar to this:

```console
...
ts=2023-07-20T02:00:26.932Z caller=tls_config.go:274 level=info msg="Listening on" address=[::]:9113
ts=2023-07-20T02:00:26.936Z caller=tls_config.go:310 level=info msg="TLS is enabled." http2=true address=[::]:9113
```

Depending on your environment, you may need to specify the full path to the binary or change the path to the web
configuration file.

## Verification

Run `curl -k https://localhost:9113/metrics` to see the metrics exposed by the exporter. The `-k` flag is needed because
the certificate is self-signed.
