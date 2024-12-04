import { Flex, Slider, Tooltip, Typography } from 'antd';
import { ReactNode } from 'react';
import { VolumeType } from '@partyhall/sdk';
import { useAuth } from '../../hooks/auth';

type Props = {
    type: VolumeType;
    tooltip: string;
    icon: ReactNode;
};

export default function VolumeSlider(props: Props) {
    const { karaoke, api, setKaraoke } = useAuth();

    if (!karaoke) {
        return;
    }

    const volume =
        props.type == 'instrumental' ? karaoke.volume : karaoke.volumeVocals;

    /**
     * @TODO: We should decouple the value so that latency is not an issue
     * i.e. do not set the value from the server directly and let the request go
     * whenever it wants so that even with bad connection the user can have
     * a responsive slider
     */
    const setVolume = async (vol: number) => {
        if (vol === volume) {
            return;
        }

        const data = await api.karaoke.setVolume(props.type, vol);
        setKaraoke(data);
    };

    return (
        <Flex align="center">
            <Tooltip title={props.tooltip}>{props.icon}</Tooltip>
            <Slider
                value={volume}
                onChange={(x) => setVolume(x)}
                className="SongCard__Slider"
            />
            <Typography.Text className="SongCard__Timecode">
                {volume}%
            </Typography.Text>
        </Flex>
    );
}
