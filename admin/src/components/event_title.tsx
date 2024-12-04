import { PhEvent } from '@partyhall/sdk';
import { useTranslation } from 'react-i18next';

export default function EventTitle({ event }: { event: PhEvent }) {
    const { t } = useTranslation();

    if (!event.author || event.author.trim().length === 0) {
        return <>{event.name}</>;
    }

    return (
        <>
            {t('generic.event_title', {
                eventName: event.name,
                eventAuthor: event.author,
            })}
        </>
    );
}
