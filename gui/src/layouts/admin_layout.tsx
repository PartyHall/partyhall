import { AppBar, BottomNavigation, BottomNavigationAction, Box, Drawer, IconButton, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Stack, Toolbar, Typography } from "@mui/material";
import { Link, Navigate, useOutlet } from "react-router-dom";

import MenuIcon from '@mui/icons-material/Menu';
import LogoutIcon from '@mui/icons-material/Logout';
import SettingsIcon from '@mui/icons-material/Settings';
import PhotoIcon from '@mui/icons-material/PhotoCamera';
import KaraokeIcon from '@mui/icons-material/Mic';

import { useApi } from "../hooks/useApi";
import { useState } from "react";
import { useAdminSocket } from "../hooks/adminSocket";
import { useTranslation } from "react-i18next";

type State = {
    menuOpen: boolean;
};

const linkStyle = {
    textDecoration: 'none',
    color: 'inherit',
};

export default function AdminLayout() {
    const { t } = useTranslation();
    const outlet = useOutlet();
    const { socketMode, isLoggedIn, hasRole } = useApi();
    const { appState } = useAdminSocket();

    if (socketMode != 'admin') {
        return <Navigate to={"/"} />
    }

    if (!isLoggedIn()) {
        return <Navigate to={"/admin/login"} />
    }

    return <Stack height="100%">
        <Stack maxWidth="sm" flex="1" margin="0 auto" spacing={2} padding={2} style={{ overflowY: 'scroll', width: '100%' }}>
            {outlet}
        </Stack>
        <BottomNavigation showLabels>
            <BottomNavigationAction
                component={Link}
                to="/admin/"
                label={t('admin_main.settings')}
                icon={<SettingsIcon />}
            />
            {
                appState && appState.current_mode && hasRole('ADMIN_PHOTOBOOTH') &&
                <BottomNavigationAction
                    component={Link}
                    to="/admin/photobooth"
                    label={t('admin_main.photobooth')}
                    icon={<PhotoIcon />}
                />
            }
            {
                appState && appState.current_mode &&
                <BottomNavigationAction
                    component={Link}
                    to="/admin/karaoke"
                    label={t('admin_main.karaoke')}
                    icon={<KaraokeIcon />}
                />
            }
        </BottomNavigation>
    </Stack>
}