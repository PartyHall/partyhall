import { useCallback, useEffect, useRef, useState } from 'react';
import CDGraphics from 'cdgraphics';
import { PhSongSession } from '@partyhall/sdk';
import { TimingObject } from 'timing-object';
import { setTimingsrc } from 'timingsrc';

const BACKDROP_PADDING = 16;
const BORDER_RADIUS = 16;

// Stolen & adapted from
// https://github.com/bhj/KaraokeEternal/blob/main/src/routes/Player/components/Player/CDGPlayer/CDGPlayer.tsx

// Converted to functional component by Claude (recipe for disaster but meh it seems to work)

interface CDGPlayerProps {
    cdgAlpha: number;
    cdgSize: number;
    isPlaying: boolean;
    volumeInstru: number;
    volumeVocals: number;
    session: PhSongSession;
    width: number;
    height: number;
    onEnd(...args: unknown[]): unknown;
    onProgress(timecode: number): Promise<void>;
}

export default function CDGPlayer({
    cdgAlpha,
    cdgSize,
    isPlaying,
    volumeInstru,
    volumeVocals,
    session,
    width,
    height,
    onEnd,
    onProgress,
}: CDGPlayerProps) {
    // Refs
    const audioRef = useRef<HTMLAudioElement>(null);
    // const vocalsRef = useRef<HTMLAudioElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const canvasCtxRef = useRef<CanvasRenderingContext2D | null>(null);
    const cdgRef = useRef<CDGraphics | null>(null);
    const frameIdRef = useRef<number | null>(null);
    const lastBitmapRef = useRef<ImageBitmap | null>(null);
    const timingObjectRef = useRef<any | null>(null);
    const cleanupFunctionsRef = useRef<(() => void)[]>([]);

    // State
    const [backgroundRGBA, setBackgroundRGBA] = useState<number[]>([0, 0, 0, 0]);
    const [contentBounds, setContentBounds] = useState<number[]>([0, 0, 0, 0]);

    // Setup initial Canvas and CDG
    useEffect(() => {
        if (!canvasRef.current) return;

        canvasCtxRef.current = canvasRef.current.getContext('2d');
        cdgRef.current = new CDGraphics();
        timingObjectRef.current = new TimingObject();

        // Setup timing sync
        if (audioRef.current) {
            const cleanup = setTimingsrc(audioRef.current, timingObjectRef.current);
            cleanupFunctionsRef.current.push(cleanup);
        }

        /*
        if (vocalsRef.current) {
            const cleanup = setTimingsrc(vocalsRef.current, timingObjectRef.current);
            cleanupFunctionsRef.current.push(cleanup);
        }
            */

        updateSources();

        return () => {
            stopCDG();
            cleanupFunctionsRef.current.forEach((cleanup) => cleanup());
        };
    }, []);

    useEffect(() => {
        updateSources();
    }, [session.id]);

    useEffect(() => {
        if (!timingObjectRef.current) {
            return;
        }

        if (isPlaying) {
            timingObjectRef.current.update({ velocity: 1.0 });
            startCDG();
        } else {
            timingObjectRef.current.update({ velocity: 0.0 });
            stopCDG();
        }
    }, [isPlaying]);

    useEffect(() => {
        if (!audioRef.current) {
            return;
        }

        audioRef.current.volume = volumeInstru / 100;
        /*
        if (vocalsRef.current) {
            vocalsRef.current.volume = volumeVocals / 100;
        }
            */
    }, [volumeInstru, volumeVocals]);

    useEffect(() => {
        if (lastBitmapRef.current) {
            paintCDG(lastBitmapRef.current);
        }
    }, [width, height, cdgSize]);

    const updateSources = useCallback(async () => {
        stopCDG();

        try {
            const resp = await fetch(session.song.getLyricsUrl());
            const buf = await resp.arrayBuffer();

            if (!audioRef.current || !cdgRef.current) return;

            cdgRef.current.load(buf);
            audioRef.current.src = session.song.getInstrumentalUrl();
            audioRef.current.load();

            /*
            if (vocalsRef.current) {
                vocalsRef.current.src = session.song.getVocalsUrl();
                vocalsRef.current.load();
            }
                */
        } catch (err) {
            console.error(err);
        }
    }, [session.song]);

    const paintCDG = useCallback(
        (bitmap: ImageBitmap) => {
            if (!canvasRef.current || !canvasCtxRef.current) return;

            const { clientWidth, clientHeight } = canvasRef.current;

            canvasCtxRef.current.imageSmoothingEnabled = false;
            canvasCtxRef.current.shadowBlur = Math.min(16, clientHeight * cdgSize * 0.0333);
            canvasCtxRef.current.shadowColor = 'rgba(0,0,0,1)';
            canvasCtxRef.current.clearRect(0, 0, clientWidth, clientHeight);
            canvasCtxRef.current.drawImage(bitmap, 0, 0, clientWidth, clientHeight);
        },
        [cdgSize]
    );

    const startCDG = useCallback(async () => {
        if (!cdgRef.current || !audioRef.current) return;

        frameIdRef.current = requestAnimationFrame(startCDG);

        const currentTime = timingObjectRef.current
            ? timingObjectRef.current.query().position
            : audioRef.current.currentTime;

        const frame = cdgRef.current.render(currentTime, { forceKey: true });
        if (!frame.isChanged) return;

        if (!frame.backgroundRGBA.every((val, i) => val === backgroundRGBA[i])) {
            setBackgroundRGBA(frame.backgroundRGBA);
        }

        if (!frame.contentBounds.every((val, i) => val === contentBounds[i])) {
            setContentBounds(frame.contentBounds);
        }

        try {
            const bitmap = await createImageBitmap(frame.imageData);
            lastBitmapRef.current = bitmap;
            paintCDG(bitmap);
        } catch (err) {
            console.error(err);
        }
    }, [backgroundRGBA, contentBounds]);

    const stopCDG = useCallback(() => {
        if (frameIdRef.current) {
            cancelAnimationFrame(frameIdRef.current);
        }
    }, []);

    const handleEnded = useCallback(() => {
        if (timingObjectRef.current) {
            timingObjectRef.current.update({ velocity: 0.0 });
        }
        onEnd();
        stopCDG();
    }, [onEnd]);

    const handleTimeUpdate = useCallback(async () => {
        if (!audioRef.current) return;
        await onProgress(Math.floor(audioRef.current.currentTime));
    }, [onProgress]);

    // Calculate canvas dimensions
    const [x1, y1, x2, y2] = contentBounds;
    const [r, g, b] = backgroundRGBA;
    const wScale = (width - BACKDROP_PADDING * 2) / 300;
    const hScale = ((height - BACKDROP_PADDING * 2) * cdgSize) / 216;
    const scale = Math.min(wScale, hScale);
    const pad = x2 - x1 && y2 - y1 ? BACKDROP_PADDING : 0;

    const filters = [
        `blur(${30 * cdgAlpha * scale}px)`,
        `brightness(${100 - 100 * cdgAlpha ** 3}%)`,
        `saturate(${100 - 100 * cdgAlpha ** 3}%)`,
    ];

    return (
        <div className="karaoke__inner">
            <div
                className="karaoke__backdrop"
                style={{
                    backdropFilter: cdgAlpha !== 1 ? filters.join(' ') : 'none',
                    backgroundColor: cdgAlpha !== 1 ? 'transparent' : `rgba(${r},${g},${b},${cdgAlpha})`,
                    borderRadius: BORDER_RADIUS * scale,
                    left: (x1 - pad) * scale,
                    top: (y1 - pad) * scale,
                    width: (x2 - x1 + pad * 2) * scale,
                    height: (y2 - y1 + pad * 2) * scale,
                }}
            />
            <canvas ref={canvasRef} width={300 * scale} height={216 * scale} className="karaoke__canvas" />
            <br />
            <audio preload="auto" onEnded={handleEnded} onTimeUpdate={handleTimeUpdate} ref={audioRef} />
            {/*session.song.has_vocals && <audio preload="auto" ref={vocalsRef} />*/}
        </div>
    );
}
