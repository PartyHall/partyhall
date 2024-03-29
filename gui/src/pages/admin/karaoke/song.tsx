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
import { useTranslation } from "react-i18next";

// interface SongProps extends React.HTMLAttributes<HTMLDivElement>{
interface SongProps extends StackProps {
    song: KaraokeSong;
    isFirst?: boolean;
    isLast?: boolean;
    type: 'SEARCH' | 'QUEUE';
}

export default function Song(props: SongProps) {
    const imgRef = useRef<HTMLImageElement>(null);
    const { appState } = useAdminSocket();
    const {t} = useTranslation();
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
                <IconButton style={{padding: ".25em"}} disabled={props.isFirst} onClick={() => {
                    sendMessage('karaoke/QUEUE_MOVE_UP', props.song.filename)
                }}>
                    <UpIcon />
                </IconButton>
                <IconButton style={{padding: ".25em"}} disabled={props.isLast} onClick={() => {
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
            {
                props.song.sung_by && props.song.sung_by.length > 0 &&
                <Typography variant="body1" fontSize="1em" color="GrayText">{t('karaoke.sung_by')} {props.song.sung_by}</Typography>
            }
        </Stack>
        <Stack direction="column" alignItems="center">
            {
                (appState.current_mode === 'KARAOKE' && (props.type === 'SEARCH' || !isCurrentSong)) &&
                <IconButton onClick={() => {
                    if (module.currentSong != null) {
                        showDialog(
                            t('karaoke.skip.title'),
                            t('karaoke.skip.content'),
                            t('karaoke.skip.skip_and_play'),
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
                (appState.current_mode === 'KARAOKE' && props.type === 'QUEUE' && isCurrentSong) &&
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
                    showSnackbar(t('karaoke.adding_to_queue'), 'success');
                }} disabled={module.queue.map(x => x.id).includes(props.song.id) || (module.currentSong && props.song.id == module.currentSong.id)}>
                    <AddToQueueIcon />
                </IconButton>
            }
            {
                props.type === 'QUEUE' &&
                <IconButton onClick={() => {
                    if (module.currentSong && module.currentSong.filename === props.song.filename) {
                        showDialog(
                            t('karaoke.current_remove.title'),
                            t('karaoke.current_remove.content', {name: songTitle(props.song)}),
                            t('karaoke.current_remove.remove'),
                            async () => sendMessage('karaoke/DEL_FROM_QUEUE', props.song.filename),
                        );
                    } else {
                        showDialog(
                            t('karaoke.skip.title'),
                            t('karaoke.skip.content'),
                            t('karaoke.skip.skip'),
                            async () => {
                                sendMessage('karaoke/DEL_FROM_QUEUE', props.song.filename);
                                showSnackbar(t('karaoke.queue_remove.removed'), 'error');
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