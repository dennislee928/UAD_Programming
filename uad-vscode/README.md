# UAD Language Support for Visual Studio Code

Official Visual Studio Code extension for the Unified Adversarial Dynamics (UAD) language.

## Features

### ✅ Syntax Highlighting
- Full syntax highlighting for UAD keywords, operators, and literals
- Support for musical DSL (Score, Track, Motif)
- String theory semantics (String, Brane, Coupling)
- Entanglement syntax

### 🚧 IntelliSense (Coming Soon)
- Smart code completion
- Function signatures
- Type information on hover
- Parameter hints

### 🚧 Diagnostics (Coming Soon)
- Real-time error checking
- Type errors
- Syntax errors
- Linting warnings

### 🚧 Code Navigation (Coming Soon)
- Go to Definition
- Find References
- Document Symbols
- Workspace Symbols

### 🚧 Refactoring (Coming Soon)
- Rename symbol
- Extract function
- Organize imports

## Requirements

- Visual Studio Code version 1.75.0 or higher
- UAD compiler (`uadc`) and interpreter (`uadi`) installed
- UAD Language Server (`uad-lsp`) for full IDE features (optional)

## Installation

### From VSIX (Manual)

1. Download the latest `.vsix` file from [Releases](https://github.com/dennislee928/UAD_Programming/releases)
2. Open VS Code
3. Go to Extensions (Ctrl+Shift+X / Cmd+Shift+X)
4. Click the `...` menu → Install from VSIX
5. Select the downloaded `.vsix` file

### From Source

```bash
cd uad-vscode
npm install
npm run compile
npm run package
code --install-extension uad-vscode-0.1.0.vsix
```

## Extension Settings

This extension contributes the following settings:

* `uad.lsp.enable`: Enable/disable the UAD Language Server
* `uad.lsp.serverPath`: Path to the `uad-lsp` executable
* `uad.lsp.trace.server`: Trace LSP communication (off/messages/verbose)
* `uad.compiler.path`: Path to `uadc` compiler
* `uad.interpreter.path`: Path to `uadi` interpreter
* `uad.formatting.indentSize`: Indentation size (default: 4)
* `uad.diagnostics.enable`: Enable compiler diagnostics

## Commands

- `UAD: Run Current File` - Run the currently open UAD file
- `UAD: Build Project` - Build the entire UAD project
- `UAD: Restart Language Server` - Restart the UAD LSP server

## Usage

### Running UAD Programs

1. Open a `.uad` file
2. Press F5 or use Command Palette → "UAD: Run Current File"
3. View output in the integrated terminal

### Building Projects

1. Open a workspace with UAD files
2. Use Command Palette → "UAD: Build Project"
3. Check the Problems panel for errors

## Language Features Example

```uad
// Function with type annotations
fn fibonacci(n: Int) -> Int {
    if n <= 1 {
        return n;
    } else {
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}

// Struct definition
struct Point {
    x: Float,
    y: Float,
}

// Musical DSL
score GameSimulation {
    tempo: 120,
    track Player {
        bars 1..4 {
            motif attack;
        }
    }
}

// String Theory Semantics
string EthicalField {
    modes {
        integrity: Float,
        transparency: Float,
    }
}

// Entanglement
let x: Int = 10;
let y: Int = 20;
entangle x, y;  // x and y now share the same value
```

## Keyboard Shortcuts

| Command | Windows/Linux | macOS |
|---------|---------------|-------|
| Run File | F5 | F5 |
| Build Project | Ctrl+Shift+B | Cmd+Shift+B |
| Toggle Terminal | Ctrl+` | Cmd+` |

## Known Issues

- Language Server features are not yet fully implemented
- Debugger support is in development
- Some advanced refactoring features are pending

See the [issue tracker](https://github.com/dennislee928/UAD_Programming/issues) for a complete list.

## Release Notes

### 0.1.0 (Initial Release)

- Basic syntax highlighting
- File association for `.uad` and `.uadmodel`
- Code snippets
- Task runner integration
- LSP client (server implementation pending)

## Contributing

Contributions are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming/uad-vscode

# Install dependencies
npm install

# Compile TypeScript
npm run compile

# Watch for changes
npm run watch

# Test the extension
# Press F5 in VS Code to launch Extension Development Host
```

## Resources

- [UAD Language Documentation](../docs/)
- [Language Specification](../docs/specs/CORE_LANGUAGE_SPEC.md)
- [GitHub Repository](https://github.com/dennislee928/UAD_Programming)
- [Issue Tracker](https://github.com/dennislee928/UAD_Programming/issues)

## License

Apache License 2.0 - See [LICENSE](../LICENSE) for details.

---

**Enjoy coding in UAD!** 🎉
