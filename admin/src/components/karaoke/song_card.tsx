import '../../assets/css/song_card.scss';

import { Button, Card, Flex, Popconfirm, Tooltip, Typography } from 'antd';
import { IconCirclePlus, IconMicrophone, IconPiano, IconPlayerPlay, IconVinyl } from '@tabler/icons-react';
import Image from '../image';
import { PhSong } from '@partyhall/sdk';
import { useAuth } from '../../hooks/auth';
import useNotification from 'antd/es/notification/useNotification';
import { useTranslation } from 'react-i18next';

type CardType = 'SEARCH' | 'IN_QUEUE';

type Props = {
    song: PhSong;
    type: CardType;
};

export default function SongCard({ song, type }: Props) {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });
    const { displayName, api, karaoke } = useAuth();
    const [notif, notifCtx] = useNotification();

    const addToQueue = async (directPlay: boolean) => {
        const name = displayName ?? api.tokenUser?.name;

        if (!name) {
            return;
        }

        try {
            await api.songSessions.addToQueue(song.nexus_id, name, directPlay);
            notif.success({
                message: t('notification_added_queue.title'),
                description: t('notification_added_queue.description', {
                    title: song.title,
                }),
            });
        } catch (e) {
            notif.error({
                message: t('notification_queue_fail.title'),
                description: t('notification_queue_fail.description'),
            });

            console.error(e);
        }
    };

    return (
        <Card className="SongCard">
            <Flex gap={8} align="center">
                <Image
                    hasImage={song.has_cover}
                    alt={t('cover_alt', { title: song.title })}
                    src={api.songs.getCoverUrl(song.nexus_id)}
                    className="SongCard__Cover"
                />

                <Flex vertical flex={1}>
                    <Typography.Text className="SongCard__Title">{song.title}</Typography.Text>
                    <Typography.Text>{song.artist}</Typography.Text>
                    <Typography.Text className="SongCard__Format">{song.format}</Typography.Text>

                    <Flex gap={8} className="SongCard__Tracks">
                        <Tooltip title={t('instrumental')}>
                            <IconPiano size={20} color="#fafa" />
                        </Tooltip>
                        <Tooltip title={t('vocals')}>
                            <IconMicrophone size={20} color={song.has_vocals ? '#fafa' : '#777'} />
                        </Tooltip>
                        <Tooltip title={t('full_song')}>
                            <IconVinyl size={20} color={song.has_combined ? '#fafa' : '#777'} />
                        </Tooltip>
                    </Flex>
                </Flex>
                {type === 'SEARCH' && (
                    <Flex vertical align="center" justify="center" gap={8}>
                        {karaoke?.current && (
                            <Tooltip title={t('tooltip_play_directly')}>
                                <Popconfirm
                                    title={t('confirm_play_directly')}
                                    okText={tG('actions.ok')}
                                    cancelText={tG('actions.cancel')}
                                    onConfirm={() => addToQueue(true)}
                                >
                                    <Button icon={<IconPlayerPlay size={20} />} shape="circle" />
                                </Popconfirm>
                            </Tooltip>
                        )}
                        <Tooltip title={t('tooltip_add_to_queue')}>
                            <Button
                                icon={<IconCirclePlus size={20} />}
                                shape="circle"
                                onClick={() => addToQueue(false)}
                            />
                        </Tooltip>
                    </Flex>
                )}
            </Flex>

            {notifCtx}
        </Card>
    );
}
