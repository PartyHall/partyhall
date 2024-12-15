import { useEffect, useRef } from 'react';
import { PhSongSession } from '@partyhall/sdk';
import { TimingObject } from 'timing-object';
import { setTimingsrc } from 'timingsrc';

type Props = {
    session: PhSongSession;
    isPlaying: boolean;
    volumeInstru: number;
    volumeVocals: number;
    onProgress: (timecode: number) => Promise<void>;
    onEnded: () => Promise<void>;
};

/**
 * Meh, sync is ugly but we'll see if it works
 * i dont like this
 */

export default function VideoPlayer({ session, isPlaying, volumeInstru, volumeVocals, onProgress, onEnded }: Props) {
    const videoRef = useRef<HTMLVideoElement | null>(null);
    const vocalsRef = useRef<HTMLAudioElement | null>(null);
    const timingObjectRef = useRef<any | null>(null);
    const cleanupRef = useRef<{ video?: () => void; vocals?: () => void }>({});

    useEffect(() => {
        timingObjectRef.current = new TimingObject();

        return () => {
            if (cleanupRef.current.video) cleanupRef.current.video();
            if (cleanupRef.current.vocals) cleanupRef.current.vocals();
        };
    }, []);

    useEffect(() => {
        if (!videoRef.current || !timingObjectRef.current) {
            return;
        }

        cleanupRef.current.video = setTimingsrc(videoRef.current, timingObjectRef.current);

        if (vocalsRef.current) {
            cleanupRef.current.vocals = setTimingsrc(vocalsRef.current, timingObjectRef.current);
        }

        return () => {
            if (cleanupRef.current.video) {
                cleanupRef.current.video();
            }

            if (cleanupRef.current.vocals) {
                cleanupRef.current.vocals();
            }
        };
    }, []);

    useEffect(() => {
        if (!timingObjectRef.current) {
            return;
        }

        if (isPlaying) {
            timingObjectRef.current.update({ velocity: 1.0 });
        } else {
            timingObjectRef.current.update({ velocity: 0.0 });
        }
    }, [isPlaying]);

    useEffect(() => {
        if (!videoRef.current) {
            return;
        }

        videoRef.current.volume = volumeInstru / 100;
        if (vocalsRef.current) {
            vocalsRef.current.volume = volumeVocals / 100;
        }
    }, [volumeInstru, volumeVocals]);

    return (
        <>
            <video
                src={session.song.getInstrumentalUrl()}
                onEnded={onEnded}
                ref={videoRef}
                controls
                onTimeUpdate={async (x: React.SyntheticEvent<HTMLVideoElement>) => {
                    const video = x.target as HTMLVideoElement;
                    await onProgress(Math.floor(video.currentTime));
                }}
            />
            {session.song.has_vocals && <audio src={session.song.getVocalsUrl()} ref={vocalsRef} />}
        </>
    );
}
