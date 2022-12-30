export default function getSocketMode(): string {
    //@ts-ignore
    const injectedMode: string|null|undefined = window.SOCKET_TYPE;

    return injectedMode ?? 'booth';
}