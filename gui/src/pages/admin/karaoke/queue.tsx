import { Card, CardContent, List, Slider, Stack, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { useAdminSocket } from "../../../hooks/adminSocket"

import Song from "./song";

const secondsToDisplay = (seconds: number) => {
    if (seconds < 0) {
        return '--:--'
    }

    const remainingSeconds = Math.floor(seconds) % 60;
    const remainingMinutes = Math.floor(seconds / 60);

    return `${remainingMinutes}`.padStart(2, '0') + ':' + `${remainingSeconds}`.padStart(2, '0');
};

export default function KaraokeQueue() {
    const { t } = useTranslation();
    const [currentPosition, setCurrentPosition] = useState<number>(0);
    const [duration, setDuration] = useState<number>(-1);

    const { appState, lastMessage, sendMessage } = useAdminSocket();
    const module = appState.modules.karaoke;

    useEffect(() => {
        if (lastMessage?.type == 'karaoke/PLAYING_STATUS') {
            setCurrentPosition(lastMessage.payload['current']);
            setDuration(lastMessage.payload.total);
        }
    }, [lastMessage]);

    return <Stack direction="column" gap={1} mt={1}>
        {
            module.currentSong && <Card elevation={2}>
                <CardContent>
                    <Stack gap={3}>
                        <Typography variant="h4" textAlign="center" fontSize="1.3em">{t('karaoke.current')}</Typography>
                        <Song song={module.currentSong.song} session={module.currentSong} type="QUEUE" mb={0} />
                        <Stack direction="row" gap={3}>
                            <Typography variant={"body1"}>{secondsToDisplay(currentPosition)}</Typography>
                            <Slider disabled min={0} max={duration} value={currentPosition} />
                            <Typography variant={"body1"}>{secondsToDisplay(duration)}</Typography>
                        </Stack>

                        <Typography variant="h4" textAlign="center" fontSize="1.3em">{t('karaoke.volume.title')}</Typography>
                        <Stack direction="column">
                            <Typography variant="body1" textAlign="center">{t('karaoke.volume.instrumental')}</Typography>
                            <Slider
                                min={0}
                                max={1}
                                step={.1}
                                value={module.volumeInstru}
                                onChange={(_, val) => sendMessage('karaoke/VOLUME_INSTRU', val)}
                            />

                            {
                                module.currentSong.song.has_vocals && <>
                                    <Typography variant="body1" textAlign="center">{t('karaoke.volume.vocals')}</Typography>
                                    <Slider
                                        min={0}
                                        max={.5}
                                        step={.01}
                                        value={module.volumeVocals}
                                        onChange={(_, val) => sendMessage('karaoke/VOLUME_VOCALS', val)}
                                    />
                                </>
                            }

                            {
                                !module.currentSong.song.has_vocals && module.currentSong.song.has_full && <>
                                    <Typography variant="body1" textAlign="center">{t('karaoke.volume.full')}</Typography>
                                    <Slider
                                        min={0}
                                        max={1}
                                        step={.1}
                                        value={module.volumeFull}
                                        onChange={(_, val) => sendMessage('karaoke/VOLUME_FULL', val)}
                                    />
                                </>
                            }
                        </Stack>
                    </Stack>
                </CardContent>
            </Card>
        }

        <Typography variant="h4" fontSize="1.3em" textAlign="center">{t('karaoke.queue')}</Typography>
        {module.queue.length == 0 && <Typography variant="body1">{t('karaoke.empty_queue')}</Typography>}
        <List>
            {module.queue.map((x, idx) => <Song
                key={x.id}
                song={x.song}
                session={x}
                type="QUEUE"
                isFirst={idx === 0}
                isLast={idx === module.queue.length - 1}
            />)
            }
        </List>
    </Stack>
}