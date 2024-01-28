import { Button, Typography } from "@mui/material";
import { useState } from "react";

import UploadIcon from '@mui/icons-material/Upload'

export default function ImportPhk() {
    const [uploadStatus, setUploadStatus] = useState<String>("Waiting to upload a song...");

    const uploadPhk = async (e: React.ChangeEvent<HTMLInputElement>) => {
        setUploadStatus('Uploading song...');
        if (!e.target.files || e.target.files.length != 1) {
            setUploadStatus('Something went wrong!');
            return;
        }
    };

    return <>
        <Typography variant="body1">{uploadStatus}</Typography>
        <Button variant="contained" component="label" color="primary">
            <UploadIcon/> Upload a song
            <input type="file" hidden onChange={uploadPhk} />
        </Button>
    </>
}