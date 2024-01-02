import { Card, CardContent, LinearProgress, List, Slider, Stack, Typography } from "@mui/material";
import { useAdminSocket } from "../../../hooks/adminSocket"
import Song from "./song";
import { useEffect, useState } from "react";

const secondsToDisplay = (seconds: number) => {
    const remainingSeconds = Math.floor(seconds) % 60;
    const remainingMinutes = Math.floor(seconds/60);

    return `${remainingMinutes}`.padStart(2, '0') + ':' + `${remainingSeconds}`.padStart(2, '0');
};

export default function KaraokeQueue() {
    const [currentPosition, setCurrentPosition] = useState<number>(0);
    const [duration, setDuration] = useState<number>(300);

    const {appState, lastMessage} = useAdminSocket();
    const module = appState.modules.karaoke;

    useEffect(() => {
        if (lastMessage?.type == 'karaoke/PLAYING_STATUS') {
            setCurrentPosition(lastMessage.payload['current']);
            if (!isNaN(lastMessage.payload.total)) {
                setDuration(lastMessage.payload.total);
            }
        }
    }, [lastMessage]);

    return <Stack direction="column" gap={1} mt={1}>
        {
            module.currentSong && <Card elevation={2}>
                <CardContent>
                    <Stack gap={3}>
                        <Typography variant="h4">Current:</Typography>
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

        <Typography variant="h4">Queue:</Typography>
        {module.queue.length == 0 && <Typography variant="body1">Empty queue</Typography>}
        <List>
            { module.queue.map((x, idx) => <Song 
                key={x.id}
                song={x}
                type="QUEUE" 
                first={idx === 0}
                last={idx === module.queue.length - 1}
            />)
            }
        </List>
    </Stack>
}