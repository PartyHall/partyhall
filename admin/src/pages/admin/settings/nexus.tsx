import { Button, Card, Flex, Popconfirm } from 'antd';
import { IconDeviceFloppy, IconTrash } from '@tabler/icons-react';
import SettingsNexus, { SettingsNexusValues } from '../../../components/settings/nexus';
import { useEffect, useState } from 'react';
import NexusSettings from '@partyhall/sdk/dist/models/nexus';
import { useAuth } from '../../../hooks/auth';
import { useSettings } from '../../../hooks/settings';
import { useTranslation } from 'react-i18next';

const initialData = {
    nexusUrl: '',
    hardwareId: '',
    apiKey: '',
    bypassSsl: false,
};

export default function SettingsNexusPage() {
    const { t } = useTranslation();
    const { setPageName } = useSettings();
    const { api } = useAuth();

    const [dirty, setDirty] = useState<boolean>(false);
    const [connectionErrorMessage, setConnectionErrorMessage] = useState<string | null>(null);
    const [nexusSettings, setNexusSettings] = useState<SettingsNexusValues | null>(initialData);
    const [data, setData] = useState<NexusSettings | null>(null);

    useEffect(() => setPageName('settings'), []);

    const testAndSave = async (newData?: SettingsNexusValues) => {
        if (!nexusSettings) {
            return;
        }

        const appliedData = newData ?? nexusSettings;

        const resp = await api.settings.setNexus(
            appliedData.nexusUrl,
            appliedData.hardwareId,
            appliedData.apiKey,
            appliedData.bypassSsl
        );

        if (resp) {
            setConnectionErrorMessage(resp.errorMessage ?? null);
            if (!resp.errorMessage) {
                setNexusSettings(null);
            }
        }

        setData(resp);
    };

    return (
        <Flex vertical gap={16} style={{ maxWidth: '60ch' }}>
            <Card
                title={t('settings.nexus.title')}
                extra={
                    <Flex gap={4}>
                        <Popconfirm
                            title={t('settings.nexus.remove_settings')}
                            onConfirm={() => testAndSave(initialData)}
                            okText={t('generic.actions.ok')}
                            cancelText={t('generic.actions.cancel')}
                        >
                            <Button type="primary" icon={<IconTrash size={18} />} />
                        </Popconfirm>

                        <Button
                            type="primary"
                            icon={<IconDeviceFloppy size={18} />}
                            disabled={!dirty}
                            onClick={() => testAndSave()}
                        />
                    </Flex>
                }
            >
                <SettingsNexus
                    onSettingsChanged={(x) => {
                        setNexusSettings(x);
                        setDirty(true);
                    }}
                    errorMessage={connectionErrorMessage}
                    initialData={data}
                />
            </Card>
        </Flex>
    );
}
