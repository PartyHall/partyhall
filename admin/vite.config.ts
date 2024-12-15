import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig(({ command }) => ({
    plugins: [react()],
    base: command === 'serve' ? '/' : '/app/',
    // Temp fix for tabler-icons
    // => https://github.com/tabler/tabler-icons/issues/1233#issuecomment-2428245119
    resolve: {
        alias: {
            // /esm/icons/index.mjs only exports the icons statically, so no separate chunks are created
            '@tabler/icons-react': '@tabler/icons-react/dist/esm/icons/index.mjs',
        },
    },
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
