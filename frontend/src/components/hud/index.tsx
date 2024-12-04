import '../../assets/hud.scss';

import DateTimeRender from './datetime';
import EventRenderer from './event';

export default function Hud() {
    return (
        <div id="hud">
            <EventRenderer />
            <DateTimeRender />
        </div>
    );
}
