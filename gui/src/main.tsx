import React from 'react'
import ReactDOM from 'react-dom/client'
import { ThemeProvider } from '@emotion/react'
import { createTheme, CssBaseline } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { AdapterLuxon } from '@mui/x-date-pickers/AdapterLuxon';

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
import AdminQuiz from './pages/admin/quiz'
import PartyHallUI from './pages/booth'
import Login from './pages/admin/login'
import ApiProvider from './hooks/useApi'
import SnackbarProvider from './hooks/snackbar'
import ConfirmDialogProvider from './hooks/dialog'
import EditEvent from './pages/admin/event/edit'

import './assets/css/index.scss';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

const router = createHashRouter(createRoutesFromElements(
  <>
    <Route element={<BoothLayout />}>
      <Route path="/" element={<PartyHallUI />} />
    </Route>

    <Route element={<UnauthedLayout />}>
      <Route path="/admin/login" element={<Login />} />
    </Route>

    <Route element={<AdminLayout />}>
      <Route path="/admin/" element={<AdminIndex />} />
      <Route path="/admin/event/edit/:id?" element={<EditEvent />} />
      <Route path="/admin/photobooth" element={<AdminPhotobooth />} />
      <Route path="/admin/quiz" element={<AdminQuiz />} />
    </Route>
  </>
));

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <LocalizationProvider dateAdapter={AdapterLuxon}>
        <SnackbarProvider>
          <ConfirmDialogProvider>
            <ApiProvider>
              <RouterProvider router={router} />
            </ApiProvider>
          </ConfirmDialogProvider>
        </SnackbarProvider>
      </LocalizationProvider>
    </ThemeProvider>
  </React.StrictMode>
)
