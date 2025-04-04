import { Button, Card, Flex } from 'antd';
import SettingsButtonMappings, { SettingsButtonMappingsValues } from '../../../components/settings/button_mappings';
import { useEffect, useState } from 'react';
import { IconDeviceFloppy } from '@tabler/icons-react';
import { useAuth } from '../../../hooks/auth';
import { useSettings } from '../../../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function SettingsButtonMappingsPage() {
    const { t } = useTranslation();
    const { setPageName } = useSettings();
    const { api } = useAuth();

    const [buttonMappings, setButtonMappings] = useState<SettingsButtonMappingsValues|null>(null);

    useEffect(() => setPageName('settings'), []);

    const save = async() => {
        if (!buttonMappings) {
            return;
        }

        await api.settings.setButtonMappings(buttonMappings);
        setButtonMappings(null);
    };

    return (
        <Flex vertical gap={16} style={{ maxWidth: '60ch' }}>
            <Card
                title={t('settings.btn_mappings.title')}
                extra={<Button type="primary" icon={<IconDeviceFloppy size={18} />} disabled={!buttonMappings} onClick={save}/>
                }
            >
                <SettingsButtonMappings
                    onSettingsChanged={(x) => setButtonMappings(x)}
                />
            </Card>
        </Flex>
    );
}
