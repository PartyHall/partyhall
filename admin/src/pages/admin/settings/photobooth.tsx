import { Button, Card, Flex } from "antd";
import SettingsHardwareFlash, { SettingsHardwareFlashValues } from "../../../components/settings/hwflash";
import SettingsUnattended, { SettingsUnattendedValues } from "../../../components/settings/unattended";
import SettingsWebcam, { SettingsWebcamValues } from "../../../components/settings/webcam";
import { IconDeviceFloppy } from "@tabler/icons-react";
import useAsyncEffect from "use-async-effect";
import { useAuth } from "../../../hooks/auth";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export default function SettingsPhotoboothPage() {
    const {t} = useTranslation();
    const {api} = useAuth();

    const [webcamSettings, setWebcamSettings] = useState<SettingsWebcamValues|null>(null);
    const [unattendedSettings, setUnattendedSettings] = useState<SettingsUnattendedValues|null>(null);
    const [hardwareFlashSettings, setHardwareFlashSettings] = useState<SettingsHardwareFlashValues|null>(null);

    const onWebcamSave = async () => {
        if (!webcamSettings) {
            return;
        }

        await api.settings.setWebcam(webcamSettings.width, webcamSettings.height);
        setWebcamSettings(null);
    }

    const onUnattendedSave = async () => {
        if (!unattendedSettings) {
            return;
        }

        await api.settings.setUnattended(unattendedSettings.enabled, unattendedSettings.interval);
        setUnattendedSettings(null);
    }

    useAsyncEffect(async () => {
        if (!hardwareFlashSettings) {
            return;
        }

        await api.settings.setFlash(hardwareFlashSettings.powered, hardwareFlashSettings.brightness);
        setHardwareFlashSettings(null);
    }, [hardwareFlashSettings]);

    return <Flex vertical gap={16} style={{maxWidth: '60ch'}}>
        <Card
            title={t('settings.webcam.title')}
            extra={<Button type="primary" icon={<IconDeviceFloppy size={18} />} disabled={!webcamSettings} onClick={onWebcamSave}/>}
        >
            <SettingsWebcam onSettingsChanged={x => setWebcamSettings(x)}/>
        </Card>
        <Card
            title={t('settings.unattended.title')}
            extra={<Button type="primary" icon={<IconDeviceFloppy size={18} />} disabled={!unattendedSettings} onClick={onUnattendedSave} />}
        >
            <SettingsUnattended onSettingsChanged={x => setUnattendedSettings(x)}/>
        </Card>
        <Card title={t('settings.hwflash.title')}>
            <SettingsHardwareFlash onSettingsChanged={x => setHardwareFlashSettings(x)}/>
        </Card>
    </Flex>
}