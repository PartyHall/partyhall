import { KaraokeSong } from "../types/appstate";

export function songTitle(song: KaraokeSong) {
    let title = `${song.id}`;

    if (!!song.artist && !!song.title) {
        title = `${song.artist} - ${song.title}`;
    } else if (!!song.title) {
        title = song.title;
    }

    return title;
}