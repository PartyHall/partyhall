import { AudioDevice, AudioDevices } from '@partyhall/sdk/dist/models/audio';
import { Flex, Select, Typography } from 'antd';
import { useEffect, useState } from 'react';
import Loader from '../loader';
import TextSlider from '../text_slider';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useTranslation } from 'react-i18next';

export default function SoundCard({ mercureDevices }: { mercureDevices: AudioDevices | null }) {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [loadingDevices, setLoadingDevices] = useState<boolean>(true);
    const [devices, setDevices] = useState<AudioDevices | null>(null);

    const [sinkId, setSinkId] = useState<number>(devices?.defaultSink?.id || 0);
    const [sourceId, setSourceId] = useState<number>(devices?.defaultSource?.id || 0);

    const updateSink = async (sinkId: number) => {
        await api.settings.setAudioDevices(sinkId, sourceId);

        setSinkId(sinkId);
    };

    const updateSource = async (sourceId: number) => {
        await api.settings.setAudioDevices(sinkId, sourceId);

        setSourceId(sinkId);
    };

    useAsyncEffect(async () => {
        setLoadingDevices(true);

        const data = await api.settings.getAudioDevices();

        setDevices(data);
        setSinkId(data?.defaultSink?.id || 0);
        setSourceId(data?.defaultSource?.id || 0);

        setLoadingDevices(false);
    }, []);

    useEffect(() => {
        if (!mercureDevices) {
            return;
        }

        setDevices(mercureDevices);
        setSinkId(mercureDevices.defaultSink?.id || 0);
        setSourceId(mercureDevices.defaultSource?.id || 0);
    }, [mercureDevices]);

    const setDeviceVolume = async (device: AudioDevice | null, volume: number) => {
        if (!device) {
            return;
        }

        const data = await api.settings.setAudioDeviceVolume(device, volume);
        if (data) {
            setDevices(data);
            setSinkId(data.defaultSink?.id || 0);
            setSourceId(data.defaultSource?.id || 0);
        }
    };

    return (
        <Loader loading={loadingDevices}>
            {devices && (
                <Flex vertical gap={8}>
                    <Typography.Title level={4} className="no-margin red">
                        {t('settings.audio.output_device')}
                    </Typography.Title>
                    <Select
                        options={[
                            ...devices.sinks.map((x) => ({
                                value: x.id,
                                label: x.description,
                            })),
                        ]}
                        value={sinkId}
                        onChange={(x) => updateSink(x)}
                    />

                    <Typography.Title level={4} className="no-margin red">
                        {t('settings.audio.input_device')}
                    </Typography.Title>
                    {/* For some reason, the device id is displayed when changed but works ok on refresh */}
                    <Select
                        options={[
                            ...devices.sources.map((x) => ({
                                value: x.id,
                                label: x.description,
                            })),
                        ]}
                        value={sourceId}
                        onChange={(x) => updateSource(x)}
                    />

                    <Typography.Title level={4} className="no-margin red">
                        {t('settings.audio.volume_output')}
                    </Typography.Title>
                    <TextSlider
                        value={devices.defaultSink ? devices.defaultSink.volume * 100 : 0}
                        rightText={
                            <Typography.Text className="SongCard__Timecode">
                                {devices.defaultSink ? Math.round(devices.defaultSink.volume * 100) : 0}%
                            </Typography.Text>
                        }
                        onChange={(x) => setDeviceVolume(devices.defaultSink, x)}
                    />

                    <Typography.Title level={4} className="no-margin red">
                        {t('settings.audio.volume_input')}
                    </Typography.Title>
                    <TextSlider
                        value={devices.karaokeSink.volume * 100}
                        rightText={
                            <Typography.Text className="SongCard__Timecode">
                                {Math.round(devices.karaokeSink.volume * 100)}%
                            </Typography.Text>
                        }
                        onChange={(x) => setDeviceVolume(devices.karaokeSink, x)}
                    />
                </Flex>
            )}
            {!devices && <span>An error occured</span>}
        </Loader>
    );
}
