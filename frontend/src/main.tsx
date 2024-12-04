import '@fontsource/poppins';
import './assets/index.scss';

import AuthProvider from './hooks/auth';
import Cookies from 'js-cookie';
import DefaultView from './components/default';
import { SnackbarProvider } from 'notistack';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

// In prod it's injected in index.html by the backend
// @ts-expect-error MERCURE_TOKEN is not defined on window
const JWT_TOKEN =
    import.meta.env.VITE_PARTYHALL_APPLIANCE_JWT ||
    window.MERCURE_TOKEN ||
    null;

Cookies.set('mercureAuthorization', JWT_TOKEN, {
    sameSite: 'lax',
    secure: true,
});

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <SnackbarProvider
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        >
            <AuthProvider token={JWT_TOKEN}>
                <DefaultView />
            </AuthProvider>
        </SnackbarProvider>
    </StrictMode>
);
