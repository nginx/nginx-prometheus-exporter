# NGINX Prometheus Exporter configuration example

This example shows how to configure NGINX Prometheus Exporter and NGINX.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
  - [NGINX](#nginx)
  - [NGINX Prometheus Exporter](#nginx-prometheus-exporter)
- [Verification](#verification)
  - [NGINX](#nginx-1)
  - [NGINX Prometheus Exporter](#nginx-prometheus-exporter-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisites

- Linux machine with [systemd](https://www.freedesktop.org/wiki/Software/systemd/).
- NGINX Prometheus Exporter binary in `/usr/local/bin/nginx-prometheus-exporter` or a location of your choice. See the
  [main README](../../README.md) for installation instructions.
- NGINX or NGINX Plus running on the same machine.

## Configuration

- NGINX Prometheus Exporter has 2 different options to get the needed metrics from the NGINX process:
    1. Via TCP/HTTP (default is via port 8080/TCP).
    2. Via a UNIX socket.
- When exposing the metrics from NGINX via TCP/HTTP you usually do not have to change anything as the
  `--nginx.scrape-uri=` setting is by default `http://127.0.0.1:8080/stub_status`.
- When exposing the metrics from NGINX via a UNIX socket you must specify the `--nginx.scrape-uri=` as
  `unix:///path/to/socket:with/location` for NGINX Prometheus Exporter.

### NGINX

You can use one of the following [NGINX server block](https://nginx.org/en/docs/http/ngx_http_core_module.html#server)
configuration examples:

**Example `stub_status` via TCP/HTTP config:**

```txt
server {
  listen [::1]:8080;
  listen 127.0.0.1:8080;
  server_name ::1 127.0.0.1 localhost;

  location / {
  # First attempt to serve request as file, then
  # as directory, then fall back to displaying a 404.
    try_files $uri $uri/ =404;
  }

  location /stub_status {
     stub_status on;
     access_log off;
     allow ::1;
     allow 127.0.0.1;
     deny all;
  }
```

**Example `stub_status` via UNIX socket config:**

```txt
server {
  listen unix:/var/run/nginx_stub_status.sock;
  server_name ::1 127.0.0.1 localhost;

  location / {
  # First attempt to serve request as file, then
  # as directory, then fall back to displaying a 404.
    try_files $uri $uri/ =404;
  }

  location /stub_status {
     stub_status on;
     access_log off;
     allow unix:;
     deny all;
  }
}
```

Depending on the Linux distribution the path you use to store the config may vary.

- When using a Debian based distribution it's usually `/etc/nginx/sites-available/`
  - Example path for config: `/etc/nginx/sites-available/nginx-prometheus-exporter.conf`
  - Enable config: `ln -s /etc/nginx/sites-available/nginx-prometheus-exporter.conf /etc/nginx/sites-enabled/`
  - Reload NGINX to make the new config take effect: `systemctl reload nginx.service`
- When using an Enterprise Linux distribution it's usually `/etc/nginx/conf.d/`
  - Example path for config: `/etc/nginx/conf.d/nginx-prometheus-exporter.conf`
  - Reload NGINX to make the new config take effect: `systemctl reload nginx.service`

> [!NOTE]
> Debian based distributions usually include `/etc/nginx/sites-enabled/*.conf` from the `/etc/nginx/nginx.conf` and
> Enterprise Linux includes `/etc/nginx/conf.d/*.conf`. Also Debian based distributions use a staging directory
> `/etc/nginx/sites-available/` where you can place your config files with server blocks and then symlink them from
> `/etc/nginx/sites-enabled/` to "enable" the specific website.

### NGINX Prometheus Exporter

If you are using the `stub_status` via TCP/HTTP via port 8080/TCP you usually do not have to change the configuration.

If you used a UNIX socket:

- On a generic Linux distribution or installed NGINX Prometheus Exporter make sure to start the service with
  the `"--nginx.scrape-uri=unix:///var/run/nginx_stub_status.sock:localhost/stub_status` argument.

- On a Debian based distribution (if you installed the NGINX Prometheus Exporter through the default repositories via
  `apt install prometheus-nginx-exporter`) you can simply adapt the config file at
  `/etc/default/prometheus-nginx-exporter` (as the systemd-service file includes it as an `EnvironmentFile`, see
  `systemctl cat prometheus-nginx-exporter.service`):

```txt
ARGS="--nginx.scrape-uri=unix:///var/run/nginx_stub_status.sock:localhost/stub_status"
```

## Verification

### NGINX

After reloading / restarting the NGINX service (usually via systemd with `systemctl reload nginx.service`) you
can test using `curl`:

```bash
$ curl localhost:8080/stub_status
Active connections: 2
server accepts handled requests
 32230 32230 62605
Reading: 0 Writing: 1 Waiting: 1
```

```bash
$ curl -X GET --unix-socket /var/run/nginx_stub_status.sock localhost/stub_status
Active connections: 2
server accepts handled requests
 32230 32230 62605
Reading: 0 Writing: 1 Waiting: 1
```

If you are having problems having a look at the systemd-journal output for the `nginx.service` is always
a good idea (`journalctl -u nginx.service`). Also make sure the `nginx.service` is actually running
(`systemctl start nginx.service`).

### NGINX Prometheus Exporter

1. Run `curl http://localhost:9113/metrics` to see the metrics exposed by the exporter.
2. Also check if `nginx_up` is `1`. If not either nginx is really not running (see `systemctl status nginx.service`)
   or something went wrong and NGINX Prometheus Exporter could not retrieve the `stub_status` from NGINX.

```bash
curl localhost:9113/metrics --silent | grep nginx_up
# HELP nginx_up Status of the last metric scrape
# TYPE nginx_up gauge
nginx_up 1
```
