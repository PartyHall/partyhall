import { Button, FormControl, Input, InputLabel, MenuItem, Select, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { MusicNote as MusicNoteIcon, Lyrics as LyricsIcon } from '@mui/icons-material';

import SearchSpotify from "./search_spotify";
import { useSnackbar } from "../../../../hooks/snackbar";

type COVER_SOURCE = 'NO_COVER'|'LINK'|'UPLOADED';
type FORMAT = 'CDG'|'WEBM'|'MP4';

type NewSong = {
    title: string;
    artist: string;
    format: string;
    cover_type: COVER_SOURCE;
    cover_data: string;

    song: FileList;
    cdg: FileList;
};

export default function CreateSong() {
    const {showSnackbar} = useSnackbar();
    const [coverSource, setCoverSource] = useState<COVER_SOURCE>('LINK');
    const [format, setFormat] = useState<FORMAT>('CDG');
    const [title, setTitle] = useState<string>("");
    const [artist, setArtist] = useState<string>("");

    const [songFilename, setSongFilename] = useState<string>("");
    const [cdgFilename, setCdgFilename] = useState<string>("");

    const { register, handleSubmit, resetField, setValue, reset } = useForm<NewSong>();

    const submit = async (data: NewSong) => {
        const fd = new FormData();
        fd.append('title', data.title);
        fd.append('artist', data.artist);
        fd.append('format', data.format);
        fd.append('cover_type', data.cover_type)

        if (data.cover_type === 'LINK') {
            fd.append('cover_url', data.cover_data);
        } else if (data.cover_type === 'UPLOADED') {

        }

        fd.append('song', data.song[0]);

        if (data.format === 'CDG') {
            fd.append('cdg', data.cdg[0]);
        }

        try {
            const resp = await fetch('/api/modules/karaoke/song', {
                method: 'POST',
                body: fd,
            });
    
            if (resp.status != 201) {
                showSnackbar('Failed to create song: ' + (await resp.text()), 'error');
            } else {
                showSnackbar('Music created', 'success');
                reset();
            }
        } catch (e) {
            console.log(e);
            //@ts-ignore
            showSnackbar('Failed to create song: ' + (await e.response.text()), 'error');
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