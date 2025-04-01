export default class SpotifySettings {
    enabled: boolean;
    name: string;

    constructor(data: Record<string, any>) {
        this.enabled = data['enabled'];
        this.name = data['name'];
    }

    public static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new SpotifySettings(data);
    }
}
