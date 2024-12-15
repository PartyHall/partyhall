import '../../assets/karaoke.scss';
import CDGPlayer from '../karaoke_player/cdg_player';
import VideoPlayer from '../karaoke_player/video_player';
import { useAuth } from '../../hooks/auth';

/**
 * @TODO:
 * Karaoke countdown should be made on the frontend, doing the loop on
 * the backend was a mistake I need to fix that
 */

export default function KaraokeHud() {
    const { karaoke, karaokeQueue, api } = useAuth();

    const onProgress = async (timecode: number) => {
        if (!karaoke.current) {
            return;
        }

        await api.karaoke.songProgress(timecode);
    };

    const onEnded = async () => {
        if (!karaoke.current) {
            return;
        }

        await api.karaoke.songEnded(karaoke.current.id);
    };

    return (
        <div id="karaoke">
            {karaoke.current && !karaoke.isPlaying && karaoke.countdown > 0 && (
                <div className="titlescreen karaokeBox">
                    <span className="blue-glow title">{karaoke.current.title}</span>
                    <span className="blue-glow artist">{karaoke.current.artist}</span>
                    <span className="countdown">{karaoke.countdown}</span>
                    <span className="singer">Singer: {karaoke.current.sung_by}</span>
                </div>
            )}

            {karaoke.current && karaoke.countdown === 0 && (
                <>
                    {karaoke.current.song.format.toLowerCase() !== 'cdg' && (
                        <VideoPlayer
                            session={karaoke.current}
                            isPlaying={karaoke.isPlaying}
                            volumeInstru={karaoke.volume}
                            volumeVocals={karaoke.volumeVocals}
                            onProgress={onProgress}
                            onEnded={onEnded}
                        />
                    )}
                    {karaoke.current.song.format.toLowerCase() === 'cdg' && (
                        <CDGPlayer
                            cdgAlpha={0.2}
                            cdgSize={window.innerHeight / 2}
                            isPlaying={karaoke.isPlaying}
                            volumeInstru={karaoke.volume}
                            volumeVocals={karaoke.volumeVocals}
                            session={karaoke.current}
                            width={window.innerWidth / 2}
                            height={window.innerHeight / 2}
                            onEnd={onEnded}
                            onProgress={onProgress}
                        />
                    )}
                </>
            )}

            {karaokeQueue.length > 0 && (
                <div className="nextup karaokeBox">
                    <span className="blue-glow">Next up:</span>
                    <span className="red-glow">{karaokeQueue[0].title}</span>
                    <span className="red-glow">{karaokeQueue[0].artist}</span>
                    <span className="red-glow">{karaokeQueue[0].sung_by}</span>
                </div>
            )}
        </div>
    );
}
