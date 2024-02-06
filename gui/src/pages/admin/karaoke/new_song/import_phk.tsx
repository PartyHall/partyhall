import { Button, Typography } from "@mui/material";
import { useState } from "react";

import UploadIcon from '@mui/icons-material/Upload'
import { useTranslation } from "react-i18next";

export default function ImportPhk() {
    const {t} = useTranslation();
    const [uploadStatus, setUploadStatus] = useState<String>(t('karaoke.waiting_song_upload'));

    const uploadPhk = async (e: React.ChangeEvent<HTMLInputElement>) => {
        setUploadStatus(t('karaoke.upload_in_progress'));
        if (!e.target.files || e.target.files.length != 1) {
            setUploadStatus(t('general.something_went_wrong'));
            return;
        }
    };

    return <>
        <Typography variant="body1">{uploadStatus}</Typography>
        <Button variant="contained" component="label" color="primary">
            <UploadIcon/> {t('karaoke.upload_a_song')}
            <input type="file" hidden onChange={uploadPhk} />
        </Button>
    </>
}