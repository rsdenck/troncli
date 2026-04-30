Name:           nux
Version:        1.0.0
Release:        1%{?dist}
Summary:        NUX - Linux CLI Manager for all distributions
License:        MIT
URL:            https://github.com/rsdenck/nux
Source0:        https://github.com/rsdenck/nux/releases/download/v%{version}/nux_%{version}_linux_amd64.tar.gz

BuildRequires:  systemd
Requires:       bash

%description
NUX is a production-grade CLI for comprehensive Linux system administration.
It supports multiple distributions (apt, dnf, yum, pacman, apk, zypper)
and provides universal management for services, packages, network, disk, and more.

Features:
- Universal Package Management
- Service Management (systemd, openrc, sysvinit)
- Network Configuration & Diagnostics
- Skill Engine for managing external CLIs
- Agent Integration (Ollama AI)

%prep
%setup -q -n nux_%{version}_linux_amd64

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 nux %{buildroot}%{_bindir}/nux

mkdir -p %{buildroot}%{_datadir}/doc/nux
install -m 644 README.md %{buildroot}%{_datadir}/doc/nux/
install -m 644 LICENSE %{buildroot}%{_datadir}/doc/nux/

%files
%{_bindir}/nux
%{_datadir}/doc/nux/*

%changelog
* Thu Apr 30 2026 Ranlens Denck <piptr@protonmail.com> - 1.0.0-1
- Initial package for NUX CLI Manager
- Support for multiple Linux distributions
- Skill engine and Ollama integration
