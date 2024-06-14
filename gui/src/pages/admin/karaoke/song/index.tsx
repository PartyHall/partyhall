import { IconButton, Stack, StackProps } from "@mui/material";
import { KaraokeSong, KaraokeSongSession } from "../../../../types/appstate";

import { useAdminSocket } from "../../../../hooks/adminSocket";
import { useConfirmDialog } from "../../../../hooks/dialog";
import { useTranslation } from "react-i18next";

import SongImageAndTitle from "./song_image_and_title";
import QueueMoveButtons from "./queue_move_buttons";
import QueueAddSong from "./queue_add_song";
import DeleteFromQueueButton from "./queue_delete_session_button";

import { PlayArrow as PlayIcon, Pause as PauseIcon } from '@mui/icons-material';

interface SongProps extends StackProps {
    song: KaraokeSong;
    session?: KaraokeSongSession;
    isFirst?: boolean;
    isLast?: boolean;
    type: 'SEARCH' | 'QUEUE';
}

export default function Song(props: SongProps) {
    const { appState } = useAdminSocket();
    const { t } = useTranslation();
    const { sendMessage } = useAdminSocket();

    const song = props.song;

    const { showDialog } = useConfirmDialog();

    const module = appState.modules.karaoke;
    const isCurrentSong = !(!module.currentSong || module.currentSong.song.uuid !== song.uuid);

    const playButtonClicked = () => {
        const action = props.type === 'SEARCH' ? 'QUEUE_AND_PLAY' : 'PLAY';
        const value = props.type === 'SEARCH' ? song.uuid : props.session?.id;

        if (module.currentSong != null) {
            showDialog(
                t('karaoke.skip.title'),
                t('karaoke.skip.content'),
                t('karaoke.skip.skip_and_play'),
                async () => sendMessage('karaoke/' + action, value)
            )
        } else {
            sendMessage('karaoke/' + action, value);
        }
    }

    return <Stack direction="row" mb={2} {...props}>
        {
            /** QUEUE: Move up & down a session in the queue */
            props.session
            && props.type === 'QUEUE'
            && (!module.currentSong || module.currentSong.id !== props.session.id)
            && <QueueMoveButtons
                session={props.session}
                isFirst={props.isFirst || false}
                isLast={props.isLast || false}
            />
        }

        <SongImageAndTitle song={song} session={props.session} />

        <Stack direction="column" alignItems="center">
            {
                /** SEARCH | OTHER SONG: Play button */
                appState.current_mode === 'KARAOKE'
                && (props.type === 'SEARCH' || !isCurrentSong)
                && <IconButton onClick={playButtonClicked}>
                    <PlayIcon />
                </IconButton>
            }
            {
                /** Queue & CURRENT SONG: Pause button */
                appState.current_mode === 'KARAOKE'
                && props.session
                && props.type === 'QUEUE'
                && isCurrentSong
                && <IconButton onClick={() => sendMessage('karaoke/PAUSE', props.session?.id)}>
                    {!module.started && <PlayIcon />}
                    {module.started && <PauseIcon />}
                </IconButton>
            }

            {
                /** SEARCH: Add a song to the queue */
                props.type === 'SEARCH'
                && <QueueAddSong
                    currentSong={module.currentSong}
                    queue={module.queue}
                    song={song}
                />
            }

            {
                /** QUEUE: Delete from queue button */
                props.session
                && props.type === 'QUEUE'
                && <DeleteFromQueueButton
                    song={props.song}
                    session={props.session}
                />
            }
        </Stack>
    </Stack>
}