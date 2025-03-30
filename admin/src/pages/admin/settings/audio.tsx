import { Card, Flex } from "antd";
import { useEffect, useState } from "react";
import { AudioDevices } from "@partyhall/sdk/dist/models/audio";
import SoundCard from "../../../components/settings/sound_card";
import { useMercureTopic } from "../../../hooks/mercure";
import { useSettings } from "../../../hooks/settings";
import { useTranslation } from "react-i18next";

export default function SettingsAudioPage() {
    const [audioDevices, setAudioDevices] = useState<AudioDevices | null>(null);
    const {t} = useTranslation();
    const {setPageName} = useSettings();

    useMercureTopic('/audio-devices', (x: any) => setAudioDevices(AudioDevices.fromJson(x)));

    useEffect(() => {
        setPageName('settings', ['/audio-devices']);
    }, []);

    return <Flex vertical gap={16} style={{ maxWidth: '60ch' }}>
        <Card
            title={t('settings.audio.title')}
        >
            <SoundCard mercureDevices={audioDevices} />
        </Card>
    </Flex>
}