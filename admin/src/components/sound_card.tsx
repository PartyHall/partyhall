import { AudioDevice, AudioDevices } from '@partyhall/sdk/dist/models/audio';
import { Card, Flex, Select, Typography } from 'antd';
import { useEffect, useState } from 'react';
import Loader from './loader';
import TextSlider from './text_slider';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../hooks/auth';
import { useTranslation } from 'react-i18next';

function AdminSoundCard({ devices }: { devices: AudioDevices }) {
    const { api } = useAuth();
    const { t } = useTranslation();

    const [sinkId, setSinkId] = useState<number>(devices.defaultSink?.id || 0);
    const [sourceId, setSourceId] = useState<number>(devices.defaultSource?.id || 0);

    const updateSink = async (sinkId: number) => {
        await api.settings.setAudioDevices(
            sinkId,
            sourceId,
        );

        setSinkId(sinkId);
    };

    const updateSource = async (sourceId: number) => {
        await api.settings.setAudioDevices(
            sinkId,
            sourceId,
        );

        setSourceId(sinkId);
    };

    return (
        <>
            <Typography.Title level={4} className="no-margin red">
                {t('home.sound_settings.output_device')}
            </Typography.Title>
            <Select
                options={[
                    ...devices.sinks.map((x) => ({
                        value: x.id,
                        label: x.description,
                    })),
                ]}
                value={sinkId}
                onChange={x => updateSink(x)}
            />

            <Typography.Title level={4} className="no-margin red">
                {t('home.sound_settings.input_device')}
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
                onChange={x => updateSource(x)}
            />
        </>
    );
}

export default function SoundCard({
    mercureDevices,
}: {
    mercureDevices: AudioDevices | null;
}) {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [loadingDevices, setLoadingDevices] = useState<boolean>(true);
    const [devices, setDevices] = useState<AudioDevices | null>(null);

    useAsyncEffect(async () => {
        setLoadingDevices(true);
        setDevices(await api.settings.getAudioDevices());
        setLoadingDevices(false);
    }, []);

    useEffect(() => {
        if (!mercureDevices) {
            return;
        }

        setDevices(mercureDevices);
    }, [mercureDevices]);

    const setDeviceVolume = async (
        device: AudioDevice | null,
        volume: number
    ) => {
        if (!device) {
            return;
        }

        const data = await api.settings.setAudioDeviceVolume(device, volume);
        setDevices(data);
    };

    return (
        <Card title={t('home.sound_settings.title')}>
            <Loader loading={loadingDevices}>
                {devices && (
                    <Flex vertical gap={8}>
                        {api.tokenUser?.roles.includes('ADMIN') && (
                            <AdminSoundCard devices={devices} />
                        )}

                        <Typography.Title level={4} className="no-margin red">
                            {t('home.sound_settings.volume_output')}
                        </Typography.Title>
                        <TextSlider
                            value={
                                devices.defaultSink
                                    ? devices.defaultSink.volume * 100
                                    : 0
                            }
                            rightText={
                                <Typography.Text className="SongCard__Timecode">
                                    {devices.defaultSink
                                        ? Math.round(
                                              devices.defaultSink.volume * 100
                                          )
                                        : 0}
                                    %
                                </Typography.Text>
                            }
                            onChange={(x) =>
                                setDeviceVolume(devices.defaultSink, x)
                            }
                        />

                        <Typography.Title level={4} className="no-margin red">
                            {t('home.sound_settings.volume_input')}
                        </Typography.Title>
                        <TextSlider
                            value={devices.karaokeSink.volume * 100}
                            rightText={
                                <Typography.Text className="SongCard__Timecode">
                                    {Math.round(
                                        devices.karaokeSink.volume * 100
                                    )}
                                    %
                                </Typography.Text>
                            }
                            onChange={(x) =>
                                setDeviceVolume(devices.karaokeSink, x)
                            }
                        />
                    </Flex>
                )}
                {!devices && <span>An error occured</span>}
            </Loader>
        </Card>
    );
}
