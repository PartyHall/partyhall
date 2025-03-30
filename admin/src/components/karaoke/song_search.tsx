import { Button, Checkbox, Flex, Input, Pagination, Popover, Segmented, Typography } from 'antd';
import { Collection, PhSong } from '@partyhall/sdk';
import { CheckboxChangeEvent } from 'antd/es/checkbox';
import { IconFilter } from '@tabler/icons-react';
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
    const [formats, setFormats] = useState<string[]>([]);
    const [hasVocals, setHasVocals] = useState<boolean | null>(null);

    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [songs, setSongs] = useState<Collection<PhSong> | null>(null);

    useAsyncEffect(async () => {
        setLoading(true);
        setSongs(await api.songs.getCollection(page, search, formats, hasVocals));
        setLoading(false);
    }, [page, search, formats, hasVocals]);

    const onCheckboxChange = (e: CheckboxChangeEvent) => {
        if (!e.target.checked) {
            setFormats(formats.filter((y) => y !== e.target.value));
            setPage(1);

            return;
        }

        if (formats.includes(e.target.value)) {
            return;
        }

        setFormats([...formats, e.target.value]);
        setPage(1);
    };

    return (
        <>
            <Flex align="center" justify="center" gap={8} style={{ width: 'min(100%, 500px)', margin: 'auto' }}>
                <Input
                    placeholder={tG('actions.search') + '...'}
                    value={search}
                    onChange={(x) => {
                        setSearch(x.target.value);
                        setPage(1);
                    }}
                />
                <Popover
                    title={t('filter.title')}
                    trigger="click"
                    content={
                        <Flex vertical gap={8}>
                            <Typography.Text>{t('filter.has_vocals')}:</Typography.Text>
                            <Segmented
                                block
                                options={[
                                    { value: false, label: t('filter.no') },
                                    { value: null, label: '-' },
                                    { value: true, label: t('filter.yes') },
                                ]}
                                value={hasVocals}
                                onChange={(x) => setHasVocals(x)}
                            />
                            <Typography.Text>{t('filter.format')}:</Typography.Text>
                            <Checkbox value="video" checked={formats.includes('video')} onChange={onCheckboxChange}>
                                {t('filter.video')}
                            </Checkbox>
                            <Checkbox value="cdg" checked={formats.includes('cdg')} onChange={onCheckboxChange}>
                                {t('filter.cdg')}
                            </Checkbox>
                            <Checkbox
                                value="transparent_video"
                                checked={formats.includes('transparent_video')}
                                onChange={onCheckboxChange}
                            >
                                {t('filter.transparent_video')}
                            </Checkbox>
                        </Flex>
                    }
                >
                    <Button icon={<IconFilter size={20} />} />
                </Popover>
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
