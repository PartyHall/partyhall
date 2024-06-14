import { Stack, Typography } from "@mui/material";
import { KaraokeSongSession } from "../../../types/appstate";

import { useTranslation } from "react-i18next";
import SongCoverImage from "../../admin/karaoke/song/song_img";

export default function OsdSong({ session }: { session: KaraokeSongSession }) {
    const {t} = useTranslation();

    const song = session.song;

    return <Stack direction="row" gap={2}>
        <SongCoverImage song={song} />

        <Stack direction="column" flex={1}>
            <Typography variant="body1">{song.title}</Typography>
            <Typography variant="body1" fontSize=".9em" color="GrayText">{song.artist}</Typography>
            {
                session.sung_by && session.sung_by.length > 0 &&
                <Typography variant="body1" fontSize="1em" color="GrayText">{t('karaoke.sung_by')} {session.sung_by}</Typography>
            }
        </Stack>
    </Stack>
}