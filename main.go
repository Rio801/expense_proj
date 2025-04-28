package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// --- Configuration ---
var (
	AppEnv        = os.Getenv("APP_ENV")
	IsDevelopment = AppEnv == "development"
	// URL of the Vite Dev server (running separately)
	ViteDevServer = "http://localhost:5173"
	// Directory where Vite builds assets (relative to main.go)
	DistDir = "public/static"
	// Path to the manifest file (relative to main.go)
	ManifestFile = filepath.Join(DistDir, ".vite", "manifest.json")
)

// --- Manifest Handling ---
type ViteManifest map[string]struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
	// Src  string `json:"src"` // Can add if needed
	// IsEntry bool `json:"isEntry"` // Can add if needed
}

var assetManifest ViteManifest

func loadManifest() {
	if IsDevelopment {
		log.Println("Development mode: Skipping manifest load.")
		return // Don't load manifest in development
	}

	log.Println("Production mode: Attempting to load manifest...")
	content, err := os.ReadFile(ManifestFile)
	if err != nil {
		// Log fatal because we NEED the manifest in production to serve assets
		log.Fatalf("FATAL: Failed to read manifest file '%s': %v", ManifestFile, err)
	}
	err = json.Unmarshal(content, &assetManifest)
	if err != nil {
		log.Fatalf("FATAL: Failed to parse manifest file '%s': %v", ManifestFile, err)
	}
	log.Println("Successfully loaded production asset manifest.")
}

// --- Template Function ---
func viteAssets(entrypoints ...string) template.HTML {
	var scripts strings.Builder
	var styles strings.Builder

	if IsDevelopment {
		// Development: Point to Vite dev server
		scripts.WriteString(`<script type="module" src="` + ViteDevServer + `/@vite/client"></script>`)
		for _, entry := range entrypoints {
			// IMPORTANT: In dev, Vite expects paths relative to project root (web/src)
			// We use the same key as used in the template call (e.g., "src/main.js")
			scripts.WriteString(`<script type="module" src="` + ViteDevServer + `/` + entry + `"></script>`)
		}
	} else {
		// Production: Use manifest.json
		for _, entry := range entrypoints {
			// Use the entry key (e.g., "src/main.js") passed from the template
			manifestEntry, ok := assetManifest[entry]
			if !ok {
				// Log a warning if an entry point isn't found in the manifest
				log.Printf("WARN: Entrypoint '%s' not found in production manifest '%s'", entry, ManifestFile)
				continue
			}

			// Add the main JS file script tag (using the path from manifest)
			// Prepends /static/ which is the base path Go serves from
			scripts.WriteString(`<script type="module" src="/static/` + manifestEntry.File + `"></script>`)

			// Add linked CSS files (using paths from manifest)
			for _, cssFile := range manifestEntry.CSS {
				styles.WriteString(`<link rel="stylesheet" href="/static/` + cssFile + `">`)
			}
		}
	}
	return template.HTML(styles.String() + scripts.String())
}

func main() {
	log.Printf("Starting app: APP_ENV='%s', IsDevelopment=%t", AppEnv, IsDevelopment)
	loadManifest() // Load manifest ONLY if in production

	// --- Serve Static Assets (Production Only) ---
	if !IsDevelopment {
		fs := http.FileServer(http.Dir(DistDir))
		// Serve files from /static/ URL path, stripping the prefix
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		log.Printf("Serving production static assets from /static/ mapped to ./%s", DistDir)
	}

	// --- Template Setup ---
	templatesRoot := "templates"
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"viteAssets": viteAssets,
	}).ParseFiles(filepath.Join(templatesRoot, "index.html"))) // Parse the specific template

	// --- HTTP Handler ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Go + Vite Reset Example",
		}
		// Execute the named template "index.html"
		err := tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// --- Start Server ---
	port := ":8080"
	log.Printf("Go server listening on http://localhost%s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
