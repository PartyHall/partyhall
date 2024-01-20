import { DetailedHTMLProps, ReactEventHandler, VideoHTMLAttributes, useEffect, useRef, useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import { Stack } from "@mui/material";

type Props = {
    isPlaying: boolean;
    song: KaraokeSong;
    onEnd(...args: unknown[]): unknown;
    onStatus(...args: unknown[]): unknown;
};

export default function VideoPlayer({isPlaying, song, onEnd, onStatus}: Props) {
    const videoRef = useRef<HTMLVideoElement>();

    const handleProgress = (e: any) => {
        if (isNaN(e.target.duration)){
            onStatus({
                position: 0,
                duration: 360,
            });

            return;
        }

        onStatus({
            position: e.target.currentTime,
            duration: e.target.duration,
        });
    };

    useEffect(() => {
        if (!videoRef.current) {
            return;
        }

        if (isPlaying) {
            videoRef.current.play();
        } else {
            videoRef.current.pause();
        }
    }, [isPlaying]);

    return <video 
        src={'/media/karaoke/' + song.filename + '/song.' + song.format}
        autoPlay
        onEnded={onEnd}
        controls
        // Fuck typescript
        //@ts-ignore
        ref={videoRef}
        onTimeUpdate={handleProgress}
    />;
}