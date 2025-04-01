export default class NexusSettings {
    baseUrl: string;
    hardwareId: string;
    bypassSsl: boolean;
    errorMessage?: string | null;

    constructor(data: Record<string, any>) {
        this.baseUrl = data['nexus_url'];
        this.hardwareId = data['hardware_id'];
        this.bypassSsl = data['bypass_ssl'];
        this.errorMessage = data['error'];
    }

    public static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new NexusSettings(data);
    }
}
