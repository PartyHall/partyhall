import { Button } from "@mui/material";
import { useState } from "react";

import ScanIcon from '@mui/icons-material/Radar'
import { useSnackbar } from "../../../../hooks/snackbar";
import { useApi } from "../../../../hooks/useApi";
import { useTranslation } from "react-i18next";

export default function SettingsKaraoke() {
    const {t} = useTranslation();
    const {api} = useApi();
    const [scanning, setScanning] = useState<boolean>(false);
    const {showSnackbar} = useSnackbar();

    const rescanSongs = async (e: any) => {
        setScanning(true);
        try {
            await api.karaoke.rescanSongs();
            showSnackbar(t('karaoke.songs_rescanned'), 'success');
        } catch (e) {
            showSnackbar(t('karaoke.song_rescan_failed') + ': ' + e, 'error');
        }
        setScanning(false);
    };

    return <>
        <Button variant="contained" component="label" color="primary" onClick={rescanSongs} disabled={scanning}>
            <ScanIcon/> {t('karaoke.rescan_songs')}
        </Button>
    </>
}