# Helm Chart for SafeLine

## Prerequisites

- Kubernetes cluster storage support RWX.

## Installation

Install the SafeLine helm chart with a release name `safeline`:
```bash
git clone https://github.com/chaitin/SafeLine
helm -n safeline upgrade -i safeline ./SafeLine/helm --create-namespace
```

## Uninstallation

To uninstall/delete the `safeline` deployment:
```bash
helm -n safeline uninstall safeline
```
