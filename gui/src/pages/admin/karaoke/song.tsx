import { IconButton, Stack, StackProps, Typography } from "@mui/material";
import { KaraokeSong } from "../../../types/appstate";

import {
    PlayArrow as PlayIcon,
    Pause as PauseIcon,
    Add as AddToQueueIcon,
    Close as DelFromQueueIcon,
    ArrowUpward as UpIcon,
    ArrowDownward as DownIcon,
} from '@mui/icons-material';
import { useAdminSocket } from "../../../hooks/adminSocket";
import { useSnackbar } from "../../../hooks/snackbar";
import { useRef } from "react";
import { useConfirmDialog } from "../../../hooks/dialog";
import { songTitle } from "../../../utils/songs";

// interface SongProps extends React.HTMLAttributes<HTMLDivElement>{
interface SongProps extends StackProps {
    song: KaraokeSong;
    first?: boolean;
    last?: boolean;
    type: 'SEARCH' | 'QUEUE';
}

export default function Song(props: SongProps) {
    const imgRef = useRef<HTMLImageElement>(null);
    const { appState } = useAdminSocket();
    const { sendMessage } = useAdminSocket();
    const { showSnackbar } = useSnackbar();

    const {showDialog} = useConfirmDialog();

    const onImgError = () => {
        if (!imgRef.current) {
            return;
        }

        imgRef.current.src = '/api/modules/karaoke/fallback-image';
    };

    const module = appState.modules.karaoke;
    const isCurrentSong = !(!module.currentSong || module.currentSong.filename !== props.song.filename);

    return <Stack direction="row" mb={2} {...props}>
        {
            (props.type === 'QUEUE' && (!module.currentSong || module.currentSong.filename !== props.song.filename)) &&
            <Stack direction="column" alignItems="center" justifyContent="center" gap={1}>
                <IconButton style={{padding: ".25em"}} disabled={props.first} onClick={() => {
                    sendMessage('karaoke/QUEUE_MOVE_UP', props.song.filename)
                }}>
                    <UpIcon />
                </IconButton>
                <IconButton style={{padding: ".25em"}} disabled={props.last} onClick={() => {
                    sendMessage('karaoke/QUEUE_MOVE_DOWN', props.song.filename)
                }}>
                    <DownIcon />
                </IconButton>
            </Stack>
        }
        <img ref={imgRef} onError={onImgError} src={`/api/modules/karaoke/medias/${props.song.filename}/cover.jpg`} alt={props.song.filename} style={{ maxHeight: '6em', objectFit: 'contain' }} />
        <Stack direction="column" flex={1} ml={2}>
            <Typography variant="body1">{props.song.title.length > 0 ? props.song.title : props.song.filename}</Typography>
            <Typography variant="body1" fontSize=".9em" color="GrayText">{props.song.artist}</Typography>
            <Typography variant="body1" fontSize=".6em" color="GrayText">{props.song.format}</Typography>
        </Stack>
        <Stack direction="column" alignItems="center">
            {
                (props.type === 'SEARCH' || !isCurrentSong) &&
                <IconButton onClick={() => {
                    if (module.currentSong != null) {
                        showDialog(
                            "Skipping song?",
                            "There is a song currently playing. You might want to put it in queue instead.",
                            "Skip & Play",
                            async () => sendMessage('karaoke/PLAY', props.song.filename)
                        )
                    } else {
                        sendMessage('karaoke/PLAY', props.song.filename);
                    }
                }}>
                    <PlayIcon />
                </IconButton>
            }
            {
                (props.type === 'QUEUE' && isCurrentSong) &&
                <IconButton onClick={() => {
                    sendMessage('karaoke/PAUSE', props.song.filename);
                }}>
                    {!module.started && <PlayIcon />}
                    {module.started && <PauseIcon />}
                </IconButton>
            }
            {
                props.type === 'SEARCH' &&
                <IconButton onClick={() => {
                    sendMessage('karaoke/ADD_TO_QUEUE', props.song.filename);
                    showSnackbar('Adding song to the queue', 'success');
                }}>
                    <AddToQueueIcon />
                </IconButton>
            }
            {
                props.type === 'QUEUE' &&
                <IconButton onClick={() => {
                    if (module.currentSong && module.currentSong.filename === props.song.filename) {
                        showDialog(
                            "Skip current song",
                            "Are you sure you want to skip the currently playing song?",
                            "Skip",
                            async () => sendMessage('karaoke/DEL_FROM_QUEUE', props.song.filename),
                        );
                    } else {
                        showDialog(
                            "Remove from the queue",
                            `Are you sure you want to remove "${songTitle(props.song)}" from the queue?`,
                            "Remove",
                            async () => {
                                sendMessage('karaoke/DEL_FROM_QUEUE', props.song.filename);
                                showSnackbar('Removing song from the queue', 'error');
                            }
                        );
                    }
                }}>
                    <DelFromQueueIcon />
                </IconButton>
            }
        </Stack>
    </Stack>
}