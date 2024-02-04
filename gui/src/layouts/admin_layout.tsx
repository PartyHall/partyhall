import { AppBar, Box, Drawer, IconButton, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Stack, Toolbar, Typography } from "@mui/material";
import { Link, Navigate, useOutlet } from "react-router-dom";

import MenuIcon from '@mui/icons-material/Menu';
import LogoutIcon from '@mui/icons-material/Logout';
import SettingsIcon from '@mui/icons-material/Settings';
import PhotoIcon from '@mui/icons-material/PhotoCamera';
import KaraokeIcon from '@mui/icons-material/Mic';

import { useApi } from "../hooks/useApi";
import { useState } from "react";
import { useAdminSocket } from "../hooks/adminSocket";

type State = {
    menuOpen: boolean;
};

const linkStyle = {
    textDecoration: 'none',
    color: 'inherit',
};

export default function AdminLayout() {
    const outlet = useOutlet();
    const { socketMode, token, logout } = useApi();
    const {appState} = useAdminSocket();
    const [state, setState] = useState<State>({
        menuOpen: false,
    });

    if (socketMode != 'admin') {
        return <Navigate to={"/"} />
    }

    if (!token) {
        return <Navigate to={"/admin/login"} />
    }

    const close = () => setState({...state, menuOpen: false});

    return <>
        <AppBar position="static">
            <Toolbar>
                <IconButton size="large" edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }} onClick={() => setState({ ...state, menuOpen: true })}><MenuIcon /></IconButton>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>Mode: {appState?.current_mode}</Typography>
                <IconButton size="large" edge="start" color="inherit" aria-label="menu" onClick={logout}><LogoutIcon /></IconButton>
            </Toolbar>
        </AppBar>

        <Drawer anchor="left" open={state.menuOpen} onClose={close}>
            <Box sx={{ width: 250 }} role="presentation" onClick={close} onKeyDown={close}>
                <List>
                    <Link to="/admin/" style={linkStyle}>
                        <ListItem disablePadding>
                            <ListItemButton>
                                <ListItemIcon><SettingsIcon /></ListItemIcon>
                                <ListItemText primary="Settings" />
                            </ListItemButton>
                        </ListItem>
                    </Link>
                    <Link to="/admin/photobooth" style={linkStyle}>
                        <ListItem disablePadding>
                            <ListItemButton>
                                <ListItemIcon><PhotoIcon /></ListItemIcon>
                                <ListItemText primary="Photobooth" />
                            </ListItemButton>
                        </ListItem>
                    </Link>
                    <Link to="/admin/karaoke" style={linkStyle}>
                        <ListItem disablePadding>
                            <ListItemButton>
                                <ListItemIcon><KaraokeIcon /></ListItemIcon>
                                <ListItemText primary="Karaoke" />
                            </ListItemButton>
                        </ListItem>
                    </Link>
                </List>
            </Box>
        </Drawer>

        <div style={{height: '100%', paddingBottom: '5em'}}>
            <Stack maxWidth="sm" spacing={2} margin="auto" paddingTop={2} style={{height: "100%"}}>
                {outlet}
            </Stack>
        </div>
    </>
}