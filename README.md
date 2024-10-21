# Request Catcher

[![ci](https://github.com/adyanth/request-catcher/actions/workflows/docker.yaml/badge.svg)](https://github.com/adyanth/request-catcher/actions/workflows/docker.yaml)
[![GitHub license](https://img.shields.io/github/license/adyanth/request-catcher?color=brightgreen)](https://github.com/adyanth/request-catcher/blob/main/LICENSE)
[![GitHub forks](https://img.shields.io/github/forks/adyanth/request-catcher)](https://github.com/adyanth/cloudflare-operator/network)
[![GitHub stars](https://img.shields.io/github/stars/adyanth/request-catcher)](https://github.com/adyanth/cloudflare-operator/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/adyanth/request-catcher)](https://github.com/adyanth/request-catcher/issues)
[![Dockerhub](https://img.shields.io/badge/package-darkblue?style=flat&logo=docker&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fadyanth%2Frequest-catcher)](https://hub.docker.com/r/adyanth/request-catcher)
[![Github Packages](https://img.shields.io/badge/package-gray?style=flat&logo=github)
](https://github.com/adyanth/request-catcher/pkgs/container/request-catcher)

[Request Catcher](https://requestcatcher.com) is a tool for catching web requests for testing webhooks, http clients and other applications that communicate over http. Request Catcher gives you a subdomain to test your application against. Keep the index page open and instantly see all incoming requests to the subdomain via WebSockets.


## To run

Configure using environment variables or config file. Localhost with subdomains work well for local testing.

```bash
docker run --rm -p 8080:8080 adyanth/request-catcher
```

## To develop

```bash
docker build -t request-catcher .
docker run --rm -p 8080:8080 request-catcher
```

or 

```bash
go run main.go [config/sample.json]
```
