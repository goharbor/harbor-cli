# Harbor CLI Documentation

Welcome to the Harbor CLI documentation guide! This document outlines the steps to generate the CLI documentation files.

## Overview

The [doc.go](./doc.go) script is designed to create markdown files that detail the functions of commands available in the Harbor CLI. These files are generated based on the existing cli commands and it will be placed in the [cli-docs](./cli-docs) directory.

## Steps to Generate CLI Documentation

1. Clone the repository
```bash
    git clone https://github.com/goharbor/harbor-cli.git
```

1. Navigate to the doc directory:
```bash
    cd harbor-cli/doc
```
change your directory to the `doc` folder where the `doc.go` script is located.

2. Run the `doc.go` file:
```bash
    go run doc.go
```
This generates markdown files for CLI commands that do not already exist in the [cli-docs](./cli-docs) directory.

## Generate Documentation using Dagger

Make sure you have latest [Dagger](https://docs.dagger.io/) installed in your system. 

```bash
git clone https://github.com/goharbor/harbor-cli.git
cd harbor-cli
dagger call run-doc --source=. export --path=doc
```
This would runs the dagger function and generate markdown files in [cli-docs](./cli-docs).

### Note

- For any newly generated markdown files, ensure to set the weight according to the order you want them to appear.
