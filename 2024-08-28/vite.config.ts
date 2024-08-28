import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import customBuildPlugin from './vite-plugin';

export default defineConfig({
    plugins: [react(), customBuildPlugin()],
    build: {
        outDir: 'dist',
        rollupOptions: {
            input: {
                main: 'index.html',
            },
        },
    },
});