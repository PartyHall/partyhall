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

const loadImage = (src: string): Promise<HTMLImageElement> => {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = reject;
        img.src = src;
    });
};

export default function DefaultView() {
    const {
        currentMode,
        user_settings,
        shouldTakePicture,
        setPictureTaken,
        api,
        backdropAlbum,
        selectedBackdrop
    } = useAuth();

    const photobooth = user_settings?.photobooth;

    const [countdown, setCountdown] = useState<number>(0);
    const [flash, setFlash] = useState<boolean>(false);
    const [lastPictureTaken, setLastPictureTaken] = useState<Blob | null>(null);

    const webcamRef = useRef<Webcam>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);

    useAsyncEffect(async () => {
        if (!shouldTakePicture || !webcamRef || !webcamRef.current || !canvasRef.current || !photobooth) {
            return;
        }

        if (shouldTakePicture == 'normal') {
            // Display the countdown
            for (let x = photobooth.countdown; x > 0; x--) {
                setCountdown(x);
                await timeout(1000);
            }

            setCountdown(0);

            await api.state.setFlash(true);
            setFlash(true);
            await timeout(1000); // We wait for the webcam to pick-up the flash
        }

        const picture = webcamRef.current.getScreenshot();
        let editedPicture = null;
        setFlash(false);

        if (shouldTakePicture === 'normal') {
            await api.state.setFlash(false);
        }

        if (picture) {
            if (shouldTakePicture === 'normal' && backdropAlbum && selectedBackdrop > 0) {
                const canvas = canvasRef.current;
                const camCanvas = webcamRef.current.getCanvas();

                canvas.width = camCanvas?.width ?? photobooth.resolution.width;
                canvas.height = camCanvas?.height ?? photobooth.resolution.height;

                const ctx = canvas.getContext('2d');
                if (ctx) {
                    ctx.clearRect(0, 0, canvas.width, canvas.height);

                    // First we draw the webcam image
                    const webcamImage = await loadImage(picture);
                    ctx.drawImage(webcamImage, 0, 0, canvas.width, canvas.height);

                    // Then we draw the backdrop
                    const backdropImage = await loadImage(
                        api.backdrops.getImageLink(backdropAlbum.backdrops[selectedBackdrop - 1].id)
                    );
                    ctx.drawImage(backdropImage, 0, 0, canvas.width, canvas.height);

                    editedPicture = canvas.toDataURL('image/jpg', 100);
                }
            }

            const resp = await api.photobooth.uploadPicture(
                editedPicture ?? picture,
                shouldTakePicture === 'unattended',
                editedPicture ? picture : null
            );

            if (shouldTakePicture === 'normal') {
                setLastPictureTaken(resp);
                await timeout(3500);
                setLastPictureTaken(null);
            }
        }

        setPictureTaken();
    }, [shouldTakePicture]);

    if (!photobooth) {
        return <div>loading</div>;
    }

    return (
        <>
            <canvas
                ref={canvasRef}
                // style={{ position: 'absolute', zIndex: 1000, top: 0, left: 0, width: '320px', height: '180px', background: 'red'}}
                style={{ display: 'none' }}
            />

            <Hud />
            <div id="webcam">
                <Webcam
                    ref={webcamRef}
                    width={photobooth.resolution.width}
                    height={photobooth.resolution.height}
                    screenshotFormat="image/jpeg"
                    videoConstraints={{
                        facingMode: 'user',
                        ...photobooth.resolution,
                    }}
                    mirrored={true}
                    forceScreenshotSourceSize
                />

                {backdropAlbum && selectedBackdrop > 0 && (
                    <img
                        id="webcam--backdrop"
                        src={api.backdrops.getImageLink(backdropAlbum.backdrops[selectedBackdrop - 1].id)}
                        alt="Backdrop"
                    />
                )}
            </div>

            <KaraokeHud />

            {currentMode === 'disabled' && <Disabled />}

            {countdown > 0 && <Countdown seconds={countdown} />}
            {flash && <Flash />}
            {lastPictureTaken && (
                <div className="lastPicture">
                    <img src={window.URL.createObjectURL(lastPictureTaken)} alt="Last picture" />
                </div>
            )}
        </>
    );
}
