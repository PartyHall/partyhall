export default function getSocketMode(): string {
    //@ts-ignore
    const injectedMode: string|null|undefined = window.SOCKET_TYPE;
    const query = new URLSearchParams(window.location.search);

    return query.get('mode') ?? injectedMode ?? 'booth';
}