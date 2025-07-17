package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	html := `
    <html>
        <head><title>ECSysLabs Go Server</title></head>
        <body>
			<h2>ECSysLabs Golang Server</h2>
            <h3>Powered by Intranet Innovations And Services LLC</h3>
            <img src="/Images/K8_KIND_Cluster_Topology.png" style="width:80%; display:block; margin:auto;" />
        </body>
    </html>`
	fmt.Fprint(w, html)
}

func main() {
	// Serve static files from the "Images" folder
	fs := http.FileServer(http.Dir("Images"))
	http.Handle("/Images/", http.StripPrefix("/Images/", fs))

	// Handle root path
	http.HandleFunc("/", handler)

	fmt.Println("Server running at http://localhost:2000")
	http.ListenAndServe(":2000", nil)
}
