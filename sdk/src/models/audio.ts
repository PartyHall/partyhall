export class AudioDevice {
    id: number;
    name: string;
    description: string;
    volume: number;

    constructor(data: Record<string, any>) {
        this.id = data['id'];
        this.name = data['name'];
        this.description = data['description'];
        this.volume = data['volume']
    }

    static fromJson(data: Record<string, any>|null) {
        if (!data) {
            return null;
        }

        return new AudioDevice(data);
    }

    static fromArray(data: Record<string, any>[]) {
        return data.map(x => new AudioDevice(x));
    }
}

export class AudioDevices {
    karaokeSource: AudioDevice;
    karaokeSink: AudioDevice;
    defaultSource: AudioDevice|null;
    defaultSink: AudioDevice|null;
    sources: AudioDevice[];
    sinks: AudioDevice[];

    constructor(data: Record<string, any>) {
        this.karaokeSource = new AudioDevice(data['karaoke_source']);
        this.karaokeSink = new AudioDevice(data['karaoke_sink']);
        this.defaultSource = AudioDevice.fromJson(data['default_source']);
        this.defaultSink = AudioDevice.fromJson(data['default_sink']);
        this.sources = AudioDevice.fromArray(data['sources']);
        this.sinks = AudioDevice.fromArray(data['sinks']);
    }

    static fromJson(data: Record<string, any>|null) {
        if (!data) {
            return null;
        }

        return new AudioDevices(data);
    }
}