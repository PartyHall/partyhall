import { Button, FormControl, Input, InputLabel, MenuItem, Select, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { MusicNote as MusicNoteIcon, Lyrics as LyricsIcon } from '@mui/icons-material';

import SearchSpotify from "./search_spotify";
import { useSnackbar } from "../../../../hooks/snackbar";
import { useApi } from "../../../../hooks/useApi";
import { COVER_SOURCE, FORMAT, SongPost } from "../../../../sdk/requests/song";

export default function CreateSong() {
    const {api} = useApi();
    const {showSnackbar} = useSnackbar();
    const [coverSource, setCoverSource] = useState<COVER_SOURCE>('LINK');
    const [format, setFormat] = useState<FORMAT>('CDG');
    const [title, setTitle] = useState<string>("");
    const [artist, setArtist] = useState<string>("");

    const [songFilename, setSongFilename] = useState<string>("");
    const [cdgFilename, setCdgFilename] = useState<string>("");

    const { register, handleSubmit, resetField, setValue, reset } = useForm<SongPost>();

    const submit = async (data: SongPost) => {
        try {
            await api.karaoke.post(data);
            reset();
            setSongFilename('');
            setCdgFilename('');

            showSnackbar('Music created', 'success');
        } catch (e) {
            console.log(e);
            showSnackbar('Failed to create song: ' + e, 'error');
        }
    };

    return <form onSubmit={handleSubmit(submit)}>
        <Stack direction="column" gap={2}>
            <TextField required placeholder="Title" {...register('title')} onChange={x => {
                setTitle(x.target.value);
                register('title').onChange(x);
            }} />
            <TextField required placeholder="Artist" {...register('artist')} onChange={x => {
                setArtist(x.target.value);
                register('artist').onChange(x);
            }} />
            <FormControl fullWidth>
                <InputLabel id="label_select_cover_source">Cover source</InputLabel>
                <Select required defaultValue='LINK' placeholder='Cover source' labelId='label_select_cover_source' {...register('cover_type')} onChange={x => {
                    setCoverSource(x.target.value as COVER_SOURCE);
                    register('cover_type').onChange(x);
                }}>
                    <MenuItem value='NO_COVER'>No cover</MenuItem>
                    <MenuItem value='LINK'>Spotify</MenuItem>
                    <MenuItem value='UPLOADED'>Upload</MenuItem>
                </Select>
            </FormControl>

            { coverSource === 'LINK' && <SearchSpotify artist={artist} title={title} onChange={x => setValue('cover_data', x)}/> }

            <FormControl fullWidth>
                <InputLabel id="label_select_format">Format</InputLabel>
                <Select required defaultValue='CDG' placeholder='Format' labelId='label_select_format' {...register('format')} onChange={x => {
                    setFormat(x.target.value as FORMAT);
                    register('format').onChange(x)
                    resetField('song')
                    resetField('cdg');
                }}>
                    <MenuItem value='CDG'>MP3+CDG</MenuItem>
                    <MenuItem value='WEBM'>webm</MenuItem>
                    <MenuItem value='MP4'>mp4</MenuItem>
                </Select>
            </FormControl>

            <Stack direction="column" gap={2} justifyContent="center">
                <Button variant="contained" component="label" color="primary">
                    <MusicNoteIcon/> Upload {format === 'CDG' ? 'mp3' : format}
                    <input type="file" hidden required {...register('song')} onChange={x => {
                        register('song').onChange(x);
                        if (x.target.files && x.target.files.length > 0) {
                            setSongFilename(x.target.files[0].name);
                        }
                    }} />
                </Button>

                <Typography variant="body1" textAlign="center">{songFilename}</Typography>

                {
                    format === 'CDG' &&
                    <>
                        <Button variant="contained" component="label" color="primary">
                            <LyricsIcon/> Upload CDG
                            <input type="file" hidden required {...register('cdg')} onChange={x => {
                                register('cdg').onChange(x);
                                if (x.target.files && x.target.files.length > 0) {
                                    setCdgFilename(x.target.files[0].name);
                                }
                            }}/>
                        </Button>

                        <Typography variant="body1" textAlign="center">{cdgFilename}</Typography>
                    </>
                }
            </Stack>

            <Button
                variant='outlined'
                type='submit'
                disabled={songFilename.length === 0 || (format === 'CDG' && cdgFilename.length === 0)}
            >Add song</Button>
        </Stack>
    </form>;
}