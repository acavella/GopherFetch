# GopherFetch

<!-- PROJECT SHIELDS -->
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/acavella/GopherFetch/gorelease.yml?logo=go)
![GitHub License](https://img.shields.io/github/license/acavella/GopherFetch)
![GitHub Release](https://img.shields.io/github/v/release/acavella/GopherFetch?include_prereleases)

## Overview
*The Gopher-powered Concurrent File Retrieval Tool*

GopherFetch (gfetch) is a high-performance, concurrent file retrieval utility written in [Go](https://go.dev/). Designed to act as a tireless digital "go-fer," it automates the synchronization of local file repositories with remote sources defined in a centralized configuration. By leveraging Go's native concurrency primitives, GopherFetch manages multiple downloads simultaneously while enforcing strict data integrity through cryptographic hashing. This makes it an ideal solution for maintaining local mirrors of frequently updated assets, such as Certificate Revocation Lists (CRLs), security definitions, or remote configuration files.

### Key Features

- **Concurrent Worker Pool**: Efficiently handles bulk downloads by distributing the workload across a configurable pool of Gopher workers, preventing system resource exhaustion while maximizing throughput.
- **Integrity Verification**: Uses SHA-256 checksums to compare remote assets against existing local files. If the hashes match, the program skips the download to save bandwidth and reduce disk I/O.
- **Heartbeat Synchronization**: Runs on a user-defined interval, ensuring your local directory stays in "Steady State" with remote sources without manual intervention.
- **Hot-Reloading Configuration**: Monitors its own YAML configuration file for changes. Updates to URLs, file paths, or worker counts are applied on the next sync cycle without requiring a process restart.
- **Validated CRL Retrieval**: Performs x509 cryptographic validation on all detected CRL files to verify their integrity and authenticity. This critical security check guards against replacing a legitimate local CRL with corrupted, or invalid data.

## Installation Instructions

### Native Deployment

1. Download the [latest release](https://github.com/acavella/GopherFetch/releases/latest/) archive for the appropriate platform 
   - Linux (amd64): gfetch-<version>-linux-amd64.tar.gz
   - Windows (amd64): gfetch-<version>-windows-amd64.zip
2. Extract the archive to the appropriate application directory
   - Linux: `/usr/local/bin`
   - Windows: `C:\Program Files\`
3. Edit the provided example configuration file `gfetch.yaml` and save it as `/etc/gfetch.yaml`
4. (optional) Create a system user for GopherFetch: `useradd --system --no-create-home --shell=/sbin/nologin gfetch`
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

### RPM Deployment

1. Download the [latest RPM build](https://github.com/acavella/GopherFetch/releases/latest/) appropriate for your platform:
   - EL9 (amd64): gfetch-<version>.fc43.x86_64.rpm
2. Create a system user for GopherFetch: `useradd --system --no-create-home --shell=/sbin/nologin gfetch`
3. Install the RPM with the appropriate package manager command, `sudo dnf install gfetch-<version>.fc43.x86_64.rpm`
4. Edit the sample configuration at `/etc/gfetch.sample.yaml` and rename to `/etc/gfetch.yaml`
   > [!NOTE]
   > It is important to make sure the `download_directory` is set to a directory that the `gfetch` user has permissions to read/write
5. Start and enable the gfetch systemd service:
   ```shell
   sudo systemctl start gfetch.service
   sudo systemctl enable gfetch.service
   ```

## Configuration
A list of all available configuration options is available in the sample yaml config file [gfetch.sample.yaml](gfetch.sample.yaml), with comments provided inline. Configuration is set via a static file, in which case the following paths are checked:

- `./gfetch.yaml`
- `/etc/gfetch.yaml`

## Security Vulnerabilities

I welcome and appreciate all responsible disclosures. To ensure the safety of our users, **please do not open a public Issue** to report a security vulnerability. Instead, use the GitHub private reporting system to submit your findings securely: https://github.com/acavella/GopherFetch/security/advisories/new

## Contributing

Contributions are the lifeblood of open-source projects. Help us keep GoperhFetch great by participating in the following ways:

- Propose Best Practices: Share your knowledge of RFC standards and security hardening to help us standardize the tool's behavior.
- Report Issues: Encountered a bug or an edge case in your deployment? Open an issue and help us squash it.
- Request Features: Have an idea to make GopherFetch faster or more versatile? Suggest an improvement or submit a PR.

**Important Links**:

- 🛡️ Security: Use our [Private Reporting System](https://www.google.com/search?q=https://github.com/acavella/GopherFetch/security/advisories/new) for vulnerabilities.
- 🐛 Bugs: Tracked via [GitHub Issues](https://www.google.com/search?q=https://github.com/acavella/GopherFetch/security/advisories/new).
- 📜 Rules: See our [Code of Conduct](https://www.google.com/search?q=https://github.com/acavella/GopherFetch%3Ftab%3Dcoc-ov-file%23) for community guidelines.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

- Tony Cavella - tony@cavella.com
- Project Link: [https://github.com/acavella/GopherFetch](https://github.com/acavella/GopherFetch)

## Acknowledgements
- [@Deliveranc3](https://github.com/Deliveranc3) - Containerfile development and additions to config logic

> [!NOTE]
> GopherFetch was developed using agentic coding methodologies. While the core architecture, security logic, and project direction were defined by the author, AI agents were utilized to assist with boilerplate generation, optimization, and documentation. This collaborative approach allows for faster iteration while maintaining a high standard of code integrity and RFC compliance.
