import { Flex, Typography } from 'antd';

import EventForm from '../../components/event_form';
import Title from 'antd/es/typography/Title';

import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function NewEventPage() {
    const { t } = useTranslation('', { keyPrefix: 'events.editor' });
    const { setPageName } = useSettings();
    const navigate = useNavigate();

    useEffect(() => setPageName('events'), []);

    return (
        <Flex vertical gap={32}>
            <Typography>
                <Title>{t('create_title')}</Title>
            </Typography>

            <EventForm event={null} onSaved={() => navigate('/events')} />
        </Flex>
    );
}
