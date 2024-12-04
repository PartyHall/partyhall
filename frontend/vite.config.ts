import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig(({ command }) => ({
    plugins: [react()],
    base: command === 'serve' ? '/' : '/appliance/',
    server: {
        proxy: {
            '/api': {
                target: 'http://host.docker.internal:8080',
            },
            '/.well-known': {
                target: 'http://host.docker.internal:8080',
            },
        },
    },
}));
