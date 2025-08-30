// wasm_init_correct.js - Initialize WASM with proper exports access

async function initWASM() {
    const go = new Go();
    
    // Load WASM module
    const result = await WebAssembly.instantiateStreaming(
        fetch("app.wasm"),
        go.importObject
    );
    
    // Store the instance globally for access to exports
    window.wasmInstance = result.instance;
    window.wasmExports = result.instance.exports;
    
    // Run the Go program (this starts the Go runtime)
    go.run(result.instance);
    
    // The exported functions are now available in result.instance.exports
    console.log("WASM exports available:", Object.keys(result.instance.exports));
    
    // Call onDOMReady if DOM is ready
    if (document.readyState === "complete" || document.readyState === "interactive") {
        if (window.wasmExports && window.wasmExports.onDOMReady) {
            window.wasmExports.onDOMReady();
        }
    } else {
        document.addEventListener("DOMContentLoaded", function() {
            if (window.wasmExports && window.wasmExports.onDOMReady) {
                window.wasmExports.onDOMReady();
            }
        });
    }
}

// Start initialization
initWASM().catch(console.error);