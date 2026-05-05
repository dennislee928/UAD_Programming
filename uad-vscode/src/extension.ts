import * as vscode from 'vscode';
import * as path from 'path';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext) {
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

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

function registerCommands(context: vscode.ExtensionContext) {
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
        const interpreterPath = config.get<string>('interpreter.path', 'uadi');
        const filePath = document.uri.fsPath;

        // Create terminal and run
        const terminal = vscode.window.createTerminal('UAD Run');
        terminal.show();
        terminal.sendText(`${interpreterPath} -i "${filePath}"`);
    });

    // Build Project
    const buildProjectCommand = vscode.commands.registerCommand('uad.buildProject', async () => {
        const config = vscode.workspace.getConfiguration('uad');
        const compilerPath = config.get<string>('compiler.path', 'uadc');

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
        } else {
            vscode.window.showWarningMessage('Language Server is not running');
        }
    });

    context.subscriptions.push(runFileCommand, buildProjectCommand, restartLSPCommand);
}

function startLanguageServer(context: vscode.ExtensionContext) {
    const config = vscode.workspace.getConfiguration('uad.lsp');
    
    if (!config.get<boolean>('enable', true)) {
        console.log('UAD Language Server is disabled');
        return;
    }

    const serverPath = config.get<string>('serverPath', 'uad-lsp');

    // Check if server exists
    const serverOptions: ServerOptions = {
        command: serverPath,
        args: ['--stdio'],
        options: {}
    };

    const clientOptions: LanguageClientOptions = {
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
        client = new LanguageClient(
            'uadLanguageServer',
            'UAD Language Server',
            serverOptions,
            clientOptions
        );

        client.start().then(() => {
            console.log('UAD Language Server started');
        }).catch((error) => {
            console.error('Failed to start UAD Language Server:', error);
            vscode.window.showWarningMessage(
                'UAD Language Server not found. Install uad-lsp for full IDE features.'
            );
        });

        context.subscriptions.push(client);
    } catch (error) {
        console.error('Error creating Language Client:', error);
    }
}

function registerTaskProvider(context: vscode.ExtensionContext) {
    const taskProvider = vscode.tasks.registerTaskProvider('uad', {
        provideTasks: () => {
            const tasks: vscode.Task[] = [];

            // Build task
            const buildTask = new vscode.Task(
                { type: 'uad', task: 'build' },
                vscode.TaskScope.Workspace,
                'build',
                'uad',
                new vscode.ShellExecution('uadc', ['.'])
            );
            buildTask.group = vscode.TaskGroup.Build;
            buildTask.presentationOptions = {
                reveal: vscode.TaskRevealKind.Always,
                panel: vscode.TaskPanelKind.Dedicated
            };
            tasks.push(buildTask);

            // Test task
            const testTask = new vscode.Task(
                { type: 'uad', task: 'test' },
                vscode.TaskScope.Workspace,
                'test',
                'uad',
                new vscode.ShellExecution('make', ['test'])
            );
            testTask.group = vscode.TaskGroup.Test;
            tasks.push(testTask);

            // Run task
            const runTask = new vscode.Task(
                { type: 'uad', task: 'run' },
                vscode.TaskScope.Workspace,
                'run',
                'uad',
                new vscode.ShellExecution('uadi', ['-i', '${file}'])
            );
            tasks.push(runTask);

            return tasks;
        },
        resolveTask: () => {
            return undefined;
        }
    });

    context.subscriptions.push(taskProvider);
}


