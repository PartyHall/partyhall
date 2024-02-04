import { useBoothSocket } from "../../hooks/boothSocket";
import Karaoke from "./karaoke";
import Photobooth from "./photobooth";

export default function PartyHallUI() {
    const { appState } = useBoothSocket();

    return <>
        { ['PHOTOBOOTH', 'DISABLED'].includes(appState.current_mode) && <Photobooth /> }
        { ['KARAOKE'].includes(appState.current_mode) && <Karaoke /> }
    </>
}