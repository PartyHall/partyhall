export class Backdrop {
    id: number;
    albumId: number;
    nexusId: number;
    title: string;
    filename: string;

    constructor(data: Record<string, any>) {
        this.id = data['id'];
        this.albumId = data['album_id'];
        this.nexusId = data['nexus_id'];
        this.title = data['title'];
        this.filename = data['filename'];
    }

    static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new Backdrop(data);
    }

    static fromArray(data: Record<string, any>[] | null) {
        if (!data) {
            return [];
        }

        return data.map((x) => new Backdrop(x));
    }
}

export class BackdropAlbum {
    id: number;
    nexusId: number;
    name: string;
    author: string;
    version: number;
    backdrops: Backdrop[];

    constructor(data: Record<string, any>) {
        this.id = data['id'];
        this.nexusId = data['nexus_id'];
        this.name = data['name'];
        this.author = data['author'];
        this.version = data['version'];
        this.backdrops = Backdrop.fromArray(data['backdrops']);
    }

    static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new BackdropAlbum(data);
    }

    static fromArray(data: Record<string, any>[]) {
        return data.map((x) => new BackdropAlbum(x));
    }
}
