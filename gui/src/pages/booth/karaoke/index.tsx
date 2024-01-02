import { useRef } from "react";
import { useBoothSocket } from "../../../hooks/boothSocket";

import '../../../assets/css/karaoke.scss';
import CDGPlayer from "./cdgplayer";
import Webcam from "react-webcam";
import { Stack, Typography } from "@mui/material";
import { songTitle } from "../../../utils/songs";
import OsdSong from "./osd_song";

export default function Karaoke() {
    const { appState, lastMessage, sendMessage } = useBoothSocket();
    const webcamRef = useRef<Webcam>(null);
    const module = appState.modules.karaoke;

    return <div className="karaoke">
        <Webcam
            forceScreenshotSourceSize
            ref={webcamRef}
            width={appState.modules.photobooth.webcam_resolution.width}
            height={appState.modules.photobooth.webcam_resolution.height}
            screenshotFormat="image/jpeg"
            videoConstraints={{ facingMode: 'user', ...appState.modules.photobooth.webcam_resolution }}
            className='karaoke__webcam'
        />
        {
            module.currentSong && module.preplayTimer == 0 &&
            <>
                {
                    module.currentSong.format === 'cdg' &&
                    <CDGPlayer
                        cdgAlpha={.8}
                        cdgSize={window.innerHeight / 2}
                        width={window.innerWidth/2}
                        height={window.innerHeight / 2}
                        isPlaying={module.started}
                        song={module.currentSong}
                        onEnd={() => sendMessage('karaoke/PLAYING_ENDED')}
                        onError={() => {}}
                        onLoad={() => {}}
                        onPlay={() => {}}
                        onStatus={(x: any) => sendMessage('karaoke/PLAYING_STATUS', {'current': x.position, 'total': x.total})}
                    /> 
                }
            </>
        }
        {
            module.currentSong && module.preplayTimer > 0 &&
            <Stack display="column" className="karaoke__no_song">
                <Typography variant="h1">Now playing:</Typography>
                <Typography variant="h2">{songTitle(module.currentSong)}</Typography>
                <Typography variant="h3">{module.preplayTimer}</Typography>
            </Stack>
        }
        {
            !module.currentSong &&
            <Stack display="column" className="karaoke__no_song">
                <Typography variant="h1">No song playing !</Typography>
            </Stack>
        }
        {
            module.queue.length > 0 &&
            <Stack className="karaoke__next_song" gap={1}>
                <Typography variant="h3">Next up:</Typography>
                <OsdSong song={module.queue[0]} />
            </Stack>
        }
    </div>;
}