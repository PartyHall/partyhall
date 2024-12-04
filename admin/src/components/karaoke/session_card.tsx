import '../../assets/css/song_card.scss';
import { Button, Card, Flex, Popconfirm, Tooltip, Typography } from 'antd';
import {
    IconCaretDownFilled,
    IconCaretUpFilled,
    IconPlayerPlayFilled,
    IconTrash,
} from '@tabler/icons-react';
import Image from '../image';
import { PhSongSession } from '@partyhall/sdk';
import { useAuth } from '../../hooks/auth';
import useNotification from 'antd/es/notification/useNotification';
import { useTranslation } from 'react-i18next';

type Props = {
    session: PhSongSession;
    isFirst: boolean;
    isLast: boolean;
};

export default function SessionCard({ session, isFirst, isLast }: Props) {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });
    const { setKaraokeQueue, api } = useAuth();
    const [notif, notifCtx] = useNotification();

    const directPlay = async () => {
        await api.karaoke.directPlay(session.id);
    };

    const moveSession = async (direction: 'up' | 'down') => {
        const data = await api.karaoke.moveInQueue(session.id, direction);
        setKaraokeQueue(data);
    };

    const cancelSession = async () => {
        try {
            await api.karaoke.cancelSession(session.id);
            notif.success({
                message: t('notification_remove_queue.title'),
                description: t('notification_remove_queue.description', {
                    title: session.title,
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
            <Flex gap={10} vertical align="stretch">
                <Flex gap={16} align="center">
                    <Flex vertical gap={8}>
                        <Button
                            icon={<IconCaretUpFilled size={18} />}
                            shape="circle"
                            disabled={isFirst}
                            onClick={() => moveSession('up')}
                        />
                        <Button
                            icon={<IconCaretDownFilled size={18} />}
                            shape="circle"
                            disabled={isLast}
                            onClick={() => moveSession('down')}
                        />
                    </Flex>
                    <Image
                        hasImage={session.song.has_cover}
                        alt={t('cover_alt', { title: session.title })}
                        src={session.song.getCoverUrl()}
                        className="SongCard__Cover"
                    />

                    <Flex vertical flex={1}>
                        <Typography.Text className="SongCard__Title">
                            {session.title}
                        </Typography.Text>
                        <Typography.Text>{session.artist}</Typography.Text>
                        <Typography.Text className="SongCard__Singer">
                            {t('singer')}: {session.sung_by}
                        </Typography.Text>
                    </Flex>
                    <Flex vertical gap={8}>
                        <Tooltip title={t('tooltip_play_directly')}>
                            <Popconfirm
                                title={t('confirm_play_directly')}
                                onConfirm={directPlay}
                                okText={tG('actions.ok')}
                                cancelText={tG('actions.cancel')}
                            >
                                <Button
                                    icon={<IconPlayerPlayFilled size={18} />}
                                    shape="circle"
                                />
                            </Popconfirm>
                        </Tooltip>
                        <Tooltip title={tG('actions.cancel')}>
                            <Popconfirm
                                title={t('confirm_cancel')}
                                onConfirm={cancelSession}
                                okText={tG('actions.ok')}
                                cancelText={tG('actions.cancel')}
                            >
                                <Button
                                    icon={<IconTrash size={18} />}
                                    shape="circle"
                                />
                            </Popconfirm>
                        </Tooltip>
                    </Flex>
                </Flex>
            </Flex>

            {notifCtx}
        </Card>
    );
}
