import { DetailedHTMLProps, ReactEventHandler, VideoHTMLAttributes, useEffect, useRef, useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import { Stack } from "@mui/material";

type Props = {
    isPlaying: boolean;
    song: KaraokeSong;
    onEnd(...args: unknown[]): unknown;
    onStatus(...args: unknown[]): unknown;

    volumeInstru: number;
    volumeVocals: number;
    volumeFull: number;
};

export default function VideoPlayer({isPlaying, song, onEnd, onStatus, volumeFull, volumeInstru, volumeVocals}: Props) {
    const videoRef = useRef<HTMLVideoElement>();
    const vocalsRef = useRef<HTMLAudioElement>();
    const fullRef = useRef<HTMLAudioElement>();

    const handleProgress = (e: any) => {
        if (isNaN(e.target.duration)){
            onStatus({
                position: 0,
                total: -1,
            });

            return;
        }

        onStatus({
            position: e.target.currentTime,
            total: Math.floor(e.target.duration),
        });

        // Thats ugly to do this there but heh
        if (!videoRef.current) {
            return
        }

        videoRef.current.volume = volumeInstru;

        if (vocalsRef.current) {
            vocalsRef.current.volume = volumeVocals;
        }

        if (fullRef.current) {
            fullRef.current.volume = volumeFull;
        }
    };

    useEffect(() => {
        if (!videoRef.current) {
            return;
        }

        if (isPlaying) {
            videoRef.current.play();
            if (vocalsRef.current && song.has_vocals) {
                vocalsRef.current.play();
            }
            if (fullRef.current && song.has_full && !song.has_vocals) {
                fullRef.current.play();
            }
        } else {
            videoRef.current.pause();
            if (vocalsRef.current && song.has_vocals) {
                vocalsRef.current.pause();
            }
            if (fullRef.current && song.has_vocals && !song.has_vocals) {
                fullRef.current.pause();
            }
        }
    }, [isPlaying]);

    return <>
        {
            song.has_vocals &&
            <audio
                preload='auto'
                // onCanPlayThrough={this.updateIsPlaying}
                src={'/api/modules/karaoke/song/' + song.uuid + '/vocals-mp3'}
                // Fuck typescript
                //@ts-ignore
                ref={vocalsRef}
            />
        }
        {
            song.has_full && !song.has_vocals &&
            <audio
                preload='auto'
                // onCanPlayThrough={this.updateIsPlaying}
                src={'/api/modules/karaoke/song/' + song.uuid + '/full-mp3'}
                controls
                // Fuck typescript
                //@ts-ignore
                ref={fullRef}
            />
        }
        <video 
            src={'/api/modules/karaoke/song/' + song.uuid + '/instrumental-webm'}
            onEnded={onEnd}
            // onCanPlayThrough={this.}
            controls
            // Fuck typescript
            //@ts-ignore
            ref={videoRef}
            onTimeUpdate={handleProgress}
        />
    </>;
}