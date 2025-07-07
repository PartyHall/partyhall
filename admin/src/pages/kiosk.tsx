import '../assets/css/kiosk.css';
import { PhKaraoke, PhSongSession } from '@partyhall/sdk';
import CurrentCard from '../components/karaoke/current_card';
import { Flex } from 'antd';
import PhLogo from '../assets/ph_logo_sd.webp';
import SongQueue from '../components/karaoke/song_queue';
import SongSearch from '../components/karaoke/song_search';
import { useAuth } from '../hooks/auth';
import { useEffect } from 'react';
import { useMercureTopic } from '../hooks/mercure';
import { useSettings } from '../hooks/settings';

export default function KioskPage() {
    const { setPageName } = useSettings();
    const { setKaraoke, setKaraokeQueue, setTimecode, event } = useAuth();

    useMercureTopic('/karaoke-timecode', (x: any) => setTimecode(x.timecode));
    useMercureTopic('/karaoke-queue', (x: any[]) => {
        setKaraokeQueue(x.map((y) => PhSongSession.fromJson(y)).filter((y) => !!y));
    });

    useMercureTopic('/karaoke', (x: any) => {
        setKaraoke(PhKaraoke.fromJson(x));
    });

    useEffect(() => {
        setPageName('karaoke', ['/karaoke-queue', '/karaoke-timecode', '/karaoke']);
    }, []);

    return (
        <div className="kioskMode">
            <div className="topBar">
                <Flex vertical align="center" gap={20}>
                    <img src={PhLogo} style={{ height: '4em' }} />
                    <span style={{ fontWeight: 'bold', color: 'white' }}>
                        Inscrivez-vous pour accéder aux photos de la soirée via ce QR code:
                    </span>
                </Flex>
                <img className="qr" src={`/api/events/${event?.id}/registration-qr`} />
            </div>

            <div className="leftSide">
                <h1>Rechercher une musique</h1>
                <SongSearch />
            </div>

            <div className="rightSide">
                <h1>File d&apos;attente</h1>
                <SongQueue hideCurrent />
            </div>

            <CurrentCard className="bottomSide" />
        </div>
    );
}
