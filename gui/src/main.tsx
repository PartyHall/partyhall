import React from 'react'
import ReactDOM from 'react-dom/client'
import { ThemeProvider } from '@emotion/react'
import { createTheme, CssBaseline } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { AdapterLuxon } from '@mui/x-date-pickers/AdapterLuxon';
import i18n from "i18next";
import { initReactI18next, I18nextProvider } from "react-i18next";

import {
  createHashRouter,
  createRoutesFromElements,
  Route,
  RouterProvider
} from 'react-router-dom'

import BoothLayout from './layouts/booth_layout'
import UnauthedLayout from './layouts/unauthed_layout'
import AdminLayout from './layouts/admin_layout'
import AdminIndex from './pages/admin'
import AdminPhotobooth from './pages/admin/photobooth'
import AdminKaraoke from './pages/admin/karaoke'
import PartyHallUI from './pages/booth'
import Login from './pages/admin/login'
import SnackbarProvider from './hooks/snackbar'
import ConfirmDialogProvider from './hooks/dialog'
import EditEvent from './pages/admin/event/edit'
import ApiProvider from './hooks/useApi'
import { ENGLISH, FRENCH } from './translations'

import './assets/css/index.scss';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: ENGLISH },
      fr: { translation: FRENCH },
    },
    lng: 'en',
    fallbackLng: 'en',
    interpolation: {
      // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
      escapeValue: false,
    }
  });

const router = createHashRouter(createRoutesFromElements(
  <>
    <Route element={<UnauthedLayout />}>
      <Route path="/admin/login" element={<Login />} />
    </Route>

    <Route element={<AdminLayout />}>
      <Route path="/admin/" element={<AdminIndex />} />
      <Route path="/admin/event/edit/:id?" element={<EditEvent />} />
      <Route path="/admin/photobooth" element={<AdminPhotobooth />} />
      <Route path="/admin/karaoke" element={<AdminKaraoke />} />
    </Route>

    <Route element={<BoothLayout />}>
      <Route path="/" element={<PartyHallUI />} />
    </Route>
  </>
));

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <I18nextProvider i18n={i18n}>
        <LocalizationProvider dateAdapter={AdapterLuxon}>
          <SnackbarProvider>
            <ConfirmDialogProvider>
              <ApiProvider>
                <RouterProvider router={router} />
              </ApiProvider>
            </ConfirmDialogProvider>
          </SnackbarProvider>
        </LocalizationProvider>
      </I18nextProvider>
    </ThemeProvider>
  </React.StrictMode>
)
