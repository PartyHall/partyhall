import { IconButton } from "@mui/material";
import { useAdminSocket } from "../../../../hooks/adminSocket";
import { useSnackbar } from "../../../../hooks/snackbar";
import { useTranslation } from "react-i18next";
import { KaraokeSong, KaraokeSongSession } from "../../../../types/appstate";

import { Add as AddToQueueIcon } from '@mui/icons-material';

type Props = {
    currentSong?: KaraokeSongSession;
    song: KaraokeSong;
    queue: KaraokeSongSession[];
};

export default function QueueAddSong({currentSong, song, queue}: Props) {
    const {t} = useTranslation();
    const {sendMessage} = useAdminSocket();
    const {showSnackbar} = useSnackbar();

    const isCurrentSong = currentSong?.song?.uuid === song.uuid;
    const isAlreadyInQueue = isCurrentSong || queue.map(x => x.song.uuid).includes(song.uuid);

    const addSongToQueue = () => {
        sendMessage('karaoke/ADD_TO_QUEUE', song.uuid);
        showSnackbar(t('karaoke.adding_to_queue'), 'success');
    };

    return <IconButton onClick={addSongToQueue} disabled={isAlreadyInQueue}>
        <AddToQueueIcon />
    </IconButton>;
}