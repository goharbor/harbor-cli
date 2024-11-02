## Overview

This document provides a step-by-step guide for maintainers to create and publish a release for the project using GoReleaser. The release process is automated via GitHub Actions, and includes generating a changelog, signing the release, and pushing artifacts to the specified container registry.

## Prerequisites

Before creating a release, ensure the following:
- You have **write access** to the repository.
- The required **repository secrets** and **environment variables** are set.
- You have **cosign** installed locally to generate the signing key-pair for release verification.

---

### 1. Set Up Cosign Key-Pair

Before releasing, you need to generate a cosign key-pair (in local env) to sign the release.

**Steps**:
1. Install cosign (if not installed):
   ```bash
   cosign install
   ```
2. Generate a new cosign key-pair:
   ```bash
   cosign generate-key-pair
   ```
   This will generate two files:
   - `cosign.key` (the private key)
   - `cosign.pub` (the public key)

3. Set the **private key** and **password** as GitHub repository secrets:
   - **COSIGN_KEY**: Content of `cosign.key`
   - **COSIGN_PASSWORD**: Password used to generate the key-pair

   Navigate to **Settings > Secrets and Variables > Repository secrets** and add both secrets.

---

### 2. Configure GitHub Environments

Next, create a new GitHub environment called **`production`** with the necessary secrets and variables for the release.

#### Secrets for the Production Environment

1. **REGISTRY_USERNAME**: The username for authenticating with the container registry.
2. **REGISTRY_PASSWORD**: The password for authenticating with the container registry.

**Steps**:
- Go to **Settings > Environments**.
- Click **Add environment** and name it `production`.
- Add the secrets **REGISTRY_USERNAME** and **REGISTRY_PASSWORD** under the `production` environment.

#### Environment Variables for the Production Environment

1. **REGISTRY_ADDRESS**: The address of the registry (e.g., `registry.bupd.xyz`).
2. **PUBLISH_ADDRESS**: The address to which the CLI artifacts will be published (e.g., `registry.bupd.xyz/harbor/cli`).

**Steps**:
- After adding secrets, add the following environment variables under `production`:
  - **REGISTRY_ADDRESS**
  - **PUBLISH_ADDRESS**

---

### 3. Create a GitHub Release

Once the secrets and environment are set, follow these steps to create a release:

1. Go to the **GitHub repository** and click on **Releases**.
2. Click **Draft a new release**.
3. In the **Tag version** field, specify the version number (e.g., `v0.2.0`).
4. **Do not add a description**—the changelog will be generated automatically via GitHub Actions.
5. Click **Publish Release**.

Once the release is created, the GitHub Actions workflow will:
- Generate the release changelog.
- Sign the release using `cosign` (with the `COSIGN_KEY` and `COSIGN_PASSWORD`).
- Push the CLI binaries to the container registry.

---
### 4. Verifying the Release

Once the release is completed, you can verify it by:

- Checking the GitHub Actions log for successful execution.
- Pulling the image or artifact from the registry using:
  ```bash
  # example
  docker pull registry.bupd.xyz/harbor/cli:v0.2.0
  ```

---

### 5. Troubleshooting

- **Missing GITHUB_TOKEN, GITLAB_TOKEN, or GITEA_TOKEN**:
  Ensure the required environment variables are set in GitHub secrets and accessible to the workflow.

- **Error Signing Release**:
  Double-check that the `COSIGN_KEY` and `COSIGN_PASSWORD` secrets are correctly set in GitHub.

---
