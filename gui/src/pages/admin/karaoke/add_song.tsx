import { Box, Button, FormControl, InputLabel, MenuItem, Select, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { useForm } from "react-hook-form";

import UploadIcon from '@mui/icons-material/Upload'

type METHOD = 'CREATE'|'IMPORT';

type NewSong = {
    title: string;
    artist: string;
    format: string;
    cover_type: 'NO_COVER' | 'SPOTIFY' | 'UPLOADED';
    cover_data: string;
};

export default function KaraokeAddSong() {
    const [uploadStatus, setUploadStatus] = useState<String>("Waiting to upload a song...");
    const [method, setMethod] = useState<METHOD>('CREATE');
    const { register, handleSubmit } = useForm<NewSong>();

    const submit = async (data: NewSong) => {
        console.log(data);
    };

    const uploadPhk = async (e: React.ChangeEvent<HTMLInputElement>) => {
        setUploadStatus('Uploading song...');
        if (!e.target.files || e.target.files.length != 1) {
            setUploadStatus('Something went wrong!');
            return;
        }
    };

    return <Stack direction="column" gap={2}>
        <Box mb={3}>
            <FormControl fullWidth>
                <InputLabel id="label_select_method">What to do?</InputLabel>
                <Select labelId="label_select_method" value={method} onChange={x => setMethod(x.target.value as METHOD)} style={{marginBottom: 3}}>
                    <MenuItem value='CREATE'>Create</MenuItem>
                    <MenuItem value='IMPORT'>Import .phk</MenuItem>
                </Select>
            </FormControl>
        </Box>

        {
            method === 'IMPORT' &&
            <>
                <Typography variant="body1">{uploadStatus}</Typography>
                <Button variant="contained" component="label" color="primary">
                    <UploadIcon/> Upload a song
                    <input type="file" hidden onChange={uploadPhk} />
                </Button>
            </>
        }
        {
            method === 'CREATE' &&
            <form onSubmit={handleSubmit(submit)}>
                <Stack direction="column" gap={2}>
                    <TextField required placeholder="Title" {...register('title')} />
                    <TextField required placeholder="Artist" {...register('artist')} />
                    <FormControl fullWidth>
                        <InputLabel id="label_select_format">Format</InputLabel>
                        <Select required defaultValue='CDG' placeholder='Format' labelId='label_select_format' {...register('format')}>
                            <MenuItem value='CDG'>MP3+CDG</MenuItem>
                            <MenuItem value='WEBM'>webm</MenuItem>
                            <MenuItem value='MP4'>mp4</MenuItem>
                        </Select>
                    </FormControl>
                    <FormControl fullWidth>
                        <InputLabel id="label_select_cover_source">Cover source</InputLabel>
                        <Select required defaultValue='SPOTIFY' placeholder='Cover source' labelId='label_select_cover_source' {...register('cover_type')}>
                            <MenuItem value='NO_COVER'>No cover</MenuItem>
                            <MenuItem value='SPOTIFY'>Spotify</MenuItem>
                            <MenuItem value='UPLOADED'>Upload</MenuItem>
                        </Select>
                    </FormControl>

                    {
                        /**
                         * TODO: Add search song cover on Spotify
                         *       Add upload song cover
                         *       Add upload file (mp3+cdg / webm / mp4)
                         */
                    }

                    <Button variant='outlined' type='submit'>Add song</Button>
                </Stack>
            </form>
        }
    </Stack>;
}