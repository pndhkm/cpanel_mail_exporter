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
./cpanel_mail_exporter --apikey="your-api-key" --endpoint="your-whm-endpoint:2087"
```

### Command-Line Arguments

| Argument            | Description                                                                 | Example Value                  |
|---------------------|-----------------------------------------------------------------------------|--------------------------------|
| `--apikey`          | WHM API key for authentication.                                             | `"your-api-key"`               |
| `--endpoint`        | WHM API endpoint URL (include the port, usually `2087`).                    | `"your-whm-endpoint:2087"`     |
| `--listen`          | Address and port to expose Prometheus metrics.                              | `"localhost:9197"`             |
| `--web.config.file` | (Optional) Path to a YAML file for configuring TLS and Basic Auth settings. | `"path/to/web-config.yml"`     |

### Example

```
./cpanel_mail_exporter --apikey="$apikeywhm" --listen="localhost:9197" --endpoint="$url-whm:2087" --web.config.file="web-config.yml"
```

---

### Security Configuration

You can enhance the security of the exporter using TLS and Basic Authentication by providing a `web-config.yml` file.

**Example `web-config.yml`:**

```
tls_server_config:
  cert_file: server.crt
  key_file: server.key

basic_auth_users:
  alice: $2y$12$wXx/B7Gc9QbN/8VZ0fGh4.O4HvF3Z0tW5Vj3xZPL5cN7x9pkl9D4u  # bcrypt hashed password
```

- **TLS Configuration:** Provide paths to your certificate (`cert_file`) and key (`key_file`).
- **Basic Authentication:** List usernames and their bcrypt-hashed passwords.

> **Note:** Both TLS and Basic Auth are optional. If not configured, the server will run without these security features.

### Creating a bcrypt hashed password

You can create bcrypt hashed passwords using Python:

```python
import getpass
import bcrypt

password = getpass.getpass("password: ")
hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt())
print(hashed_password.decode())

```

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
