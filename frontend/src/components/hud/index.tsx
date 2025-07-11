import '../../assets/hud.scss';
import DateTimeRender from './datetime';
import DisplayText from './display_text';
import EventRenderer from './event';
import { FORCE_DEBUG } from './debug';
import WifiRenderer from './wifi_renderer';
import { useAuth } from '../../hooks/auth';

export default function Hud() {
    const { debug } = useAuth();

    return (
        <div id="hud">
            <div className="hud-top">
                <EventRenderer />
                <DateTimeRender />
            </div>

            <div className="hud-bottom">
                {(debug || FORCE_DEBUG) && <WifiRenderer />}
                {!debug && !FORCE_DEBUG && <DisplayText />}
            </div>
        </div>
    );
}
