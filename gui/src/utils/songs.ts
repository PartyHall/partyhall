export function songTitle(song: any) {
    let title = song.id;

    if (!!song.artist && !!song.title) {
        title = `${song.artist} - ${song.title}`;
    } else if (!!song.title) {
        title = song.title;
    }

    return title;
}