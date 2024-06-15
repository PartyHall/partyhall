import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import Webcam from "react-webcam";
import { Stack, Typography } from "@mui/material";

import CDGPlayer from "./cdgplayer";
import VideoPlayer from "./videoplayer";
import OsdSong from "./osd_song";
import { b64ImageToBlob } from "../../../utils/files";
import { songTitle } from "../../../utils/songs";
import LockedModal from "../../../components/locked_modal";
import { useBoothSocket } from "../../../hooks/boothSocket";

import '../../../assets/css/photobooth.scss';
import '../../../assets/css/karaoke.scss';

type LastPicture = {
    url: string;
    loaded: boolean;
};

export default function Photobooth() {
    const webcamRef = useRef<Webcam>(null);
    const {t} = useTranslation();
    const { appState, lastMessage, sendMessage } = useBoothSocket();
    const [timer, setTimer] = useState(-1);
    const [flash, setFlash] = useState<boolean>(false);
    const [lastPicture, setLastPicture] = useState<LastPicture|null>(null);

    const module_photobooth = appState.modules.photobooth;
    const module_karaoke = appState.modules.karaoke;
    const resolution = module_photobooth.webcam_resolution;

    const takePicture = async (unattended: boolean, karaoke: boolean) => {
        if (!webcamRef || !webcamRef.current) {
            return;
        }

        const imageSrc = webcamRef.current.getScreenshot();
        if (imageSrc) {
            let form = new FormData();

            form.append('image', b64ImageToBlob(imageSrc));
            if (!karaoke) {
                form.append('unattended', unattended ? 'true' : 'false')
                form.append('event', ''+appState?.app_state?.current_event?.id)
            }

            const endpoint = `/api/modules/${karaoke ? 'karaoke' : 'photobooth'}/picture`;

            try {
                const resp = await fetch(endpoint, {
                    method: 'POST',
                    body: form,
                });
         
                setTimer(-1);

                if (!unattended) {
                    const image = await resp.blob();
                    const url = URL.createObjectURL(image);

                    setLastPicture({ url, loaded: false});
                    setTimeout(() => setLastPicture(null), 3500);
                }
            } catch {
                setTimer(-1);
            }
        }
    };

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        if (lastMessage.type == 'TAKE_PICTURE' && timer === -1) {
            if (appState.current_mode === 'DISABLED') {
                return;
            }

            setTimer(module_photobooth.default_timer)
            return
        }

        if (lastMessage.type == 'UNATTENDED_PICTURE') {
            takePicture(true, false);
        }

        if (lastMessage.type == 'UNATTENDED_KARAOKE_PICTURE') {
            takePicture(true, true);
        }
    }, [lastMessage]);

    useEffect(() => {
        if (timer > 0) {
            setTimeout(() => {
                setTimer(timer-1);
            }, 1000);
        }

        if (timer == 0) {
            setFlash(true);
            setTimeout(() => {
                takePicture(false, false);
                setFlash(false);
                setTimer(-1);
            }, 500);
        }
    }, [timer]);

    return <div className="karaoke">
        <Webcam
            forceScreenshotSourceSize
            ref={webcamRef}
            width={resolution.width}
            height={resolution.height}
            onClick={() => appState.current_mode !== 'DISABLED' && sendMessage('photobooth/TAKE_PICTURE')}
            screenshotFormat="image/jpeg"
            videoConstraints={{ facingMode: 'user', ...resolution }}
        />

        {
            module_karaoke.currentSong && module_karaoke.preplayTimer == 0 &&
            <>
                {
                    module_karaoke.currentSong.song.format.toLowerCase() === 'cdg' &&
                    <CDGPlayer
                        cdgAlpha={.8}
                        cdgSize={window.innerHeight / 2}
                        width={window.innerWidth/2}
                        height={window.innerHeight / 2}
                        isPlaying={module_karaoke.started}
                        song={module_karaoke.currentSong.song}
                        onEnd={() => sendMessage('karaoke/PLAYING_ENDED', module_karaoke.currentSong.id)}
                        onError={() => {}}
                        onLoad={() => {}}
                        onPlay={() => {}}
                        onStatus={(x: any) => sendMessage('karaoke/PLAYING_STATUS', {
                            'session_id': module_karaoke.currentSong.id,
                            'current': x.position,
                            'total': x.total,
                        })}
                    /> 
                }
                {
                    module_karaoke.currentSong.song.format.toLowerCase() !== 'cdg' && <VideoPlayer
                        isPlaying={module_karaoke.started}
                        song={module_karaoke.currentSong.song}
                        onEnd={() => sendMessage('karaoke/PLAYING_ENDED', module_karaoke.currentSong.id)}
                        onStatus={(x: any) => sendMessage('karaoke/PLAYING_STATUS', {
                            'session_id': module_karaoke.currentSong.id,
                            'current': x.position,
                            'total': x.total,
                        })}
                    />
                }
            </>
        }
        {
            module_karaoke.currentSong && module_karaoke.preplayTimer > 0 &&
            <Stack display="column" className="karaoke__no_song">
                <Typography variant="h1">{t('karaoke.now_playing')}:</Typography>
                <Typography variant="h2">{songTitle(module_karaoke.currentSong.song)}</Typography>
                <Typography variant="h3">{module_karaoke.preplayTimer}</Typography>
                {
                    module_karaoke.currentSong.sung_by && module_karaoke.currentSong.sung_by.length > 0 &&
                    <Typography variant="h2">{t('karaoke.sung_by')} {module_karaoke.currentSong.sung_by}</Typography>
                }
            </Stack>
        }
        {
            module_karaoke.queue.length > 0 &&
            <Stack className="karaoke__next_song" gap={1}>
                <Typography variant="h3">{t('karaoke.next_up')}:</Typography>
                <OsdSong session={module_karaoke.queue[0]} />
            </Stack>
        }

        { timer >= 0 && <div className={`timer`}>{timer > 0 && timer}</div> }
        { flash && <div className="timer flash"></div> }
        { appState.current_mode === 'DISABLED' && <LockedModal /> }

        {
            lastPicture && <div className="picture_frame" style={lastPicture.loaded ? {} : {display: 'none'}}>
                <img src={lastPicture.url} onLoad={() => setLastPicture({...lastPicture, loaded: true})} alt="Last picture" />
            </div>
        }
    </div>
}