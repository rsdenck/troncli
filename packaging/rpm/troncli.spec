Name:           troncli
Version:        0.2.21
Release:        1%{?dist}
Summary:        Production Grade Linux System Administration CLI
License:        MIT
URL:            https://github.com/rsdenck/troncli
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.22
Requires:       iproute2, nftables, lvm2

%description
TronCLI is a comprehensive tool for Linux system administration,
offering real-time monitoring, LVM management, security auditing,
and network configuration in a unified TUI/CLI interface.

%prep
%setup -q

%build
go build -ldflags="-s -w" -o troncli cmd/troncli/main.go

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
install -m 755 troncli %{buildroot}/usr/bin/troncli
mkdir -p %{buildroot}/usr/share/man/man1
install -m 644 docs/man/troncli.1 %{buildroot}/usr/share/man/man1/

%files
/usr/bin/troncli
/usr/share/man/man1/troncli.1*
%doc README.md
%license LICENSE

%changelog
* Fri Feb 20 2026 Ranlens Denck <ranlens.denck@protonmail.com> - 0.2.21-1
- Fix RPM build dependencies and changelog date
