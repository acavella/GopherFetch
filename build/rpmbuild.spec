%global debug_package %{nil}

Name:           gfetch
Version:        1.0.0
Release:        2%{?dist}
Summary:        GopherFetch File Retrieval Server
License:        MIT
URL:            https://github.com/acavella/GopherFetch

Source0:        gfetch-%{version}-linux-amd64.tar.gz
Source1:        fapolicy.rules

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
mkdir -p %{buildroot}%{_sharedstatedir}/gfetch
mkdir -p %{buildroot}%{_sysconfdir}/fapolicyd/rules.d/
cp %{SOURCE1} %{buildroot}%{_sysconfdir}/fapolicyd/rules.d/80-gfetch.rules

# Install binary to /usr/bin
install -m 0755 gfetch %{buildroot}%{_bindir}/gfetch

# Install config file to /etc
install -m 0644 gfetch.sample.yaml %{buildroot}%{_sysconfdir}/gfetch.yaml

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
if [ -x /usr/sbin/fagenrules ]; then
    /usr/sbin/fagenrules --load > /dev/null 2>&1 || :
fi

%preun
%systemd_preun gfetch.service

%postun
%systemd_postun_with_restart gfetch.service
if [ -x /usr/sbin/fagenrules ]; then
    /usr/sbin/fagenrules --load > /dev/null 2>&1 || :
fi

%files
%{_bindir}/gfetch
%{_unitdir}/gfetch.service
%config(noreplace) %{_sysconfdir}/gfetch.yaml
%dir %attr(0755, gfetch, gfetch) %{_sharedstatedir}/gfetch
%{_sysconfdir}/fapolicyd/rules.d/80-gfetch.rules


%changelog
* Mon Mar 09 2026 Tony Cavella <tony@cavella.com> - 1.0.0-1
- Initial release
* Mon Mar 11 2026 Tony Cavella <tony@cavella.com> - 1.0.0-2
- Added fapolicy rules
- Added default download directory creation
