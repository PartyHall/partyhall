import { useAuth } from '../../hooks/auth';

export default function WifiRenderer() {
    const { userSettings } = useAuth();

    if (!userSettings || !userSettings.wirelessAp.enabled) {
        return <></>;
    }

    return (
        <span>
            Wifi: {userSettings?.wirelessAp.ssid} / {userSettings?.wirelessAp.password}
        </span>
    );
}
