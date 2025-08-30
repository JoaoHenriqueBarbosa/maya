#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🚀 Building Maya Examples"
echo "========================="

# Get wasm_exec.js
WASM_EXEC_JS="$(go env GOROOT)/lib/wasm/wasm_exec.js"
if [ ! -f "$WASM_EXEC_JS" ]; then
    WASM_EXEC_JS="$(go env GOROOT)/misc/wasm/wasm_exec.js"
fi

# Build simple example
echo -e "${YELLOW}Building simple example...${NC}"
cd simple
GOOS=js GOARCH=wasm go build -o app.wasm main.go
cp "$WASM_EXEC_JS" wasm_exec.js
cd ..

# Create HTML for simple example
cat > simple/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Maya Simple Example</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        #app {
            max-width: 600px;
            margin: 50px auto;
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        .loading {
            text-align: center;
            color: #666;
        }
    </style>
</head>
<body>
    <div id="app">
        <div class="loading">Loading Maya Framework...</div>
    </div>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
            .then((result) => {
                go.run(result.instance);
            });
    </script>
</body>
</html>
EOF

# Build advanced example
echo -e "${YELLOW}Building advanced example...${NC}"
cd advanced
GOOS=js GOARCH=wasm go build -o app.wasm main.go
cp "$WASM_EXEC_JS" wasm_exec.js
cd ..

# Create HTML for advanced example
cat > advanced/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Maya Advanced Example</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        #app {
            max-width: 800px;
            margin: 50px auto;
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        .loading {
            text-align: center;
            color: #666;
        }
    </style>
</head>
<body>
    <div id="app">
        <div class="loading">Loading Maya Framework...</div>
    </div>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
            .then((result) => {
                go.run(result.instance);
            });
    </script>
</body>
</html>
EOF

echo -e "${GREEN}✅ Build complete!${NC}"
echo ""
echo "To run examples:"
echo "  cd simple && python3 -m http.server 8080"
echo "  cd advanced && python3 -m http.server 8081"