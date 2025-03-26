import { DebugLeft, FORCE_DEBUG } from './debug';
import { useAuth } from '../../hooks/auth';

export default function EventRenderer() {
    const { currentEvent: current_event, debug } = useAuth();

    return (
        <div>
            <div>{current_event?.name}</div>
            {(debug || FORCE_DEBUG) && <DebugLeft />}
        </div>
    );
}
