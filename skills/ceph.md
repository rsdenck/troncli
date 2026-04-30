---
name: ceph
binary: ceph
category: storage
subcategory: distributed-storage
repo: https://github.com/ceph/ceph
sdk_repo: https://github.com/ceph/go-ceph
website: https://ceph.io
language: go
license: LGPL-2.1 / LGPL-3.0 / MIT (go-ceph)

description: Plataforma distribuída de armazenamento object, block e file com alta disponibilidade e escalabilidade.

install:
  rocky_linux:
    - sudo dnf install -y epel-release
    - sudo dnf install -y ceph ceph-common
  ubuntu:
    - sudo apt update
    - sudo apt install -y ceph ceph-common
  container:
    - podman run --rm quay.io/ceph/ceph --version

sdk:
  golang:
    module: github.com/ceph/go-ceph
    install:
      - go get github.com/ceph/go-ceph
    requirements:
      - librados-devel
      - librbd-devel
      - libcephfs-devel

commands:
  verify:
    - ceph --version
  cluster_status:
    - ceph -s
  health:
    - ceph health detail
  osd:
    - ceph osd tree
  pools:
    - ceph osd pool ls
  fs:
    - ceph fs ls

nux:
  install:
    - nux skill install ceph
  info:
    - nux skill info ceph
  health:
    - nux ceph status
  pools:
    - nux ceph pools

tags:
  - storage
  - cluster
  - s3
  - cephfs
  - rbd
  - kubernetes
  - devops
---

# Ceph

## Visão Geral

Ceph é uma plataforma distribuída de armazenamento que fornece:

- Object Storage (S3/Swift via RGW)
- Block Storage (RBD)
- File Storage (CephFS)
- Replicação distribuída
- Auto-healing
- Alta disponibilidade

Muito usada em:

- Kubernetes / Rook
- Proxmox VE
- OpenStack
- Bare Metal Clusters
- Backup targets
- Large scale infra

## Repositórios Oficiais

Core:
https://github.com/ceph/ceph

Go SDK:
https://github.com/ceph/go-ceph

O `go-ceph` fornece bindings Go para:

- librados
- librbd
- cephfs
- rgw admin API

## Comandos Essenciais

```bash
ceph -s
ceph health detail
ceph osd tree
ceph osd pool ls
ceph fs ls
rbd ls
```

## Uso com Golang

```bash
go get github.com/ceph/go-ceph
```

```go
import "github.com/ceph/go-ceph/rados"
```

## Casos de Uso

- Storage para VMs
- PVCs Kubernetes
- Backup distribuído
- S3 privado
- NAS escalável
- Data lakes

## Integração com NUX

```bash
nux skill install ceph
nux ceph status
nux ceph health
nux ceph pools
nux ceph osd
nux ceph fs
```

## Observações para Rocky Linux

- MAS DEVE SER MULTI DISTRO!
- No Rocky/RHEL normalmente usar pacotes:
  - librados-devel
  - librbd-devel
  - libcephfs-devel
EOF'
cat skills/ceph.md | head -50
