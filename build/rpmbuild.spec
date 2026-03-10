%global debug_package %{nil}

Name:           gfetch
Version:        1.0.0
Release:        1%{?dist}
Summary:        GopherFetch File Retrieval Server
License:        MIT
URL:            https://github.com/acavella/GopherFetch

Source0:        gfetch-%{version}-linux-amd64.tar.gz

BuildRequires:  systemd-rpm-macros
%{?systemd_requires}

%description
GopherFetch is a file retrieval server designed for efficient data fetching.
Note: This package assumes the 'gfetch' user already exists on the system.

%prep
%setup -q -c

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_sysconfdir}
mkdir -p %{buildroot}%{_unitdir}

# Install binary to /usr/bin
install -m 0755 gfetch %{buildroot}%{_bindir}/gfetch

# Install config file to /etc
install -m 0644 gfetch.sample.yaml %{buildroot}%{_sysconfdir}/gfetch.sample.yaml

# Create systemd service file
cat <<EOF > %{buildroot}%{_unitdir}/gfetch.service
[Unit]
Description=GopherFetch File Retrieval Server
After=network-online.target

[Service]
Type=simple
ExecStart=%{_bindir}/gfetch
User=gfetch
Restart=always

[Install]
WantedBy=multi-user.target
EOF

%post
%systemd_post gfetch.service

%preun
%systemd_preun gfetch.service

%postun
%systemd_postun_with_restart gfetch.service

%files
%{_bindir}/gfetch
%{_unitdir}/gfetch.service
# Marking as config(noreplace) to protect user edits
%config(noreplace) %{_sysconfdir}/gfetch.sample.yaml

%changelog
* Mon Mar 09 2026 Developer <mail@example.com> - 1.0.0-1
- Removed automatic user creation per user request
