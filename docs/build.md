# Build

## How to debug locally

For Visual Studio Code:

1. Configure `kubectl` context
2. Configure debug
   1. Create folder `.vscode` in root directory
   2. Create file `launch.json` with content

        ```json
        {
            // Use IntelliSense to learn about possible attributes.
            // Hover to view descriptions of existing attributes.
            // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
            "version": "0.2.0",
            "configurations": [
                {
                    "name": "Launch Operator",
                    "type": "go",
                    "request": "launch",
                    "mode": "auto",
                    "program": "${workspaceFolder}/main.go"
                }
            ]
        }
        ```

   3. Save file
3. Run debug with using `F5` or `Debug -> Launch Operator`
