"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = __importStar(require("vscode"));
const node_1 = require("vscode-languageclient/node");
let client;
function activate(context) {
    console.log('UAD Language Extension activated');
    // Register commands
    registerCommands(context);
    // Start LSP client (if server is available)
    startLanguageServer(context);
    // Register task provider
    registerTaskProvider(context);
    // Status bar item
    const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBarItem.text = '$(check) UAD';
    statusBarItem.tooltip = 'UAD Language Extension Active';
    statusBarItem.show();
    context.subscriptions.push(statusBarItem);
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
function registerCommands(context) {
    // Run Current File
    const runFileCommand = vscode.commands.registerCommand('uad.runFile', async () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) {
            vscode.window.showErrorMessage('No active editor');
            return;
        }
        const document = editor.document;
        if (document.languageId !== 'uad') {
            vscode.window.showErrorMessage('Current file is not a UAD file');
            return;
        }
        // Save file before running
        await document.save();
        const config = vscode.workspace.getConfiguration('uad');
        const interpreterPath = config.get('interpreter.path', 'uadi');
        const filePath = document.uri.fsPath;
        // Create terminal and run
        const terminal = vscode.window.createTerminal('UAD Run');
        terminal.show();
        terminal.sendText(`${interpreterPath} -i "${filePath}"`);
    });
    // Build Project
    const buildProjectCommand = vscode.commands.registerCommand('uad.buildProject', async () => {
        const config = vscode.workspace.getConfiguration('uad');
        const compilerPath = config.get('compiler.path', 'uadc');
        const terminal = vscode.window.createTerminal('UAD Build');
        terminal.show();
        terminal.sendText(`${compilerPath} .`);
    });
    // Restart LSP
    const restartLSPCommand = vscode.commands.registerCommand('uad.restartLSP', async () => {
        if (client) {
            await client.stop();
            await client.start();
            vscode.window.showInformationMessage('UAD Language Server restarted');
        }
        else {
            vscode.window.showWarningMessage('Language Server is not running');
        }
    });
    context.subscriptions.push(runFileCommand, buildProjectCommand, restartLSPCommand);
}
function startLanguageServer(context) {
    const config = vscode.workspace.getConfiguration('uad.lsp');
    if (!config.get('enable', true)) {
        console.log('UAD Language Server is disabled');
        return;
    }
    const serverPath = config.get('serverPath', 'uad-lsp');
    // Check if server exists
    const serverOptions = {
        command: serverPath,
        args: ['--stdio'],
        options: {}
    };
    const clientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'uad' },
            { scheme: 'file', language: 'uad', pattern: '**/*.uad' },
            { scheme: 'file', language: 'uad', pattern: '**/*.uadmodel' }
        ],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.uad')
        }
    };
    try {
        client = new node_1.LanguageClient('uadLanguageServer', 'UAD Language Server', serverOptions, clientOptions);
        client.start().then(() => {
            console.log('UAD Language Server started');
        }).catch((error) => {
            console.error('Failed to start UAD Language Server:', error);
            vscode.window.showWarningMessage('UAD Language Server not found. Install uad-lsp for full IDE features.');
        });
        context.subscriptions.push(client);
    }
    catch (error) {
        console.error('Error creating Language Client:', error);
    }
}
function registerTaskProvider(context) {
    const taskProvider = vscode.tasks.registerTaskProvider('uad', {
        provideTasks: () => {
            const tasks = [];
            // Build task
            const buildTask = new vscode.Task({ type: 'uad', task: 'build' }, vscode.TaskScope.Workspace, 'build', 'uad', new vscode.ShellExecution('uadc', ['.']));
            buildTask.group = vscode.TaskGroup.Build;
            buildTask.presentationOptions = {
                reveal: vscode.TaskRevealKind.Always,
                panel: vscode.TaskPanelKind.Dedicated
            };
            tasks.push(buildTask);
            // Test task
            const testTask = new vscode.Task({ type: 'uad', task: 'test' }, vscode.TaskScope.Workspace, 'test', 'uad', new vscode.ShellExecution('make', ['test']));
            testTask.group = vscode.TaskGroup.Test;
            tasks.push(testTask);
            // Run task
            const runTask = new vscode.Task({ type: 'uad', task: 'run' }, vscode.TaskScope.Workspace, 'run', 'uad', new vscode.ShellExecution('uadi', ['-i', '${file}']));
            tasks.push(runTask);
            return tasks;
        },
        resolveTask: () => {
            return undefined;
        }
    });
    context.subscriptions.push(taskProvider);
}
//# sourceMappingURL=extension.js.map