/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import * as path from 'path';
import { stderr } from 'process';
import { workspace, ExtensionContext, languages } from 'vscode';
import { DocumentSemanticsTokensSignature, SemanticTokensFeature, SemanticTokensMiddleware } from 'vscode-languageclient/lib/common/semanticTokens';
import { execSync } from "child_process"

import {
	LanguageClient,
	LanguageClientOptions,
	ServerOptions,
	TransportKind,
	SemanticTokens,
	SemanticTokenTypes,
	SemanticTokenModifiers,
	Middleware,
	SemanticTokensParams,
	SemanticTokensRequest
} from 'vscode-languageclient/node';

let client: LanguageClient;

const legend = { tokenTypes: Object.values(SemanticTokenTypes), tokenModifiers: Object.values(SemanticTokenModifiers) };

export async function activate(context: ExtensionContext) {
	// The server is implemented in node
	/*let serverModule = context.asAbsolutePath(
		path.join('server', 'out', 'server.js')
	);*/
	// The debug options for the server
	// --inspect=6009: runs the server in Node's Inspector mode so VS Code can attach to the server for debugging
	//let debugOptions = { execArgv: ['--nolazy', '--inspect=6009'] };

	// If the extension is launched in debug mode then the debug server options are used
	// Otherwise the run options are used

	try {
		execSync("cd && go get -u github.com/clls-dev/clls/cmd/clls")
	} catch (e) {
		console.error("failed to fetch clls:", e)
	}

	const serverOptions: ServerOptions = { command: "clls" };
	// Options to control the language client
	const clientOptions: LanguageClientOptions = {
		// Register the server for plain text documents
		documentSelector: [{ scheme: 'file', language: 'chialisp' }],
		synchronize: {
			// Notify the server about file changes to '.clientrc files contained in the workspace
			fileEvents: workspace.createFileSystemWatcher('**/.clientrc')
		}
	};

	// Create the language client and start the client.
	client = new LanguageClient(
		'chialispLanguageServer',
		'Chialisp Language Server',
		serverOptions,
		clientOptions
	);


	// Start the client. This will also launch the server
	const disposable = client.start();
	context.subscriptions.push(disposable);

	languages.registerDocumentSemanticTokensProvider({ scheme: 'file', language: 'chialisp' }, {
		provideDocumentSemanticTokens: (document, token) => {
			const middleware = client.clientOptions.middleware! as Middleware & SemanticTokensMiddleware;
			const provideDocumentSemanticTokens: DocumentSemanticsTokensSignature = (document, token) => {
				const params: SemanticTokensParams = {
					textDocument: client.code2ProtocolConverter.asTextDocumentIdentifier(document)
				};
				return client.sendRequest(SemanticTokensRequest.type, params, token).then((result) => {
					return client.protocol2CodeConverter.asSemanticTokens(result);
				}, (error: any) => {
					return client.handleFailedRequest(SemanticTokensRequest.type, token, error);
				});
			};
			return middleware.provideDocumentSemanticTokens
				? middleware.provideDocumentSemanticTokens(document, token, provideDocumentSemanticTokens)
				: provideDocumentSemanticTokens(document, token);
		},
	}, legend);

	await client.onReady();
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}
