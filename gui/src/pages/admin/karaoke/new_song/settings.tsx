import { Button } from "@mui/material";
import { useState } from "react";

import ScanIcon from '@mui/icons-material/Radar'
import { useSnackbar } from "../../../../hooks/snackbar";

export default function SettingsKaraoke() {
    const [scanning, setScanning] = useState<boolean>(false);
    const {showSnackbar} = useSnackbar();

    const rescanSongs = async (e: any) => {
        setScanning(true);
        try {
            const resp = await fetch('/api/modules/karaoke/rescan', {
                method: 'POST',
            });
    
            if (resp.status != 200) {
                showSnackbar('Failed to re-scan songs: ' + (await resp.text()), 'error');
            } else {
                showSnackbar('Songs re-scanned', 'success');
            }
        } catch {
            showSnackbar('Failed to re-scan songs', 'error');
        }
        setScanning(false);
    };

    return <>
        <Button variant="contained" component="label" color="primary" onClick={rescanSongs} disabled={scanning}>
            <ScanIcon/> Re-scan songs
        </Button>
    </>
}