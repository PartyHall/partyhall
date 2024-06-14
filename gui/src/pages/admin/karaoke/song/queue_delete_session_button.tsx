import { IconButton } from "@mui/material";
import { useTranslation } from "react-i18next";
import { useConfirmDialog } from "../../../../hooks/dialog";
import { useAdminSocket } from "../../../../hooks/adminSocket";
import { useSnackbar } from "../../../../hooks/snackbar";
import { songTitle } from "../../../../utils/songs";
import { Close as DelFromQueueIcon } from '@mui/icons-material';
import { KaraokeSong, KaraokeSongSession } from "../../../../types/appstate";

export default function DeleteFromQueueButton({song, session}: {song: KaraokeSong, session: KaraokeSongSession}) {
    const {appState} = useAdminSocket();
    const {t} = useTranslation();
    const {showDialog} = useConfirmDialog();
    const {sendMessage} = useAdminSocket();
    const {showSnackbar} = useSnackbar();

    const currentSong = appState.modules.karaoke.currentSong;

    return <IconButton onClick={() => {
        if (currentSong?.id === session.id) {
            showDialog(
                t('karaoke.current_remove.title'),
                t('karaoke.current_remove.content', { name: songTitle(song) }),
                t('karaoke.current_remove.remove'),
                async () => sendMessage('karaoke/DEL_FROM_QUEUE', session.id),
            );
        } else {
            showDialog(
                t('karaoke.skip.title'),
                t('karaoke.skip.content'),
                t('karaoke.skip.skip'),
                async () => {
                    sendMessage('karaoke/DEL_FROM_QUEUE', session.id);
                    showSnackbar(t('karaoke.queue_remove.removed'), 'error');
                }
            );
        }
    }}>
        <DelFromQueueIcon />
    </IconButton>;
}