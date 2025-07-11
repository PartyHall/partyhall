import { useAuth } from '../../hooks/auth';

export default function DisplayText() {
    const { currentEvent } = useAuth();

    if (!currentEvent || (!currentEvent.displayTextAppliance && !currentEvent.display_text_appliance)) {
        return null;
    }

    return (
        <div
            style={{
                display: 'flex',
                flexDirection: 'column',
                gap: '1em',
                alignItems: 'center',
                justifyItems: 'center',
                maxWidth: '200px',
            }}
        >
            <span style={{ textAlign: 'center' }}>{currentEvent.displayText ?? currentEvent.display_text}</span>
            {(currentEvent.registrationUrl || currentEvent.registration_url) && (
                <img
                    className="qr"
                    src={`/api/events/${currentEvent.id}/registration-qr`}
                    style={{ width: '150px', height: '150px' }}
                />
            )}
        </div>
    );
}
