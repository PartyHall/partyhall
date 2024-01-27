import { Button, FormControl, InputLabel, MenuItem, Select, Stack, TextField } from "@mui/material";
import { useForm } from "react-hook-form";

type NewSong = {
    title: string;
    artist: string;
    format: string;
    cover_type: 'NO_COVER' | 'SPOTIFY' | 'UPLOADED';
    cover_data: string;
};

export default function KaraokeAddSong() {
    const { register, handleSubmit } = useForm<NewSong>();

    const submit = async (data: NewSong) => {
        console.log(data);
    };

    return <form onSubmit={handleSubmit(submit)}>
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