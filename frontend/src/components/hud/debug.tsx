import { ReactNode } from 'react';
import { useAuth } from '../../hooks/auth';

const D = (title: string, child: ReactNode) => (
    <div>
        <span style={{ fontWeight: 'bold' }}>{title}</span>: {child}
    </div>
);

export function DebugLeft() {
    const { currentMode, currentEvent, debug } = useAuth();

    if (!debug) {
        return <></>;
    }

    return (
        <div className="debug">
            {D('Current mode', <span>{currentMode}</span>)}
            {D('Event ID', `${currentEvent?.id}`)}
            {D(
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

    if (!debug) {
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
