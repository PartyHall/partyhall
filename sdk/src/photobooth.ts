import { SDK } from './index';
import { b64ImageToBlob } from './utils';

export default class Photobooth {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async takePicture() {
        await this.sdk.post('/api/photobooth/take-picture');
    }

    public async uploadPicture(b64Image: string, unattended: boolean, b64AlternateImage: string | null) {
        const form = new FormData();

        form.append('picture', b64ImageToBlob(b64Image));

        if (b64AlternateImage) {
            form.append('alternate_picture', b64ImageToBlob(b64AlternateImage));
        }

        form.append('unattended', unattended ? 'true' : 'false');

        const resp = await fetch('/api/photobooth/upload-picture', {
            method: 'POST',
            body: form,
            headers: {
                Authorization: 'Bearer ' + this.sdk.token,
            },
        });

        return await resp.blob();
    }
}
