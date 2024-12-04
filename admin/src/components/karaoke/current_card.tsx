import '../../assets/css/song_card.scss';
import { Button, Card, Flex, Slider, Tooltip, Typography } from 'antd';
import {
    IconPlayerPause,
    IconPlayerPlay,
    IconRecordMail,
    IconVolume,
} from '@tabler/icons-react';
import Image from '../image';
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

export default function CurrentCard() {
    const { t } = useTranslation('', { keyPrefix: 'karaoke' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });
    const { karaoke, api } = useAuth();

    if (!karaoke || !karaoke.current) {
        return;
    }

    const setStatus = async (status: boolean) =>
        await api.karaoke.togglePlay(status);

    return (
        <Card className="SongCard">
            <Flex gap={10} vertical align="stretch">
                <Flex gap={8} align="center">
                    <Image
                        hasImage={karaoke.current.song.has_cover}
                        alt={t('cover_alt', {
                            title: karaoke.current.song.title,
                        })}
                        src={karaoke.current.song.getCoverUrl()}
                        className="SongCard__Cover"
                    />

                    <Flex vertical flex={1}>
                        <Typography.Text className="SongCard__Title">
                            {karaoke.current.title}
                        </Typography.Text>
                        <Typography.Text>
                            {karaoke.current.artist}
                        </Typography.Text>
                        <Typography.Text className="SongCard__Singer">
                            {t('singer')}: {karaoke.current.sung_by}
                        </Typography.Text>
                    </Flex>
                    <Flex vertical gap={4}>
                        {!karaoke.isPlaying && (
                            <Tooltip title={tG('actions.resume')}>
                                <Button
                                    icon={<IconPlayerPlay size={20} />}
                                    shape="circle"
                                    onClick={() => setStatus(true)}
                                />
                            </Tooltip>
                        )}
                        {karaoke.isPlaying && (
                            <Tooltip title={tG('actions.pause')}>
                                <Button
                                    icon={<IconPlayerPause size={20} />}
                                    shape="circle"
                                    onClick={() => setStatus(false)}
                                />
                            </Tooltip>
                        )}
                    </Flex>
                </Flex>

                <Flex align="center">
                    <Typography.Text className="SongCard__Timecode">
                        {formatSeconds(karaoke.timecode)}
                    </Typography.Text>
                    <Slider
                        value={karaoke.timecode}
                        max={karaoke.current.song.duration}
                        className="SongCard__Slider"
                        disabled
                    />
                    <Typography.Text className="SongCard__Timecode">
                        {formatSeconds(karaoke.current.song.duration)}
                    </Typography.Text>
                </Flex>

                <VolumeSlider
                    type="instrumental"
                    icon={<IconVolume />}
                    tooltip={t('volume')}
                />
                {karaoke.current.song.has_vocals && (
                    <VolumeSlider
                        type="vocals"
                        icon={<IconRecordMail />}
                        tooltip={t('voices')}
                    />
                )}
            </Flex>
        </Card>
    );
}
