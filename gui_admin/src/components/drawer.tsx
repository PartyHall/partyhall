import { Box, Divider, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText } from '@mui/material';
import { Link } from 'react-router-dom';
import SettingsIcon from '@mui/icons-material/Settings';
import PhotoIcon from '@mui/icons-material/PhotoCamera';
import QuizIcon from '@mui/icons-material/Quiz';

type DrawerProps = {
    open: boolean;
    close: () => void;
};

export default function CustomDrawer({ open, close }: DrawerProps) {
    const linkStyle = {
        textDecoration: 'none',
        color: 'inherit',
    };

    return <Drawer anchor="left" open={open} onClose={close}>
        <Box sx={{ width: 250 }} role="presentation" onClick={close} onKeyDown={close}>
            <List>
                <Link to="/" style={linkStyle}>
                    <ListItem disablePadding>
                        <ListItemButton>
                            <ListItemIcon><SettingsIcon /></ListItemIcon>
                            <ListItemText primary="Settings" />
                        </ListItemButton>
                    </ListItem>
                </Link>
                <Link to="/photobooth" style={linkStyle}>
                    <ListItem disablePadding>
                        <ListItemButton>
                            <ListItemIcon><PhotoIcon /></ListItemIcon>
                            <ListItemText primary="Photobooth" />
                        </ListItemButton>
                    </ListItem>
                </Link>
                <Link to="/quiz" style={linkStyle}>
                    <ListItem disablePadding>
                        <ListItemButton>
                            <ListItemIcon><QuizIcon /></ListItemIcon>
                            <ListItemText primary="Quiz" />
                        </ListItemButton>
                    </ListItem>
                </Link>
            </List>
        </Box>
    </Drawer>
}