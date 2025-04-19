export class ApSettings {
    wiredInterface: string;
    wirelessInterface: string;

    enabled: boolean;
    ssid: string;
    password: string;

    constructor(data: Record<string, any>) {
        this.enabled = data['enabled'];
        this.ssid = data['ssid'];
        this.password = data['password'];

        this.wiredInterface = data['wired_interface'];
        this.wirelessInterface = data['wireless_interface'];
    }

    public static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new ApSettings(data);
    }
}

export class Interface {
    friendlyName: string;
    name: string;
    wireless: boolean;
    ips: string[];

    constructor(data: Record<string, any>) {
        this.friendlyName = data['friendly_name'];
        this.name = data['name'];
        this.wireless = data['wireless'];
        this.ips = data['ips'];
    }

    public static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new Interface(data);
    }
}

export class InterfacesSettings {
    accessPointSettings: ApSettings;
    wiredInterfaces: Interface[];
    wirelessInterfaces: Interface[];

    constructor(data: Record<string, any>) {
        if (!data['ap_settings']) {
            throw 'NO AP SETTINGS!';
        }

        this.accessPointSettings = new ApSettings(data['ap_settings']);

        if (!data['interfaces']) {
            throw 'NO INTERFACES';
        }

        this.wiredInterfaces = data['interfaces']['ethernet'].map((x: Record<string, any>) => new Interface(x));
        this.wirelessInterfaces = data['interfaces']['wifi'].map((x: Record<string, any>) => new Interface(x));
    }

    public withApSettings(data: Record<string, any>) {
        this.accessPointSettings = new ApSettings(data);

        return this;
    }

    public static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new InterfacesSettings(data);
    }
}
