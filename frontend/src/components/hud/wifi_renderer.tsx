import { useAuth } from '../../hooks/auth';

export default function WifiRenderer() {
    const { userSettings } = useAuth();

    if (!userSettings || !userSettings.wirelessAp.enabled) {
        return <></>;
    }

    return (
        <div>
            <div className="wifi-renderer">
                <h1>Wifi</h1>
                <span>{userSettings?.wirelessAp.ssid}</span>
                <span>{userSettings?.wirelessAp.password}</span>
                <img src="/api/state/ap-qr" alt="Wifi QR Code" />
            </div>
        </div>
    );
}
