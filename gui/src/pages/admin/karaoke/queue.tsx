import { Card, CardContent, LinearProgress, List, Slider, Stack, Typography } from "@mui/material";
import { useAdminSocket } from "../../../hooks/adminSocket"
import Song from "./song";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

const secondsToDisplay = (seconds: number) => {
    if (seconds < 0) {
        return '--:--'
    }

    const remainingSeconds = Math.floor(seconds) % 60;
    const remainingMinutes = Math.floor(seconds/60);

    return `${remainingMinutes}`.padStart(2, '0') + ':' + `${remainingSeconds}`.padStart(2, '0');
};

export default function KaraokeQueue() {
    const {t} = useTranslation();
    const [currentPosition, setCurrentPosition] = useState<number>(0);
    const [duration, setDuration] = useState<number>(-1);

    const {appState, lastMessage} = useAdminSocket();
    const module = appState.modules.karaoke;

    useEffect(() => {
        if (lastMessage?.type == 'karaoke/PLAYING_STATUS') {
            setCurrentPosition(lastMessage.payload['current']);
            setDuration(lastMessage.payload.total);
        }
    }, [lastMessage]);

    return <Stack direction="column" gap={1} mt={1} flex="1 1 0">
        {
            module.currentSong && <Card elevation={2}>
                <CardContent>
                    <Stack gap={3}>
                        <Typography variant="h4">{t('karaoke.current')}:</Typography>
                        <Song song={module.currentSong} type="QUEUE" mb={0} />
                        <Stack direction="row" gap={3}>
                            <Typography variant={"body1"}>{secondsToDisplay(currentPosition)}</Typography>
                            <Slider disabled min={0} max={duration} value={currentPosition}/>
                            <Typography variant={"body1"}>{secondsToDisplay(duration)}</Typography>
                        </Stack>
                    </Stack>
                </CardContent>
            </Card>
        }

        <Typography variant="h4">{t('karaoke.queue')}:</Typography>
        {module.queue.length == 0 && <Typography variant="body1">{t('karaoke.empty_queue')}</Typography>}
        <Stack flex="1 1 0" style={{overflowY: 'scroll'}}>
            <List>
                { module.queue.map((x, idx) => <Song 
                    key={x.id}
                    song={x}
                    type="QUEUE" 
                    isFirst={idx === 0}
                    isLast={idx === module.queue.length - 1}
                />)
                }
            </List>
        </Stack>
    </Stack>
}