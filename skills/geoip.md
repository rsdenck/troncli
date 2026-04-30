# skills/geoip.md

## Skill Name
geoip

## Category
networking, security, observability

## Status
active

## Priority
high

## Objective

Desenvolver completamente a nova skill **geoip** para a CLI **NUX**, seguindo o padrão visual, estrutural e arquitetural do projeto.

A skill deve permitir geolocalização de IPs, análise de logs, investigação de conexões ativas, enriquecimento de eventos de segurança e integração com firewall.

Essa skill deve parecer nativa do NUX.

---

# Requisitos obrigatórios

- Usar o mesmo padrão visual dos demais comandos NUX
- Outputs limpos, profissionais e consistentes
- Seguir arquitetura Go atual do projeto
- Implementação modular
- Código pronto para produção
- Validado e testado
- Compatível com Rocky Linux, RHEL, Ubuntu, Debian
- Performance alta
- Sem dependência desnecessária
- Tratamento robusto de erros
- UX estilo enterprise CLI

---

# Stack Técnica Obrigatória

## Banco GeoIP

Usar:

- GeoLite2-City.mmdb
- GeoLite2-ASN.mmdb
- GeoLite2-Country.mmdb

## Biblioteca Go

Usar:

```go
github.com/oschwald/geoip2-golang/v2
```

# Campos:

geoip:
  enabled: true
  database_path: /opt/nux/geoip/GeoLite2-City.mmdb
  country_db: /opt/nux/geoip/GeoLite2-Country.mmdb
  asn_db: /opt/nux/geoip/GeoLite2-ASN.mmdb
  cache: true
  cache_ttl: 24h

---------------------------------------------------------------------------------
# NOVA SKILL A SER DESENVOLVIDA, VALIDADA E TESTADA! 
- USAR O MESMO PADRÃO DE OUTPUT DOS DEMAIS COMANDOS!
- MANTER O PADRÃO DA CLI NUX!
- skills/geoip.md
# Usando:
MaxMind GeoLite2
oschwald/geoip2-golang
banco .mmdb


# Casos reais no NUX:
nux geoip 8.8.8.8
nux geoip nginx.log
nux geoip sshd.log
nux geoip top attackers
nux geoip connections
nux geoip whois 1.1.1.1

# Segurança
- Ver origem de ataques SSH:
nux logs ssh --geo

# Firewall inteligente no NUX
nux firewall block-country RU
nux firewall allow-country BR


# Network Analytics
nux net connections --geo


# Apache / Nginx Logs
nux logs nginx --geo --top

# Fail2ban Turbo
nux security attackers --map


# stack técnica
Banco

Use:

GeoLite2-City.mmdb
GeoLite2-ASN.mmdb
GeoLite2-Country.mmdb
Go lib

Use:

github.com/oschwald/geoip2-golang/v2
