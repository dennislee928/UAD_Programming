# Landing Page 設置指南

本文檔說明如何設置 UAD 語言的 GitHub Pages Landing Page。

## 工具選擇

推薦使用 **MkDocs**，因為：
- 適合技術文檔
- 自動生成導航
- 支援 Markdown
- 易於部署到 GitHub Pages
- 主題豐富（Material 主題推薦）

## 設置步驟

### 1. 安裝 MkDocs

```bash
pip install mkdocs mkdocs-material
```

### 2. 初始化 MkDocs 專案

```bash
cd docs
mkdocs new .
```

### 3. 配置 `mkdocs.yml`

創建 `docs/mkdocs.yml`:

```yaml
site_name: UAD Programming Language
site_description: Domain-Specific Language for Adversarial Modeling, Ethical Risk, and Cognitive Security
site_author: UAD Team
site_url: https://dennislee928.github.io/UAD_Programming/

theme:
  name: material
  palette:
    - scheme: default
      primary: blue
      accent: blue
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - scheme: slate
      primary: blue
      accent: blue
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - navigation.tabs
    - navigation.sections
    - navigation.expand
    - navigation.top
    - search.suggest
    - search.highlight

repo_name: dennislee928/UAD_Programming
repo_url: https://github.com/dennislee928/UAD_Programming
edit_uri: edit/main/docs/

nav:
  - Home: index.md
  - Quick Start: LANGUAGE_GUIDE.md#快速開始
  - Language Guide: LANGUAGE_GUIDE.md
  - Examples: ../examples/
  - Specification: 
    - Core Language: specs/CORE_LANGUAGE_SPEC.md
    - Model DSL: specs/MODEL_DSL_SPEC.md
    - IR: specs/IR_SPEC.md
  - Paradigm: PARADIGM.md
  - Whitepaper: WHITEPAPER.md
  - Contributing: ../CONTRIBUTING.md

markdown_extensions:
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite
  - pymdownx.snippets
  - pymdownx.superfences
  - admonition
  - pymdownx.details
  - pymdownx.superfences:
      custom_fences:
        - name: uad
          class: uad
          format: !!python/name:pymdownx.superfences.fence_code_format
```

### 4. 創建 Landing Page (`docs/index.md`)

```markdown
# UAD Programming Language

<div align="center">

**Unified Adversarial Dynamics Language**  
_Domain-Specific Language for Adversarial Modeling, Ethical Risk, and Cognitive Security_

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

[Get Started](#quick-start) • [Documentation](LANGUAGE_GUIDE.md) • [Examples](../examples/) • [GitHub](https://github.com/dennislee928/UAD_Programming)

</div>

---

## Why UAD?

UAD is designed from the ground up for modeling adversarial dynamics, ethical risk, and cognitive security systems.

Unlike general-purpose languages, UAD treats **decisions, risks, and time** as first-class citizens, providing native semantics to describe complex adversarial behaviors, quantify ethical consequences, and simulate long-term system evolution.

### Key Features

- ✅ **Strong Type System**: Static type checking ensures model correctness
- ✅ **Temporal Semantics**: Musical DSL (Score/Track/Bar/Motif) for precise time control
- ✅ **Adversarial Modeling**: Built-in Action/Judge/Agent concepts
- ✅ **String Theory Semantics**: Describe complex coupling and resonance relationships
- ✅ **Quantum Entanglement**: Variable synchronization mechanism
- ✅ **Dual Execution**: Interpreter + Virtual Machine

---

## Quick Start

### Installation

```bash
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming
make build
```

### First Program

```uad
fn main() {
    println("Hello, UAD!");
}
```

Run it:

```bash
./bin/uadi hello.uad
```

---

## Typical Use Cases

### 1. Cybersecurity Threat Modeling

Model complete attack kill-chains with temporal precision:

```uad
score RansomwareAttack {
    track attacker {
        bars 1..4 { use reconnaissance; }
        bars 5..8 { use initial_access; }
        // ...
    }
}
```

### 2. Ethical Risk Analysis

Quantify ethical consequences using ERH (Ethical Riemann Hypothesis):

```uad
string EthicalField {
    modes {
        integrity: Float,
        transparency: Float,
    }
}

coupling EthicalField.integrity EthicalField.transparency 
    with strength 0.7
```

### 3. System Evolution Simulation

Model large-scale system dynamics:

```uad
resonance when system_stability > 8.0 {
    emit Event {
        type: "positive_feedback",
        effect: "accelerated_growth",
    };
}
```

---

## Resources

- 📘 [Language Guide](LANGUAGE_GUIDE.md) - Complete language reference
- 📋 [Specification](specs/CORE_LANGUAGE_SPEC.md) - Formal language specification
- 🎯 [Examples](../examples/) - Example programs
- 🔌 [VS Code Extension](../uad-vscode/) - IDE support

---

## Get Involved

- 🐛 [Report Issues](https://github.com/dennislee928/UAD_Programming/issues)
- 💡 [Suggest Features](https://github.com/dennislee928/UAD_Programming/discussions)
- 📝 [Contribute](CONTRIBUTING.md)

---

<div align="center">

**Made with ❤️ by the UAD Community**

</div>
```

### 5. 部署到 GitHub Pages

創建 `.github/workflows/docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/docs.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v4
        with:
          python-version: '3.x'
      - run: pip install mkdocs mkdocs-material
      - run: mkdocs gh-deploy --force
```

## 替代方案：Docsify

如果偏好更簡單的方案，可以使用 Docsify：

### 設置 Docsify

```bash
npm install -g docsify-cli
docsify init docs
```

### 配置 `docs/index.html`

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>UAD Programming Language</title>
  <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
  <meta name="description" content="Unified Adversarial Dynamics Language">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0">
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/docsify@4/lib/themes/vue.css">
</head>
<body>
  <div id="app"></div>
  <script>
    window.$docsify = {
      name: 'UAD Programming Language',
      repo: 'dennislee928/UAD_Programming',
      loadSidebar: true,
      subMaxLevel: 2,
      homepage: 'README.md'
    }
  </script>
  <script src="//cdn.jsdelivr.net/npm/docsify@4"></script>
</body>
</html>
```

## 推薦內容結構

```
docs/
├── index.md              # Landing page
├── LANGUAGE_GUIDE.md     # Language guide
├── PARADIGM.md          # Language paradigm
├── specs/               # Specifications
│   ├── CORE_LANGUAGE_SPEC.md
│   ├── MODEL_DSL_SPEC.md
│   └── IR_SPEC.md
├── mkdocs.yml           # MkDocs config (if using MkDocs)
└── _sidebar.md          # Sidebar (if using Docsify)
```

## 當前狀態

- ✅ 設置指南已建立
- ⏳ 等待實際設置和部署
- ⏳ 等待 Landing Page 內容創建

---

**最後更新**: 2025-01-07


