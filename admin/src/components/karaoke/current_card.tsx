import '../../assets/css/song_card.scss';
import { Button, Card, Flex, Popconfirm, Tooltip, Typography } from 'antd';
import { IconPlayerPause, IconPlayerPlay, IconRecordMail, IconVolume, IconX } from '@tabler/icons-react';
import Image from '../image';
import TextSlider from '../text_slider';
import VolumeSlider from './volume_slider';
import { useAuth } from '../../hooks/auth';
import { useTranslation } from 'react-i18next';

function formatSeconds(seconds: number): string {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;

    const formattedMinutes = minutes.toString().padStart(2, '0');
    const formattedSeconds = remainingSeconds.toString().padStart(2, '0');

    return `${formattedMinutes}:${formattedSeconds}`;
}

type Props = {
    className?: string;
};

export default function CurrentCard({ className }: Props) {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });
    const { karaoke, api, kioskMode } = useAuth();

    if (!karaoke || !karaoke.current) {
        return;
    }

    const setStatus = async (status: boolean) => await api.karaoke.togglePlay(status);

    const cancelCurrentSong = async () => {
        if (!karaoke || !karaoke.current) {
            return;
        }

        await api.songSessions.cancelSession(karaoke.current.id);
    };

    return (
        <Card className={`SongCard ${className || ''}`}>
            <Flex gap={10} vertical align="stretch">
                <Flex gap={8} align="center">
                    <Image
                        hasImage={karaoke.current.song.has_cover}
                        alt={t('cover_alt', {
                            title: karaoke.current.song.title,
                        })}
                        src={api.songs.getCoverUrl(karaoke.current.song.nexus_id)}
                        className="SongCard__Cover"
                    />

                    <Flex vertical flex={1}>
                        <Typography.Text className="SongCard__Title">{karaoke.current.title}</Typography.Text>
                        <Typography.Text>{karaoke.current.artist}</Typography.Text>
                        {karaoke.current.sung_by.toLowerCase() !== 'kiosk' && (
                            <Typography.Text className="SongCard__Singer">
                                {t('singer')}: {karaoke.current.sung_by}
                            </Typography.Text>
                        )}
                    </Flex>
                    <Flex vertical gap={8}>
                        {!karaoke.isPlaying && (
                            <Tooltip title={tG('actions.resume')}>
                                <Button
                                    style={{ padding: kioskMode ? '1.5em' : 0 }}
                                    icon={<IconPlayerPlay size={kioskMode ? 30 : 20} />}
                                    shape="circle"
                                    onClick={() => setStatus(true)}
                                />
                            </Tooltip>
                        )}
                        {karaoke.isPlaying && (
                            <Tooltip title={tG('actions.pause')}>
                                <Button
                                    style={{ padding: kioskMode ? '1.5em' : 0 }}
                                    icon={<IconPlayerPause size={kioskMode ? 30 : 20} />}
                                    shape="circle"
                                    onClick={() => setStatus(false)}
                                />
                            </Tooltip>
                        )}
                        <Tooltip title={tG('actions.cancel')}>
                            <Popconfirm
                                title={<span className={kioskMode ? 'tooltipKiosk' : ''}>{t('confirm_cancel')}</span>}
                                okButtonProps={{ style: { padding: kioskMode ? '1.5em' : 0 } }}
                                cancelButtonProps={{ style: { padding: kioskMode ? '1.5em' : 0 } }}
                                onConfirm={cancelCurrentSong}
                                okText={tG('actions.ok')}
                                cancelText={tG('actions.cancel')}
                            >
                                <Button
                                    style={{ padding: kioskMode ? '1.5em' : 0 }}
                                    icon={<IconX size={kioskMode ? 30 : 20} />}
                                    shape="circle"
                                />
                            </Popconfirm>
                        </Tooltip>
                    </Flex>
                </Flex>

                <TextSlider
                    leftText={
                        <Typography.Text className="SongCard__Timecode">
                            {formatSeconds(karaoke.timecode)}
                        </Typography.Text>
                    }
                    rightText={
                        <Typography.Text className="SongCard__Timecode">
                            {formatSeconds(karaoke.current.song.duration)}
                        </Typography.Text>
                    }
                    value={karaoke.timecode}
                    max={karaoke.current.song.duration}
                    className="SongCard__Slider"
                    disabled
                />

                <VolumeSlider type="instrumental" icon={<IconVolume size={20} />} tooltip={t('volume')} />
                {karaoke.current.song.has_vocals && (
                    <VolumeSlider type="vocals" icon={<IconRecordMail size={20} />} tooltip={t('voices')} />
                )}
            </Flex>
        </Card>
    );
}
