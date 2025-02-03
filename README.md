# cPanel Mail Exporter

cPanel Mail Exporter is a tool that fetches email logs from a WHM/cPanel server and exposes them as Prometheus metrics.

---

## WHM API Configuration

### Enabled Features

1. Perform cPanel API and UAPI functions through the WHM API `cpanel-api`
1. Track Email `track-email`
1. Troubleshoot Mail Delivery `mailcheck`

## Usage

Run the `cpanel_mail_exporter` binary with the required command-line arguments:

```
./cpanel_mail_exporter --apikey="your-api-key" --listen="localhost:9197" --endpoint="your-whm-endpoint:2087"
```

### Command-Line Arguments

| Argument     | Description                                                                 | Example Value                  |
|--------------|-----------------------------------------------------------------------------|--------------------------------|
| `--apikey`   | WHM API key for authentication.                                             | `"your-api-key"`               |
| `--endpoint` | WHM API endpoint URL (include the port, usually `2087`).                    | `"your-whm-endpoint:2087"`     |
| `--listen`   | Address and port to expose Prometheus metrics.                              | `"localhost:9197"`             |

### Example

```
./cpanel_mail_exporter --apikey="$apikeywhm" --listen="localhost:9197" --endpoint="$url-whm:2087"
```

---

### Prometheus Configuration
```
  - job_name: 'cpanel_mail_exporter'
    scrape_interval: 60s
    metrics_path: /metrics
    static_configs:
      - targets: ['localhost:9197']
        labels:
           instance: 'serverx'
```

### Grafana
Dashboard ID:
```
22796
```
or [manually import](grafana/cpanel_mail_exporter.json)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

---


## See

- [Prometheus](https://prometheus.io/) for the metrics framework.
- [cPanel/WHM](https://api.docs.cpanel.net/) for the API.
