import { useEffect, useRef, useState } from "react";
import Webcam from "react-webcam";
import LockedModal from "../../components/locked_modal";
import { useWebsocket } from "../../hooks/boothSocket";

import '../../assets/css/photobooth.scss';

export default function Photobooth() {
    const webcamRef = useRef<Webcam>(null);
    const { appState, lastMessage, sendMessage } = useWebsocket();
    const [timer, setTimer] = useState(-1);
    const [flash, setFlash] = useState<boolean>(false);
    const [lastPicture, setLastPicture] = useState<string|null>(null);

    const resolution = appState.photobooth.webcam_resolution;

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

                    setLastPicture(url);

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

        if (lastMessage.type == 'TIMER') {
            setTimer(lastMessage.payload)
            if (lastMessage.payload === 0) {
                setFlash(true);
                setTimeout(() => {
                    takePicture(false);
                    setFlash(false);
                }, 150);
            }

            return
        }

        if (lastMessage.type == 'UNATTENDED_PICTURE') {
            takePicture(true);
        }
    }, [lastMessage]);

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
            lastPicture && <div className="picture_frame">
                <img src={lastPicture} alt="Last picture" />
            </div>
        }
    </div>
}