import { Card } from 'antd';
import EventTitle from './event_title';
import KeyVal from './keyval';
import { PhEvent } from '@partyhall/sdk';
import { useTranslation } from 'react-i18next';

export default function EventCard({ event }: { event: PhEvent }) {
    const { t } = useTranslation();

    return (
        <Card title={<EventTitle event={event} />}>
            {event.date && <KeyVal label={t('generic.date')}>{event?.date.toFormat('yyyy/MM/dd - hh:mm')}</KeyVal>}
            {event.location && <KeyVal label={t('generic.location')}>{event.location}</KeyVal>}
        </Card>
    );
}
