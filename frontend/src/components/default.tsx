import { useRef, useState } from 'react';
import Countdown from './hud/countdown';
import Disabled from './hud/disabled';
import Flash from './hud/flash';
import Hud from './hud';
import KaraokeHud from './hud/karaoke';
import Webcam from 'react-webcam';

import { timeout } from '../utils';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../hooks/auth';

export default function DefaultView() {
    const {
        currentMode,
        modulesSettings,
        shouldTakePicture,
        setPictureTaken,
        api,
    } = useAuth();

    const [countdown, setCountdown] = useState<number>(0);
    const [flash, setFlash] = useState<boolean>(false);
    const [lastPictureTaken, setLastPictureTaken] = useState<Blob | null>(null);

    const webcamRef = useRef<Webcam>(null);

    useAsyncEffect(async () => {
        if (!shouldTakePicture || !webcamRef || !webcamRef.current) {
            return;
        }

        if (shouldTakePicture == 'normal') {
            // Display the countdown
            for (let x = modulesSettings.photobooth.countdown; x > 0; x--) {
                setCountdown(x);
                await timeout(1000);
            }

            setCountdown(0);

            setFlash(true);
            await timeout(1000); // We wait for the webcam to pick-up the flash
        }
        const picture = webcamRef.current.getScreenshot();
        setFlash(false);

        if (picture) {
            const resp = await api.photobooth.uploadPicture(
                picture,
                shouldTakePicture === 'unattended'
            );

            if (shouldTakePicture === 'normal') {
                setLastPictureTaken(resp);
                await timeout(3500);
                setLastPictureTaken(null);
            }
        }

        setPictureTaken();
    }, [shouldTakePicture]);

    return (
        <div>
            <Hud />
            <Webcam
                id="webcam"
                ref={webcamRef}
                width={modulesSettings.photobooth.resolution.width}
                height={modulesSettings.photobooth.resolution.height}
                screenshotFormat="image/jpeg"
                videoConstraints={{
                    facingMode: 'user',
                    ...modulesSettings.photobooth.resolution,
                }}
                forceScreenshotSourceSize
            />
            <KaraokeHud />

            {currentMode === 'disabled' && <Disabled />}

            {countdown > 0 && <Countdown seconds={countdown} />}
            {flash && <Flash />}
            {lastPictureTaken && (
                <div className="lastPicture">
                    <img
                        src={window.URL.createObjectURL(lastPictureTaken)}
                        alt="Last picture"
                    />
                </div>
            )}
        </div>
    );
}
