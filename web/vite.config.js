import { defineConfig } from 'vite';

export default defineConfig({
    // Set the base URL for assets in production
    // Must match the path Go serves static files from
    base: '/static/',

    build: {
        // Output directory relative to vite.config.js
        // This goes *outside* the web/ directory into the Go project's structure
        outDir: '../public/static',
        emptyOutDir: true, // Clean the directory before building

        // Generate manifest.json
        manifest: true,

        rollupOptions: {
            // Define the entry point(s)
            // The key ('main') is logical, the value is the path
            input: {
                main: './src/main.js',
            },
        },
    },

    server: {
        // Port for the Vite development server
        port: 5173,
        strictPort: true, // Error if port is already in use
    },
});