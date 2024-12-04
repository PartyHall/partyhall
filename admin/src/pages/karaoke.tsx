import { Flex, Tabs } from 'antd';
import { PhKaraoke, PhSongSession } from '@partyhall/sdk';
import { useEffect, useState } from 'react';

import SongQueue from '../components/karaoke/song_queue';
import SongSearch from '../components/karaoke/song_search';

import { useAuth } from '../hooks/auth';
import { useMercureTopic } from '../hooks/mercure';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Karaoke() {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });

    const { setPageName } = useSettings();
    const { setKaraoke, setKaraokeQueue, setTimecode } = useAuth();

    useMercureTopic('/karaoke-timecode', (x: any) => setTimecode(x.timecode));

    useMercureTopic('/karaoke-queue', (x: any[]) => {
        setKaraokeQueue(
            x.map((y) => PhSongSession.fromJson(y)).filter((y) => !!y)
        );
    });

    useMercureTopic('/karaoke', (x: any) => {
        setKaraoke(PhKaraoke.fromJson(x));
    });

    const [active, setActive] = useState<string>('search');

    useEffect(() => {
        setPageName('karaoke', [
            '/karaoke-queue',
            '/karaoke-timecode',
            '/karaoke',
        ]);
    }, []);

    // Not using tabs the correct way
    // because it breaks the CSS

    return (
        <Flex vertical className="fullwidth fullheight overflowAuto" gap={8}>
            <Tabs
                defaultActiveKey="search"
                onChange={(x) => setActive(x)}
                items={[
                    { key: 'search', label: t('tabs.search') },
                    { key: 'queue', label: t('tabs.queue') },
                ]}
            />

            {active === 'search' && <SongSearch />}
            {active === 'queue' && <SongQueue />}
        </Flex>
    );
}
