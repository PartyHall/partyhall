export type COVER_SOURCE = 'NO_COVER'|'LINK'|'UPLOADED';
export type FORMAT = 'CDG'|'WEBM'|'MP4';

export type SongPost = {
    title: string;
    artist: string;
    format: string;
    cover_type: COVER_SOURCE;
    cover_data: string;

    song: FileList;
    cdg: FileList;
};