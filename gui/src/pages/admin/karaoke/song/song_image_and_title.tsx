import { Stack, Typography } from "@mui/material";

import { KaraokeSong, KaraokeSongSession } from "../../../../types/appstate";
import { useTranslation } from "react-i18next";
import { SOCKET_MODE_DEBUG } from "../../../../hooks/useApi";

import SongCoverImage from "./song_img";

import {
    Piano as InstruIcon,
    RecordVoiceOver as VocalsIcon,
    MusicVideo as FullIcon,
} from '@mui/icons-material';

export default function SongImageAndTitle({song, session}: {song: KaraokeSong, session?: KaraokeSongSession}) {
    const {t} = useTranslation();

    return <>
        <SongCoverImage song={song} />

        <Stack direction="column" flex={1} ml={2}>
            <Typography variant="body1">{song.title}</Typography>
            <Typography variant="body1" fontSize=".9em" color="GrayText">{song.artist}</Typography>
            <Typography variant="body1" fontSize=".6em" color="GrayText">{song.format}</Typography>
            {
                session?.sung_by
                && <Typography variant="body1" fontSize="1em" color="GrayText">
                    {t('karaoke.sung_by')} {session?.sung_by}
                </Typography>
            }
            {
                SOCKET_MODE_DEBUG && <Typography variant="body1" fontSize="1em" color="GrayText">
                    {
                        session?.id ? `Session id ${session?.id}` : "Not a session"
                    }
                </Typography>
            }
            <Stack direction="row" gap={1}>
                <InstruIcon />
                <VocalsIcon color={song.has_vocals ? "action" : "disabled"} />
                <FullIcon color={song.has_full ? "action" : "disabled"} />
            </Stack>
        </Stack>
    </>
}