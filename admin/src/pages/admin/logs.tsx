import { Flex, Table } from 'antd';
import { DateTime } from 'luxon';
import { Log } from '@partyhall/sdk/dist/models/log';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useMercureTopic } from '../../hooks/mercure';
import { useSettings } from '../../hooks/settings';
import { useState } from 'react';

function RenderText({ type, msg }: { type: string; msg: string }) {
    const style = {
        background: 'inherit',
    };

    if (type.toLowerCase().startsWith('error')) {
        style['background'] = '#ff000044';
    }

    return <span style={style}>{msg}</span>;
}

export default function Logs() {
    const { setPageName } = useSettings();
    const { api } = useAuth();
    const [logs, setLogs] = useState<Log[]>([]);

    useMercureTopic('/logs', (x: any) => {
        setLogs((oldLogs) => [
            {
                ...x,
                timestamp: DateTime.fromISO(x.timestamp).toFormat(
                    'yyyy-MM-dd HH:mm:ss'
                ),
            },
            ...oldLogs,
        ]);
    });

    useAsyncEffect(async () => {
        setPageName('logs', ['/logs']);

        const logs = await api.global.getLogs();
        if (!logs) {
            return;
        }

        setLogs(logs);
    }, []);

    const columns = [
        {
            title: 'TS',
            dataIndex: 'timestamp',
            render: (_: any, x: any) => (
                <RenderText type={x.type} msg={x.timestamp} />
            ),
        },
        {
            title: 'Type',
            dataIndex: 'type',
            render: (_: any, x: any) => (
                <RenderText type={x.type} msg={x.type} />
            ),
        },
        { title: 'Message', dataIndex: 'text', key: 'id' },
    ];

    /** @TODO: Make a real console like in reposilite */
    /** @TODO: https://blog.bowen.cool/posts/infinite-scrolling-antd-table/ */
    /** https://ahooks.js.org/hooks/use-infinite-scroll */
    /** ahooks could also permit removing useAsyncEffect as they have one */
    return (
        <Flex vertical>
            <Table
                pagination={{
                    position: ['bottomLeft'],
                    current: 1,
                    total: 100,
                    pageSize: 100,
                    showSizeChanger: false,
                }}
                dataSource={logs}
                columns={columns}
            />
        </Flex>
    );
}
