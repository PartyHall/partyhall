import React, { useRef } from 'react';
import CDGraphics from 'cdgraphics';
import Webcam from 'react-webcam';
import { useBoothSocket } from '../../../hooks/boothSocket';
import { KaraokeSong } from '../../../types/appstate';

const BACKDROP_PADDING = 10 // px at 1:1 scale
const BORDER_RADIUS = parseInt(getComputedStyle(document.body).getPropertyValue('--border-radius'))

// Stolen & adapted from
// https://github.com/bhj/KaraokeEternal/blob/main/src/routes/Player/components/Player/CDGPlayer/CDGPlayer.tsx

// I failed to write it myself in a FC so when I found out that the dev of the lib I wanted to use
// had a full karaoke software in React I took a look in it

interface CDGPlayerProps {
    cdgAlpha: number;
    cdgSize: number;
    isPlaying: boolean;
    volumeInstru: number;
    volumeVocals: number;
    volumeFull: number;
    song: KaraokeSong;
    width: number;
    height: number;
    onEnd(...args: unknown[]): unknown;
    onError(...args: unknown[]): unknown;
    onLoad(...args: unknown[]): unknown;
    onPlay(...args: unknown[]): unknown;
    onStatus(...args: unknown[]): unknown;
}

class CDGPlayer extends React.Component<CDGPlayerProps> {
    audio = React.createRef<HTMLAudioElement>();
    vocals = React.createRef<HTMLAudioElement>();
    full = React.createRef<HTMLAudioElement>();

    canvas = React.createRef<HTMLCanvasElement>();
    canvasCtx: any = null;
    cdg = null;
    frameId: number|null = null;
    lastBitmap: ImageBitmap|null = null;

    state = {
        backgroundRGBA: [0, 0, 0, 0],
        contentBounds: [0, 0, 0, 0], // x1, y1, x2, y2
    };

    componentDidMount() {
        if (!this.canvas.current) {
            return;
        }

        this.canvasCtx = this.canvas.current.getContext('2d');
        this.cdg = new CDGraphics();

        this.updateSources();
    }

    componentDidUpdate(prevProps: CDGPlayerProps) {
        if (prevProps.song.id !== this.props.song.id) {
            this.updateSources();
            return;
        }

        if (prevProps.isPlaying !== this.props.isPlaying) {
            this.updateIsPlaying();
        }

        if (prevProps.width !== this.props.width || prevProps.height !== this.props.height || prevProps.cdgSize !== this.props.cdgSize) {
            if (this.lastBitmap){
                this.paintCDG(this.lastBitmap);
            }
        }
    }

    componentWillUnmount() {
        this.stopCDG();
    }

    render() {
        const { cdgAlpha, cdgSize, width, height } = this.props;
        const [x1, y1, x2, y2] = this.state.contentBounds;
        const [r, g, b] = this.state.backgroundRGBA;

        // apply sizing as % of max height, leaving room for the backdrop
        const wScale = (width - (BACKDROP_PADDING * 2)) / 300;
        const hScale = ((height - (BACKDROP_PADDING * 2)) * cdgSize) / 216;
        const scale = Math.min(wScale, hScale);
        const pad = (x2 - x1) && (y2 - y1) ? BACKDROP_PADDING : 0;

        const filters = [
            `blur(${30 * cdgAlpha * scale}px)`,
            `brightness(${100 - (100 * (cdgAlpha ** 3))}%)`,
            `saturate(${100 - (100 * (cdgAlpha ** 3))}%)`,
        ];

        return (
            <div className='karaoke__inner'> 
                <div className='karaoke__backdrop' style={{
                    backdropFilter: cdgAlpha !== 1 ? filters.join(' ') : 'none',
                    backgroundColor: cdgAlpha !== 1 ? 'transparent' : `rgba(${r},${g},${b},${cdgAlpha})`,
                    borderRadius: BORDER_RADIUS * scale,
                    left: (x1 - pad) * scale,
                    top: (y1 - pad) * scale,
                    width: ((x2 - x1) + pad * 2) * scale,
                    height: ((y2 - y1) + pad * 2) * scale,
                }}></div>
                <canvas
                    ref={this.canvas}
                    width={300 * scale}
                    height={216 * scale}
                    className='karaoke__canvas'
                />
                <br />
                <audio
                    preload='auto'
                    onCanPlayThrough={this.updateIsPlaying}
                    onEnded={this.handleEnded}
                    onError={this.handleError}
                    onLoadStart={this.props.onLoad}
                    onPlay={this.handlePlay}
                    onTimeUpdate={this.handleTimeUpdate}
                    ref={this.audio}
                />

                {
                    this.props.song.has_vocals &&
                    <audio
                        preload='auto'
                        onCanPlayThrough={this.updateIsPlaying}
                        ref={this.vocals}
                    />
                }
                {
                    !this.props.song.has_vocals && this.props.song.has_full &&
                    <audio
                        preload='auto'
                        onCanPlayThrough={this.updateIsPlaying}
                        ref={this.full}
                    />
                }
            </div>
        )
    }

    updateSources = async () => {
        this.stopCDG();

        try {
            const resp = await fetch('/api/modules/karaoke/song/' + this.props.song.uuid + '/cdg');
            const buf = await resp.arrayBuffer();
    
            if (!this.audio.current || !this.cdg) {
                return;
            }
    
            //@ts-ignore
            this.cdg.load(buf);
            this.audio.current.src = '/api/modules/karaoke/song/' + this.props.song.uuid + '/instrumental-mp3';
            this.audio.current.load();

            // HTML Audio are not perfectly in sync this could lead to issues
            // If this happens, need to implement this
            // https://stackoverflow.com/questions/54509959/how-do-i-play-audio-files-synchronously-in-javascript
            // After some test it looks like YEAH THIS WILL BE REQUIRED
            // No clue how we'll do this in video mode...
            if (this.vocals.current && this.props.song.has_vocals) {
                this.vocals.current.src = '/api/modules/karaoke/song/' + this.props.song.uuid + '/vocals-mp3';
                this.vocals.current?.load();
            }

            if (this.full.current && !this.props.song.has_vocals && this.props.song.has_full) {
                this.full.current.src = '/api/modules/karaoke/song/' + this.props.song.uuid + '/full-mp3';
                this.full.current?.load();
            }
        } catch (err: any) {
            this.props.onError(err.message);
        }
    }

    updateIsPlaying = () => {
        if (this.props.isPlaying) {
            if (this.audio.current) {
                try {
                    this.audio.current.play();

                    if (this.vocals.current) {
                        this.vocals.current.play();
                    }

                    if (this.full.current) {
                        this.full.current.play();
                    }
                } catch (err: any) {
                    this.props.onError(err.message);
                }
            }
        } else {
            if (this.audio.current) {
                try {
                    this.audio.current.pause();

                    if (this.vocals.current) {
                        this.vocals.current.pause();
                    }

                    if (this.full.current) {
                        this.full.current.pause();
                    }
                } catch (err: any) {
                    this.props.onError(err.message);
                }
            }

            this.stopCDG();
        }
    }

    handleEnded = () => {
        this.props.onEnd();
        this.stopCDG();
    }

    //@ts-ignore
    handleError = (el) => {
        const { message, code } = el.target.error;
        this.props.onError(`${message} (code ${code})`);
    }

    handlePlay = () => {
        this.props.onPlay();
        this.startCDG();
    }

    handleTimeUpdate = () => {
        if (!this.audio.current) {
            return;
        }
        
        // Thats ugly to do this there but heh
        this.audio.current.volume = this.props.volumeInstru;

        if (this.vocals.current) {
            this.vocals.current.volume = this.props.volumeVocals;
        }

        if (this.full.current) {
            this.full.current.volume = this.props.volumeFull;
        }

        this.props.onStatus({
            position: this.audio.current.currentTime,
            total: this.audio.current.duration,
        });
    }

    paintCDG = (bitmap: ImageBitmap) => {
        if (!this.canvas.current || !this.canvasCtx) {
            return;
        }

        const { clientWidth, clientHeight } = this.canvas.current;

        this.canvasCtx.imageSmoothingEnabled = false;
        this.canvasCtx.shadowBlur = Math.min(16, clientHeight * this.props.cdgSize * 0.0333);
        this.canvasCtx.shadowColor = 'rgba(0,0,0,1)';
        this.canvasCtx.clearRect(0, 0, clientWidth, clientHeight);
        this.canvasCtx.drawImage(bitmap, 0, 0, clientWidth, clientHeight);
    }

    startCDG = async () => {
        if (!this.cdg || !this.audio.current) {
            return;
        }

        this.frameId = requestAnimationFrame(this.startCDG)

        //@ts-ignore
        const frame = this.cdg.render(this.audio.current.currentTime, { forceKey: true })
        if (!frame.isChanged){
            return;
        }

        if (!frame.backgroundRGBA.every((val: any, i: number) => val === this.state.backgroundRGBA[i])) {
            this.setState({ backgroundRGBA: frame.backgroundRGBA })
        }

        if (!frame.contentBounds.every((val: any, i: number) => val === this.state.contentBounds[i])) {
            this.setState({ contentBounds: frame.contentBounds })
        }

        try {
            const bitmap = await createImageBitmap(frame.imageData);

            this.lastBitmap = bitmap; // cache for re-painting if canvas size changes
            this.paintCDG(bitmap);
        } catch (err: any){
            this.props.onError(err.message);
        }
    }

    stopCDG = () => {
        if (this.frameId) {
            cancelAnimationFrame(this.frameId);
        }
    }
}

export default CDGPlayer