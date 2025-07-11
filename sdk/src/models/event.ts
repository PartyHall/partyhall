import { DateTime } from 'luxon';

export class PhEvent {
    id: number | null;
    name: string;
    author: string | null;
    date: DateTime;
    location: string | null;
    nexusId: string | null;
    registrationUrl: string | null;
    displayText: string | null = null;
    displayTextAppliance: boolean = false;

    // This fucking frontend sucks so much
    // that it isn't parsed the same way in mercure & api
    // so fuck that sshit i dont care i need it working for tomorrow
    // this will stay like this until the rewrite
    display_text: string | null = null;
    display_text_appliance: boolean = false;
    registration_url: string | null = null;

    amtImagesHandtaken: number = 0;
    amtImagesUnattended: number = 0;

    public constructor(data: Record<string, any>) {
        this.id = data.id;
        this.name = data.name;
        this.author = data.author;
        this.date = DateTime.fromISO(data.date);
        this.location = data.location;
        this.nexusId = data.nexus_id;
        this.registrationUrl = data.registration_url;
        this.displayText = data.display_text;
        this.displayTextAppliance = data.display_text_appliance;
        this.amtImagesHandtaken = data.amt_images_handtaken;
        this.amtImagesUnattended = data.amt_images_unattended;
    }

    static fromJson(json: Record<string, any> | null): PhEvent | null {
        if (!json) {
            return null;
        }

        return new PhEvent(json);
    }

    asJson(): Record<string, any> {
        return {
            id: this.id,
            name: this.name,
            author: this.author,
            date: this.date.toISO(),
            location: this.location,
            nexus_id: this.nexusId,
            registration_url: this.registrationUrl,
            display_text: this.displayText,
            display_text_appliance: this.displayTextAppliance,
        };
    }
}
