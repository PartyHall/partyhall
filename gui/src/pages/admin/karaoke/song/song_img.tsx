import { useRef } from "react";
import { KaraokeSong } from "../../../../types/appstate";

const FALLBACK_IMG = '/api/modules/karaoke/fallback-image';

export default function SongCoverImage({song}: {song: KaraokeSong}) {
    const imgRef = useRef<HTMLImageElement>(null);

    let imgSource = `/api/modules/karaoke/song/${song.uuid}/cover`;
    if (!song.has_cover) {
        imgSource = FALLBACK_IMG;
    }

    const onImageError = () => {
        if (!imgRef.current) {
            return;
        }

        imgRef.current.src = FALLBACK_IMG;
    };

    return <img
        ref={imgRef}
        onError={onImageError}
        src={imgSource}
        alt={song.title}
        style={{ maxHeight: '6em', objectFit: 'contain' }}
    />;
}