// SuperClaude Cursor IDE Integration
// Enables SuperClaude commands directly in Cursor

const vscode = require('vscode');
const { spawn } = require('child_process');
const path = require('path');

let superclaudeProcess = null;
let outputChannel = null;

function activate(context) {
    console.log('SuperClaude for Cursor is now active!');
    
    // Create output channel
    outputChannel = vscode.window.createOutputChannel('SuperClaude');
    
    // Register commands
    const commands = [
        { cmd: 'superclaude.analyze', handler: () => runSuperClaudeCommand('/user:analyze', true) },
        { cmd: 'superclaude.build', handler: () => runSuperClaudeCommand('/user:build') },
        { cmd: 'superclaude.test', handler: () => runSuperClaudeCommand('/user:test') },
        { cmd: 'superclaude.improve', handler: () => runSuperClaudeCommand('/user:improve', true) },
        { cmd: 'superclaude.review', handler: () => runSuperClaudeCommand('/user:review', true) },
        { cmd: 'superclaude.custom', handler: runCustomCommand },
        { cmd: 'superclaude.start', handler: startSuperClaude },
        { cmd: 'superclaude.stop', handler: stopSuperClaude },
    ];
    
    commands.forEach(({ cmd, handler }) => {
        context.subscriptions.push(vscode.commands.registerCommand(cmd, handler));
    });
    
    // Register code actions provider
    context.subscriptions.push(
        vscode.languages.registerCodeActionsProvider(
            { scheme: 'file', pattern: '**/*' },
            new SuperClaudeCodeActionProvider(),
            { providedCodeActionKinds: [vscode.CodeActionKind.QuickFix, vscode.CodeActionKind.RefactorRewrite] }
        )
    );
    
    // Status bar item
    const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBarItem.text = '$(rocket) SuperClaude';
    statusBarItem.command = 'superclaude.custom';
    statusBarItem.tooltip = 'Run SuperClaude Command';
    statusBarItem.show();
    context.subscriptions.push(statusBarItem);
    
    // Auto-start if configured
    const config = vscode.workspace.getConfiguration('superclaude');
    if (config.get('autoStart')) {
        startSuperClaude();
    }
}

function deactivate() {
    stopSuperClaude();
}

async function runSuperClaudeCommand(command, useSelection = false) {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showErrorMessage('No active editor');
        return;
    }
    
    let text = '';
    if (useSelection && !editor.selection.isEmpty) {
        text = editor.document.getText(editor.selection);
    } else {
        text = editor.document.getText();
    }
    
    const config = vscode.workspace.getConfiguration('superclaude');
    const provider = config.get('provider', 'openrouter');
    const model = config.get('model', 'mistralai/mixtral-8x7b-instruct');
    
    outputChannel.show();
    outputChannel.appendLine(`Running: ${command}`);
    
    // Execute SuperClaude
    const superclaude = spawn('superclaude', [
        '--provider', provider,
        '--model', model,
        '--non-interactive'
    ], {
        cwd: vscode.workspace.rootPath,
        env: {
            ...process.env,
            OPENROUTER_API_KEY: config.get('apiKey')
        }
    });
    
    // Send command
    superclaude.stdin.write(`${command} ${text}\n`);
    superclaude.stdin.end();
    
    // Collect output
    let output = '';
    superclaude.stdout.on('data', (data) => {
        output += data.toString();
        outputChannel.append(data.toString());
    });
    
    superclaude.stderr.on('data', (data) => {
        outputChannel.append(`ERROR: ${data.toString()}`);
    });
    
    superclaude.on('close', (code) => {
        if (code === 0) {
            vscode.window.showInformationMessage('SuperClaude command completed');
            
            // Apply changes if it's an improve command
            if (command.includes('improve') && output.includes('```')) {
                applyCodeChanges(editor, output);
            }
        } else {
            vscode.window.showErrorMessage(`SuperClaude command failed with code ${code}`);
        }
    });
}

async function runCustomCommand() {
    const command = await vscode.window.showInputBox({
        prompt: 'Enter SuperClaude command',
        placeHolder: '/user:build --react my app',
        value: '/user:'
    });
    
    if (command) {
        runSuperClaudeCommand(command);
    }
}

function startSuperClaude() {
    if (superclaudeProcess) {
        vscode.window.showInformationMessage('SuperClaude is already running');
        return;
    }
    
    const config = vscode.workspace.getConfiguration('superclaude');
    const terminal = vscode.window.createTerminal({
        name: 'SuperClaude',
        env: {
            OPENROUTER_API_KEY: config.get('apiKey')
        }
    });
    
    terminal.sendText(`superclaude --provider ${config.get('provider')} --model ${config.get('model')}`);
    terminal.show();
}

function stopSuperClaude() {
    if (superclaudeProcess) {
        superclaudeProcess.kill();
        superclaudeProcess = null;
    }
}

function applyCodeChanges(editor, output) {
    // Extract code from output
    const codeMatch = output.match(/```[\w]*\n([\s\S]*?)```/);
    if (codeMatch) {
        const newCode = codeMatch[1];
        editor.edit(editBuilder => {
            const fullRange = new vscode.Range(
                editor.document.positionAt(0),
                editor.document.positionAt(editor.document.getText().length)
            );
            editBuilder.replace(fullRange, newCode);
        });
    }
}

class SuperClaudeCodeActionProvider {
    provideCodeActions(document, range, context) {
        const actions = [];
        
        // Quick fix for errors
        if (context.diagnostics.length > 0) {
            actions.push({
                title: '🚀 Fix with SuperClaude',
                kind: vscode.CodeActionKind.QuickFix,
                command: {
                    command: 'superclaude.improve',
                    arguments: [true]
                }
            });
        }
        
        // Refactor action
        actions.push({
            title: '🔧 Refactor with SuperClaude',
            kind: vscode.CodeActionKind.RefactorRewrite,
            command: {
                command: 'superclaude.improve',
                arguments: [true]
            }
        });
        
        // Test generation
        actions.push({
            title: '🧪 Generate Tests with SuperClaude',
            kind: vscode.CodeActionKind.Empty,
            command: {
                command: 'superclaude.test',
                arguments: []
            }
        });
        
        return actions;
    }
}

module.exports = {
    activate,
    deactivate
};