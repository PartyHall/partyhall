import { useAuth } from '../../hooks/auth';

export default function ButtonSetupModal() {
    const { lastBtnPressed } = useAuth();

    return (
        <div className="modal">
            <h1>Button setup</h1>
            <span>Press any button on the appliance to get its ID</span>
            {lastBtnPressed === null && <h1>No button pressed</h1>}
            {lastBtnPressed !== null && <h1>BTN_{lastBtnPressed}</h1>}
        </div>
    );
}
