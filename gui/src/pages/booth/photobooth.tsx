import { useEffect, useRef, useState } from "react";
import Webcam from "react-webcam";
import LockedModal from "../../components/locked_modal";
import { useBoothSocket } from "../../hooks/boothSocket";

import '../../assets/css/photobooth.scss';

type LastPicture = {
    url: string;
    loaded: boolean;
};

export default function Photobooth() {
    const webcamRef = useRef<Webcam>(null);
    const { appState, lastMessage, sendMessage } = useBoothSocket();
    const [timer, setTimer] = useState(-1);
    const [flash, setFlash] = useState<boolean>(false);
    const [lastPicture, setLastPicture] = useState<LastPicture|null>(null);

    const module = appState.modules.photobooth;
    const resolution = module.webcam_resolution;

    const takePicture = async (unattended: boolean) => {
        if (!webcamRef || !webcamRef.current) {
            return;
        }

        const imageSrc = webcamRef.current.getScreenshot();
        if (imageSrc) {
            let form = new FormData();

            form.append('image', imageSrc);
            form.append('unattended', unattended ? 'true' : 'false')
            form.append('event', ''+appState?.app_state?.current_event?.id)

            try {
                const resp = await fetch('/api/picture', {
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

            setTimer(module.default_timer)
            return
        }

        if (lastMessage.type == 'UNATTENDED_PICTURE') {
            takePicture(true);
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
                takePicture(false);
                setFlash(false);
                setTimer(-1);
            }, 500);
        }
    }, [timer]);

    return <div className="photobooth">
        <Webcam
            forceScreenshotSourceSize
            ref={webcamRef}
            width={resolution.width}
            height={resolution.height}
            onClick={() => appState.current_mode !== 'DISABLED' && sendMessage('TAKE_PICTURE')}
            screenshotFormat="image/jpeg"
            videoConstraints={{ facingMode: 'user', ...resolution }}
        />

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