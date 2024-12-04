import { DebugRight } from './debug';
import { useAuth } from '../../hooks/auth';

export default function DateTimeRender() {
    const { time, debug } = useAuth();

    return (
        <div>
            <div>
                {!time && 'Failed to reach the server'}
                {time && time.toFormat('HH:mm:ss - dd/MM/yyyy')}
            </div>
            {debug && <DebugRight />}
        </div>
    );
}
