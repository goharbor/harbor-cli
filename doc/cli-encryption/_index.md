---
title: Harbor CLI Encryption
weight: 25
---

# Harbor CLI Encryption Guide

This document provides an overview of how Harbor CLI encrypts sensitive data and offers recommendations for different runtime environments. While the Harbor CLI configuration file manages your credentials, passwords themselves are never stored in plain text. Instead, they are secured using the AES-GCM encryption described below. For details on how to configure and manage these credentials, please refer to the [Harbor CLI Config Management documentation](../cli-config).

## Overview of Encryption
Harbor CLI secures credentials using AES-GCM encryption. An encryption key is automatically generated and stored in one of several keyring backends, depending on your environment:

1. **Environment-based Keyring**  
2. **System Keyring**  
3. **File-based Keyring (Fallback)**

### How It Works
1. **Key Retrieval**  
    The CLI checks if an encryption key already exists.  
    If not present, Harbor CLI generates a new 32-byte key.
2. **Encrypt**  
    Harbor CLI uses AES-GCM to encrypt credentials with a random nonce.  
    The CLI stores ciphertext as a Base64-encoded string.
3. **Decrypt**  
    Harbor CLI decodes and decrypts the stored ciphertext whenever credentials are needed.

## Recommended Backends
### Desktop Environments
**Use the system keyring** if available. This is the most secure method, as most desktop operating systems keep secrets confidentially within their native keyring service.

### Docker / Kubernetes Environments
**Use the environment-based keyring** in production. System keyrings typically aren’t available in Docker or Kubernetes. Supplying encryption keys as environment variables or secrets is straightforward.

#### Example: Running Harbor CLI in Docker
Use a randomly generated key as an environment variable:
```bash
docker run -it --rm \
  -e HARBOR_ENCRYPTION_KEY="$(openssl rand -base64 32)" \
  registry.goharbor.io/harbor-cli/harbor-cli \
  login https://demo.goharbor.io -u username -p password
```

#### Example: Running Harbor CLI in Kubernetes
Use a randomly generated key / strong password and store in a kubernetes secret. Pass the secret content as an environment variable to the deployment pod:
```yaml
# secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: harbor-cli-secrets
type: Opaque
data:
  HARBOR_ENCRYPTION_KEY: "<base64-encoded-key>"

# deployment.yaml
spec:
  containers:
  - name: harbor-cli
    image: registry.goharbor.io/harbor-cli/harbor-cli
    envFrom:
    - secretRef:
        name: harbor-cli-secrets
```

### File-based Keyring (Fallback)

In environments where neither a system keyring nor environment-based keyring is available, Harbor CLI can store the encryption key in a file. By default, this file-based keyring is located at `~/.harbor/keyring`. The stored key is protected by file permissions set to allow only the owner to read or write the file. However, **this approach is less secure** and is **not recommended for production**.
