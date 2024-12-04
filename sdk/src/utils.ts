// Thx https://stackoverflow.com/questions/60188877/how-to-send-captured-webcam-image-and-save-to-server-via-input-field
export function b64ImageToBlob(b64data: string) {
    const data = b64data.split(';');
    const contentType = data[0].split(':')[1];
    b64data = data[1].split(',')[1];

    const bytesCharacter = atob(b64data);
    const byteArr: any[] = [];

    for (let offset = 0; offset < bytesCharacter.length; offset += 512) {
        const slice = bytesCharacter.slice(offset, offset + 512);
        const byteNumbers = new Array(slice.length);

        for (let i = 0; i < slice.length; i++) {
            byteNumbers[i] = slice.charCodeAt(i);
        }

        byteArr.push(new Uint8Array(byteNumbers));
    }

    return new Blob(byteArr, { type: contentType });
}
