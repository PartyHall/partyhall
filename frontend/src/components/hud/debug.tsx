import { ReactNode } from 'react';
import { useAuth } from '../../hooks/auth';

export const FORCE_DEBUG = true;

const D = (title: string, child: ReactNode) => (
    <div>
        <span style={{ fontWeight: 'bold' }}>{title}</span>: {child}
    </div>
);

export function DebugLeft() {
    const { currentMode, currentEvent, debug, backdropAlbum } = useAuth();

    if (!debug && !FORCE_DEBUG) {
        return <></>;
    }

    return (
        <div className="debug">
            {D('Current mode', <span>{currentMode}</span>)}
            {D('Event ID', `${currentEvent?.id}`)}
            {D('Backdrop Album', !backdropAlbum ? '' : `${backdropAlbum.name} (${backdropAlbum.id})`)}
            {debug &&
                D(
                    'IPs',
                    <ul>
                        {Object.entries(debug.ip_addresses)
                            .filter(([, x]) => x.length > 0)
                            .map(([key, inter]) => (
                                <li key={key}>
                                    {key}: {inter.join(', ')}
                                </li>
                            ))}
                    </ul>
                )}
        </div>
    );
}

export function DebugRight() {
    const { hwid, version, commit, debug } = useAuth();

    if (!debug && !FORCE_DEBUG) {
        return <></>;
    }

    return (
        <div className="debug">
            {D('HWID', hwid)}
            {D('Version', version)}
            {D('Commit', commit)}
        </div>
    );
}
