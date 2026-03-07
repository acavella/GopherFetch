# GopherFetch

<!-- PROJECT SHIELDS -->
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/acavella/gopherfetch/gorelease.yml?logo=go)
![GitHub License](https://img.shields.io/github/license/acavella/gopherfetch)
![GitHub release (with filter)](https://img.shields.io/github/v/release/acavella/gopherfetch)

## Overview
*The Gopher-powered Concurrent File Retrieval Tool*

GopherFetch (gfetch) is a standalone file sync service used to retrieve files from remote http(s) destinations to a local directory. GopherFetch is written in [Go](https://go.dev/), designed to be lightweight and fully self-contained using a simple configuration. 

### Key Features

- Cross-platform compatiblity; tested on Linux and Windows
- Native and containerized deployment options
- Retrieve remote CRL data via HTTP or HTTPS
- Validation and confirmation of CRL data
- Built-in webserver alleviates the need for additional servers
- Ability to retrieve and serve an unlimited number of CRL sources
- Support for full and delta CRLs

### Planned Features

- OCSP responder

## Installation Instructions

### Native Deployment

1. Download the [latest release](https://github.com/acavella/GopherFetch/releases/latest/) archive for the appropriate platform 
   - Linux (amd64): gfetch-<version>-linux-amd64.tar.gz
   - Windows (amd64): gfetch-<version>-windows-amd64.zip
2. Extract the archive to the appropriate application directory
   - Linux: `/usr/local/bin`
   - Windows: `C:\Program Files\`
3. Edit the provided example configuration file `gfetch.yaml` and save it as `/etc/gfetch.yaml`
4. (optional) Create a system user for GoRevoke: `useradd --system --no-create-home --shell=/sbin/nologin gfetch`
5. Create a systemd service file `/etc/systemd/service/gfetch.service`. Example unit files:
```ini
### Using a static-file configuration
[Unit]
Description=GopherFetch File Retrieval Server
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/gfetch
User=gfetch
Restart=always

[Install]
WantedBy=multi-user.target default.target
```
6. Set the permissions `sudo chmod 664 /etc/systemd/service/gfetch.service`
7. Reload the systemd configuration `sudo systemctl daemon-reload`
8. Enable and start the service:
```shell
sudo systemctl enable --now gfetch.service
```

## Configuration
A list of all available configuration options is available at [gorevoke.yml](gorevoke.yml), with comments provided inline. Configuration can be set via a static file, in which case the following paths are checked:

- `$PWD/gorevoke.yml`
- `$HOME/.gorevoke/gorevoke.yml`
- `/etc/gorevoke.yml`

Optionally, all configuration values can be specified via environment variables, upper-cased and prefixed with `GOREVOKE`. For example, the configuration item `default.interval` can be set via the `GOREVOKE_DEFAULT_INTERVAL` variable. If specifying the list of CRLs as an environment var (`GOREVOKE_CRLS`), the CRLs must be provided as a json dict. See the systemd unit example, above.

## Container Performance
![Docker Container Performance](assets/docker-stats.png)

## Security Vulnerabilities

I welcome welcome all responsible disclosures. Please do not open an ISSUE to report a security problem. Please use the private reporting system to report security related issues responsibly: https://github.com/acavella/gorevoke/security/advisories/new

## Contributing

Contributions are essential to the success of open-source projects. In other words, we need your help to keep GoRevoke great!

What is a contribution? All the following are highly valuable:

1. **Let us know of the best-practices you believe should be standardized**   
   GoRevoke is designed to be compliant with applicable RFCs out-of-the box. By sharing your experiences and knowledge you help us build a solution that takes into account best-practices and user experience.

2. **Let us know if things aren't working right**   
   We aim to provide a perfect application and test it extensively, however, we can't imagine or replicate every deployment scenario possible. If you run into an issue that you think isn't normal, please let us know.

3. **Add or improve features**   
   Have an idea to add or improve functionality, then let us know! We want to make GoRevoke the best total solution it can be.

**General information about contributions:**

Check our [Security Policy](https://github.com/acavella/gorevoke#).   
Found a bug? Open a [GitHub issue](https://github.com/acavella/gorevoke/issues).   
Read our [Contributing Code of Conduct](https://github.com/acavella/gorevoke?tab=coc-ov-file#), which contains all the information you need to contribute to GoRevoke!

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

- Tony Cavella - tony@cavella.com
- Project Link: [https://github.com/acavella/gorevoke](https://github.com/acavella/gorevoke)

## Acknowledgements
- [@Deliveranc3](https://github.com/Deliveranc3) - Containerfile development and additions to config logic
