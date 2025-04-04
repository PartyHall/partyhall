import { Button, Tooltip } from 'antd';
import { IconBulb, IconBulbOff } from '@tabler/icons-react';
import { useAuth } from '../hooks/auth';
import { useTranslation } from 'react-i18next';

export function HardwareFlashToggle() {
    const { hardwareFlashPowered } = useAuth();
    const { t } = useTranslation();
    const { api } = useAuth();

    const toggle = async () => await api.state.setFlash(!hardwareFlashPowered);

    return (
        <Tooltip title={t('generic.hw_flash.turn_' + (hardwareFlashPowered ? 'off' : 'on'))}>
            <Button
                icon={hardwareFlashPowered ? <IconBulb /> : <IconBulbOff />}
                onClick={toggle}
            />
        </Tooltip>
    );
}
