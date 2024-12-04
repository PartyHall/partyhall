import { DebugLeft } from './debug';
import { useAuth } from '../../hooks/auth';

export default function EventRenderer() {
    const { currentEvent: current_event, debug } = useAuth();

    return (
        <div>
            <div>{current_event?.name}</div>
            {debug && <DebugLeft />}
        </div>
    );
}
