// wasm_init.js - Initialize WASM with exported functions

// This will be set by the WASM module
window.handleEvent = null;
window.handleButtonClick = null;
window.onDOMReady = null;

// Initialize WASM with proper exports
async function initWASM() {
    const go = new Go();
    
    // Load WASM module
    const result = await WebAssembly.instantiateStreaming(
        fetch("app.wasm"),
        go.importObject
    );
    
    // Run the Go program
    go.run(result.instance);
    
    // The exported functions should now be available
    console.log("WASM module loaded with exports:", {
        handleEvent: typeof window.handleEvent,
        handleButtonClick: typeof window.handleButtonClick,
        onDOMReady: typeof window.onDOMReady
    });
}

// Start initialization
initWASM().catch(console.error);