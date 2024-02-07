import { Stack, Typography } from "@mui/material";
import { KaraokeSong } from "../../../types/appstate";

import { useRef } from "react";
import { useTranslation } from "react-i18next";

export default function OsdSong({ song }: { song: KaraokeSong }) {
    const imgRef = useRef<HTMLImageElement>(null);
    const {t} = useTranslation();

    const onImgError = () => {
        if (!imgRef.current) {
            return;
        }

        imgRef.current.src = '/api/modules/karaoke/fallback-image';
    };

    return <Stack direction="row" gap={2}>
        <img ref={imgRef} onError={onImgError} src={`/api/modules/karaoke/medias/${song.filename}/cover.jpg`} alt={song.filename} style={{ maxHeight: '6em', objectFit: 'contain' }} />
        <Stack direction="column" flex={1}>
            <Typography variant="body1">{song.title.length > 0 ? song.title : song.filename}</Typography>
            <Typography variant="body1" fontSize=".9em" color="GrayText">{song.artist}</Typography>
            {
                song.sung_by && song.sung_by.length > 0 &&
                <Typography variant="body1" fontSize="1em" color="GrayText">{t('karaoke.sung_by')} {song.sung_by}</Typography>
            }
        </Stack>
    </Stack>
}