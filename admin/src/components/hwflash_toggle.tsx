import { Button, Spin, Tooltip } from "antd";
import { IconBulb, IconBulbOff } from "@tabler/icons-react";
import useAsyncEffect from "use-async-effect";
import { useAuth } from "../hooks/auth";
import { useMercureTopic } from "../hooks/mercure";
import { useState } from "react";
import { useTranslation } from "react-i18next";

type MercureFlashState = {
    powered: boolean;
    brightness: number;
}

type Props = {
    setPoweredState?: (powered: boolean) => void;
};

export function HardwareFlashToggle({setPoweredState}: Props) {
    const { t } = useTranslation();
    const { api } = useAuth();
    const [loading, setLoading] = useState<boolean>(true);
    const [powered, setPowered] = useState<boolean>(false);

    const toggle = async () => await api.state.setFlash(!powered);

    useMercureTopic('/flash', (data: MercureFlashState) => {
        if (!data) {
            return;
        }

        setPowered(data.powered);
        if (setPoweredState) {
            setPoweredState(data.powered);
        }
    });

    useAsyncEffect(async () => {
        setLoading(true);
        setPowered(await api.state.getFlash());
        setLoading(false);
    }, []);

    return <Tooltip title={t('generic.hw_flash.turn_' + (powered ? 'off' : 'on'))}>
        <Button
            disabled={loading}
            icon={loading ? <Spin spinning /> : powered ? <IconBulb /> : <IconBulbOff />}
            onClick={toggle}
        />
    </Tooltip>;
}