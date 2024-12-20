import { SDK } from './index';
import { b64ImageToBlob } from './utils';

export default class Photobooth {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async takePicture() {
        await this.sdk.post('/api/webapp/picture');
    }

    public async uploadPicture(b64Image: string, unattended: boolean) {
        const form = new FormData();

        form.append('picture', b64ImageToBlob(b64Image));
        form.append('unattended', unattended ? 'true' : 'false');

        const resp = await fetch('/api/appliance/picture', {
            method: 'POST',
            body: form,
            headers: {
                Authorization: 'Bearer ' + this.sdk.token,
            },
        });

        return await resp.blob();
    }

    public async setFlash(powered: boolean, brightness: number) {
        await this.sdk.post(`/api/webapp/flash`, {
            powered: powered,
            brightness: brightness,
        });
    }
}
