import { Flex, Typography } from 'antd';
import CurrentCard from './current_card';
import SessionCard from './session_card';

import { useAuth } from '../../hooks/auth';
import { useTranslation } from 'react-i18next';

export default function SongQueue() {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { karaoke, karaokeQueue } = useAuth();

    if (!karaoke) {
        return;
    }

    return (
        <Flex
            vertical
            gap={8}
            className="fullheight overflowAuto"
            align="center"
        >
            <Flex
                vertical
                style={{ width: 'min(500px, 100%)' }}
                className="fullheight overflowAuto"
            >
                {karaoke.current && <CurrentCard />}

                <Typography.Title>{t('in_queue')}</Typography.Title>
                <Flex vertical gap={8} style={{ overflowY: 'scroll' }}>
                    {karaokeQueue.map((x) => (
                        <SessionCard
                            key={x.id}
                            session={x}
                            isFirst={karaokeQueue[0].id === x.id}
                            isLast={
                                karaokeQueue[karaokeQueue.length - 1].id ===
                                x.id
                            }
                        />
                    ))}
                    {karaokeQueue.length === 0 && (
                        <Typography.Title level={3}>
                            {t('no_songs_in_queue')}
                        </Typography.Title>
                    )}
                </Flex>
            </Flex>
        </Flex>
    );
}
