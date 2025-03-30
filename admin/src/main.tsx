import './assets/css/index.scss';

import { ConfigProvider, theme } from 'antd';
import { I18nextProvider, initReactI18next } from 'react-i18next';
import { RouterProvider, createHashRouter } from 'react-router-dom';

import AuthProvider from './hooks/auth';
import AuthedLayout from './pages/layout/authed_layout';
import Backend from 'i18next-http-backend';
import EditEventPage from './pages/admin/edit_event';
import EventsPage from './pages/admin/events';
import Index from './pages';
import Karaoke from './pages/karaoke';
import LoginPage from './pages/login';
import LogsPage from './pages/admin/logs';
import NewEventPage from './pages/admin/new_event';
import Photobooth from './pages/photobooth';
import SettingsAudioPage from './pages/admin/settings/audio';
import SettingsPage from './pages/admin/settings';
import SettingsPhotoboothPage from './pages/admin/settings/photobooth';
import SettingsProvider from './hooks/settings';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import detector from 'i18next-browser-languagedetector';
import i18n from 'i18next';

const basePath = import.meta.env.MODE === 'development' ? '' : '/app';

i18n.use(Backend)
    .use(detector)
    .use(initReactI18next)
    .init({
        fallbackLng: 'en',
        interpolation: {
            escapeValue: false,
        },
        backend: {
            loadPath: `${basePath}/locales/{{lng}}/{{ns}}.json`,
        },
    });

const router = createHashRouter([
    {
        path: '/login',
        element: <LoginPage />,
    },
    {
        path: '/',
        element: <AuthedLayout />,
        children: [
            {
                path: '/',
                element: <Index />,
            },
            {
                path: '/logs',
                element: <LogsPage />,
            },
            {
                path: '/settings',
                children: [
                    {
                        path: '',
                        element: <SettingsPage />,
                    },
                    {
                        path: 'photobooth',
                        element: <SettingsPhotoboothPage />
                    },
                    {
                        path: 'audio',
                        element: <SettingsAudioPage />
                    },
                ]
            },
            {
                path: '/karaoke',
                element: <Karaoke />,
            },
            {
                path: '/photobooth',
                element: <Photobooth />,
            },
            {
                path: '/events',
                children: [
                    {
                        path: '',
                        element: <EventsPage />,
                    },
                    {
                        path: 'new',
                        element: <NewEventPage />,
                    },
                    {
                        path: ':id',
                        element: <EditEventPage />,
                    },
                ],
            },
        ],
    },
]);

const phTheme = {
    token: {
        colorBgBase: '#262335',
        colorTextBase: '#8a8692',
        colorError: '#db3e4b',
        colorSuccess: '#5db793',
        colorPrimary: '#f92aa9',
        colorInfo: '#f92aa9',
        sizeStep: 4,
        sizeUnit: 4,
        borderRadius: 3,
        colorBgContainer: '#241b2f',
        colorBgElevated: '#2a2139',
        fontSize: 16,
    },
    components: {
        Typography: {
            algorithm: true,
        },
        Layout: {
            headerBg: 'rgb(23,21,32)',
            headerColor: 'rgb(211,208,212)',
            siderBg: 'rgb(23,21,32)',
        },
        Menu: {
            darkItemBg: 'rgb(23,21,32)',
        },
        Modal: {
            contentBg: 'rgb(23,21,32)',
        },
        Segmented: {
            itemActiveBg: '#f92aa9',
            itemSelectedBg: '#f92aa9',
            itemSelectedColor: 'rgb(23,21,32)',
        },
    },
    algorithm: [theme.darkAlgorithm, theme.compactAlgorithm],
};

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <I18nextProvider i18n={i18n}>
            <SettingsProvider>
                {/** Settings outside of the config if-ever we want to buid server-provided themes (?) */}
                <ConfigProvider theme={phTheme}>
                    <AuthProvider>
                        <RouterProvider router={router} />
                    </AuthProvider>
                </ConfigProvider>
            </SettingsProvider>
        </I18nextProvider>
    </StrictMode>
);
