import { Select, Spin } from 'antd';
import { BackdropAlbum } from '@partyhall/sdk';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export default function BackdropSelector() {
    const { t } = useTranslation();
    const { api, backdropAlbum, selectedBackdrop } = useAuth();
    const [backdropAlbums, setBackdropAlbums] = useState<BackdropAlbum[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useAsyncEffect(async () => {
        setLoading(true);
        setBackdropAlbums((await api.backdrops.getAlbumCollection(1, '')).results);
        setLoading(false);
    }, []);

    const setSelectedBackdropAlbum = async (id: number) =>
        await api.state.setBackdrops(id !== 0 ? id : null, selectedBackdrop);

    return (
        <Spin spinning={loading}>
            <Select
                style={{ width: '100%' }}
                allowClear
                options={[
                    { value: 0, label: `-- ${t('generic.none')} --` },
                    ...backdropAlbums.map((x) => ({
                        value: x.id,
                        label: x.name,
                    })),
                ]}
                value={backdropAlbum?.id ?? 0}
                onChange={setSelectedBackdropAlbum}
            />
        </Spin>
    );
}
