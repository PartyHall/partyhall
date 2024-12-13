import { Card, Flex, Select, Typography } from "antd";
import { AudioDevices } from "@partyhall/sdk/dist/models/audio";
import Loader from "./loader";
import TextSlider from "./text_slider";
import useAsyncEffect from "use-async-effect";
import { useAuth } from "../hooks/auth";
import { useState } from "react";
import { useTranslation } from "react-i18next";

function AdminSoundCard({ devices }: { devices: AudioDevices }) {
    const { t } = useTranslation();

    return <>
        <Typography.Title level={4} className="no-margin red">{t('home.sound_settings.output_device')}</Typography.Title>
        <Select
            options={[
                ...devices.sinks.map(x => ({
                    value: x.id,
                    label: x.description,
                }))
            ]}
            defaultValue={devices.defaultSink?.id}
        />

        <Typography.Title level={4} className="no-margin red">{t('home.sound_settings.input_device')}</Typography.Title>
        <Select
            options={[
                ...devices.sources.map(x => ({
                    value: x.id,
                    label: x.description,
                }))
            ]}
            defaultValue={devices.defaultSource?.id}
        />
    </>;
}

export default function SoundCard() {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [loadingDevices, setLoadingDevices] = useState<boolean>(true);
    const [devices, setDevices] = useState<AudioDevices | null>(null);

    useAsyncEffect(async () => {
        setLoadingDevices(true);
        setDevices(await api.settings.getAudioDevices());
        setLoadingDevices(false);
    }, []);

    return <Card title={t('home.sound_settings.title')}>
        <Loader loading={loadingDevices}>
            {
                devices && <Flex vertical gap={8}>
                    {api.tokenUser?.roles.includes('ADMIN') && <AdminSoundCard devices={devices} />}

                    <Typography.Title level={4} className="no-margin red">{t('home.sound_settings.volume_output')}</Typography.Title>
                    <TextSlider
                        value={devices.defaultSink ? devices.defaultSink.volume * 100 : 0}
                        rightText={
                            <Typography.Text className="SongCard__Timecode">
                                {devices.defaultSink ? Math.round(devices.defaultSink.volume * 100) : 0}%
                            </Typography.Text>
                        }
                    />

                    <Typography.Title level={4} className="no-margin red">{t('home.sound_settings.volume_input')}</Typography.Title>
                    <TextSlider
                        value={devices.karaokeSink.volume * 100}
                        rightText={
                            <Typography.Text className="SongCard__Timecode">
                                {Math.round(devices.karaokeSink.volume * 100)}%
                            </Typography.Text>
                        }
                    />
                </Flex>
            }
            {
                !devices && <span>An error occured</span>
            }
        </Loader>
    </Card>;
}