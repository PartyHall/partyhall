import { Flex, Typography } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';

import EventForm from '../../components/event_form';
import Loader from '../../components/loader';
import { PhEvent } from '@partyhall/sdk';
import Title from 'antd/es/typography/Title';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useSettings } from '../../hooks/settings';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

type Data = {
    loading: boolean;
    event: PhEvent | null;
};

export default function EditEvent() {
    const { t } = useTranslation('', { keyPrefix: 'events.editor' });
    const { setPageName } = useSettings();

    const { id } = useParams();
    const { api } = useAuth();
    const navigate = useNavigate();

    const [data, setData] = useState<Data>({
        loading: true,
        event: null,
    });

    useAsyncEffect(async () => {
        setPageName('events');

        if (!api.tokenUser?.hasRole('ADMIN')) {
            // @TODO: Make a 403 error screen with react router
            navigate('/');
            return;
        }

        if (!id) {
            // @TODO: Make a 400 error screen with react router
            navigate('/');
            return;
        }

        const event = await api.events.get(id);
        setData((oldCtx) => ({
            ...oldCtx,
            loading: false,
            event: event,
        }));
    }, [id]);

    return (
        <>
            <Loader loading={data.loading}>
                <Flex vertical gap={32}>
                    <Typography>
                        <Title>
                            {t('edit_title', { name: data.event?.name })}
                        </Title>
                    </Typography>

                    <EventForm event={data.event ?? null} />
                </Flex>
            </Loader>
        </>
    );
}
