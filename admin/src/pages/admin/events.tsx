import { Button, Flex, Popconfirm, Table, Typography, notification } from 'antd';
import { Collection, PhEvent, SdkError } from '@partyhall/sdk';
import { Link, useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

import Loader from '../../components/loader';
import Title from 'antd/es/typography/Title';
import useAsyncEffect from 'use-async-effect';

import { useAuth } from '../../hooks/auth';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

type Status = {
    loading: boolean;
    page: number;
    resp: Collection<PhEvent> | null;
};

export default function Events() {
    const { t } = useTranslation('', { keyPrefix: 'generic' });
    const { t: tE } = useTranslation('', { keyPrefix: 'events' });
    const { setPageName } = useSettings();
    const { event, api, setEvent } = useAuth();
    const navigate = useNavigate();
    const [notif, ctxHolder] = notification.useNotification();

    const [ctx, setCtx] = useState<Status>({
        loading: true,
        page: 1,
        resp: null,
    });

    useEffect(() => setPageName('events', ['/ignore']), []);

    const fetchPage = async () => {
        setCtx((oldCtx) => ({ ...oldCtx, loading: true }));

        const resp = await api.events.getCollection(ctx.page);

        setCtx((oldCtx) => ({
            ...oldCtx,
            loading: false,
            resp,
        }));
    };

    const deleteEvent = async (id: number, name: string) => {
        await api.events.delete(id);
        await fetchPage();

        notif.success({
            message: 'Event removed',
            description: `The event ${name} was removed!`,
        });
    };

    const useEvent = async (id: number) => {
        try {
            const event = await api.settings.setEvent(id);
            if (event) {
                setEvent(event);
            }

            notif.success({
                message: 'Settings updated',
                description: 'The new event was set',
            });
        } catch (e: any) {
            if (!(e instanceof SdkError)) {
                return;
            }

            notif.error({
                message: `Failed to set event`,
                description: e.message,
            });
        }
    };

    const columns = [
        { title: t('id'), dataIndex: 'id', key: 'id', width: 64 },
        { title: t('name'), dataIndex: 'name', key: 'name', width: 256 },
        { title: t('author'), dataIndex: 'author', key: 'author', width: 128 },
        {
            title: t('actions.noun'),
            render: (_: any, x: PhEvent) => (
                <Flex gap={4}>
                    {event && event.id !== x.id && (
                        <Popconfirm
                            title={tE('select_event.title')}
                            description={tE('select_event.description')}
                            okText={t('actions.ok')}
                            cancelText={t('actions.cancel')}
                            onConfirm={() => useEvent(x.id as number)}
                        >
                            <Button>{t('actions.use')}</Button>
                        </Popconfirm>
                    )}
                    <Link to={`/events/${x.id}`}>
                        <Button>{t('actions.edit')}</Button>
                    </Link>
                    {event && event.id !== x.id && (
                        <Popconfirm
                            title={tE('delete_event.title')}
                            description={tE('delete_event.description', {
                                name: event.name,
                            })}
                            okText={t('actions.ok')}
                            cancelText={t('actions.cancel')}
                            onConfirm={() => deleteEvent(x.id as number, x.name)}
                        >
                            <Button>{t('actions.delete')}</Button>
                        </Popconfirm>
                    )}
                </Flex>
            ),
        },
    ];

    useEffect(() => {
        if (!api.tokenUser?.hasRole('ADMIN')) {
            // @TODO: Make a 403 error screen with react router
            navigate('/');
            return;
        }

        setPageName('events');
    }, []);
    useAsyncEffect(async () => await fetchPage(), [ctx.page]);

    return (
        <Flex vertical gap={2}>
            <Loader loading={ctx.loading}>
                <Typography>
                    <Title>
                        <span className="green">{tE('title')}</span>
                    </Title>
                </Typography>
                {ctx.resp?.totalCount === 0 && (
                    <Flex vertical align="center" justify="center" style={{ marginBottom: '1em' }}>
                        <Link to="/events/new">
                            <Button>{t('actions.new')}</Button>
                        </Link>
                    </Flex>
                )}
                <Table
                    pagination={{
                        position: ['bottomLeft'],
                        current: ctx.page,
                        total: ctx.resp?.totalCount,
                        pageSize: ctx.resp?.perPageCount,
                        showSizeChanger: false,
                        showTotal: (total) => {
                            return (
                                <Flex gap={16} align="center" style={{ flex: '1' }}>
                                    <Link to="/events/new">
                                        <Button>{t('actions.new')}</Button>
                                    </Link>
                                    <Typography>{tE('events_count', { amt: total })}</Typography>
                                </Flex>
                            );
                        },
                    }}
                    dataSource={ctx.resp?.results}
                    columns={columns}
                    onChange={(x) =>
                        setCtx((oldCtx) => ({
                            ...oldCtx,
                            loading: true,
                            page: x.current || 1,
                        }))
                    }
                />
                {ctxHolder}
            </Loader>
        </Flex>
    );
}
