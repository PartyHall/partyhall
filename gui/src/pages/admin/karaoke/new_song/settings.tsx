import { Button } from "@mui/material";
import { useState } from "react";

import ScanIcon from '@mui/icons-material/Radar'
import { useSnackbar } from "../../../../hooks/snackbar";
import { useApi } from "../../../../hooks/useApi";

export default function SettingsKaraoke() {
    const {api} = useApi();
    const [scanning, setScanning] = useState<boolean>(false);
    const {showSnackbar} = useSnackbar();

    const rescanSongs = async (e: any) => {
        setScanning(true);
        try {
            await api.karaoke.rescanSongs();
            showSnackbar('Songs re-scanned', 'success');
        } catch (e) {
            showSnackbar('Failed to re-scan songs: ' + e, 'error');
        }
        setScanning(false);
    };

    return <>
        <Button variant="contained" component="label" color="primary" onClick={rescanSongs} disabled={scanning}>
            <ScanIcon/> Re-scan songs
        </Button>
    </>
}