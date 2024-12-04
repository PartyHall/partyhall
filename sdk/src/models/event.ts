import { DateTime } from 'luxon';

export class PhEvent {
    id: number | null;
    name: string;
    author: string | null;
    date: DateTime;
    location: string | null;
    nexusId: string | null;

    amtImagesHandtaken: number = 0;
    amtImagesUnattended: number = 0;

    public constructor(
        id: number | null,
        name: string,
        author: string | null,
        date: DateTime,
        location: string | null,
        nexusId: string | null
    ) {
        this.id = id;
        this.name = name;
        this.author = author;
        this.date = date;
        this.location = location;
        this.nexusId = nexusId;
    }

    static fromJson(json: Record<string, any> | null): PhEvent | null {
        if (!json) {
            return null;
        }

        const phEvent = new PhEvent(
            json.id,
            json.name,
            json.author,
            DateTime.fromISO(json.date),
            json.location,
            json.nexus_id
        );

        phEvent.amtImagesHandtaken = json['amt_images_handtaken'];
        phEvent.amtImagesUnattended = json['amt_images_unattended'];

        return phEvent;
    }

    asJson(): Record<string, any> {
        return {
            id: this.id,
            name: this.name,
            author: this.author,
            date: this.date.toISO(),
            location: this.location,
            nexus_id: this.nexusId,
        };
    }
}
