import { useBoothSocket } from "../../hooks/boothSocket";
import Photobooth from "./photobooth";
import Quiz from "./quiz";

export default function PartyHallUI() {
    const { appState } = useBoothSocket();

    return <>
        { ['PHOTOBOOTH', 'DISABLED'].includes(appState.current_mode) && <Photobooth /> }
        { ['QUIZ'].includes(appState.current_mode) && <Quiz /> }
    </>
}