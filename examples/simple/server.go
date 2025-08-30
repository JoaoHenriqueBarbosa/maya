package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve static files from current directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Printf("Server running at http://localhost:%s\n", port)
	fmt.Println("Available pages:")
	fmt.Printf("  http://localhost:%s/index_dom.html     - DOM Renderer\n", port)
	fmt.Printf("  http://localhost:%s/index_canvas.html  - Canvas Renderer\n", port)
	fmt.Printf("  http://localhost:%s/compare.html       - Side by side comparison\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}