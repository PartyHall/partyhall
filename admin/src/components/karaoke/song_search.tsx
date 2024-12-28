import { Collection, PhSong } from '@partyhall/sdk';
import { Flex, Input, Pagination, Typography } from 'antd';

import Loader from '../loader';
import SongCard from './song_card';

import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export default function SongSearch() {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });
    const { api } = useAuth();

    const [search, setSearch] = useState<string>('');

    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [songs, setSongs] = useState<Collection<PhSong> | null>(null);

    useAsyncEffect(async () => {
        setLoading(true);
        setSongs(await api.karaoke.getCollection(page, search));
        setLoading(false);
    }, [page, search]);

    return (
        <>
            <Flex align="center" justify="center">
                <Input
                    style={{ width: 'min(100%, 500px)' }}
                    placeholder={tG('actions.search') + '...'}
                    value={search}
                    onChange={(x) => {
                        setSearch(x.target.value);
                        setPage(1);
                    }}
                />
            </Flex>
            <Loader loading={loading}>
                {(!songs || songs.totalCount === 0) && (
                    <Flex align="center" justify="center">
                        <Typography.Title level={2}>{t('no_songs_found')}</Typography.Title>
                    </Flex>
                )}
                {songs && songs.results.length > 0 && (
                    <>
                        <Flex vertical style={{ overflowY: 'scroll', flex: 1 }} align="center">
                            <Flex vertical gap={8} align="stretch" style={{ width: 'min(100%, 500px)' }}>
                                {songs.results.map((x) => (
                                    <SongCard key={x.nexus_id} song={x} type="SEARCH" />
                                ))}
                            </Flex>
                        </Flex>
                    </>
                )}
                <Flex align="center" justify="space-between">
                    {songs && songs.results.length !== 0 && (
                        <Pagination
                            style={{ margin: 'auto' }}
                            total={songs.totalCount ?? 10}
                            pageSize={30} // @TODO: Default API platform one but we should add it to the hydra thing so that the front knows it
                            showTotal={(total) => t('amt_songs', { amt: total })}
                            showSizeChanger={false}
                            current={page}
                            onChange={(x) => setPage(x)}
                        />
                    )}
                </Flex>
            </Loader>
        </>
    );
}
