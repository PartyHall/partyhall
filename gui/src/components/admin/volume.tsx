import { InputLabel, MenuItem, Select, Slider, Stack, ToggleButton, Typography } from "@mui/material";
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { useEffect, useRef, useState } from "react";
import { useAdminSocket } from "../../hooks/adminSocket";
import { debounce } from "lodash";
import { useApi } from "../../hooks/useApi";
import { useTranslation } from "react-i18next";

import VolumeOffIcon from '@mui/icons-material/VolumeOff';

export default function VolumeAdmin() {
    const { t } = useTranslation();
    const { hasRole } = useApi();
    const { sendMessage, appState } = useAdminSocket();
    // No clue why I have to do this as the server sends back the info theorically?
    const [mute, setMute] = useState<boolean>(appState.pulseaudio_selected?.mute || false);
    const [volume, setVolume] = useState<number>(appState.pulseaudio_selected?.volume.value ?? 0);

    const debouncedSetVolume = useRef(
        debounce((vol: number) => sendMessage('SET_VOLUME', Math.floor((vol / 65536) * 100)), 250)
    ).current;

    const toggleMute = (x: boolean) => {
        sendMessage("SET_MUTE", x);
        setMute(x);
    }

    useEffect(() => {
        if (!appState.pulseaudio_selected?.volume?.value) {
            return;
        }

        setVolume(appState.pulseaudio_selected.volume.value);
        setMute(appState.pulseaudio_selected.mute || false);
    }, [appState]);

    useEffect(() => {
        return () => {
            debouncedSetVolume.cancel();
        }
    }, [debouncedSetVolume]);

    if (!appState.pulseaudio_devices || appState.pulseaudio_devices.length === 0) {
        return <>
            <Typography variant="h2" fontSize={18}>{t('admin_main.volume')}</Typography>
            <Typography variant="body1" color="GrayText">{t('admin_main.no_devices')}</Typography>
        </>
    }

    return <>
        <Typography variant="h2" fontSize={18}>{t('admin_main.volume')}</Typography>
        {
            hasRole('ADMIN') && appState.pulseaudio_devices && appState.pulseaudio_devices.length > 0 &&
            <>
                <InputLabel id="volume-device-label">{t('admin_main.device')}</InputLabel>
                <Select value={appState.pulseaudio_selected?.index} labelId="volume-device-label" onChange={x => sendMessage('SET_SOUND_CARD', x.target.value)}>
                    {
                        appState.pulseaudio_devices.map(x => <MenuItem key={x.index} value={x.index}>
                            {
                                x.description !== "(null)" ? `${x.description}` : `${x.name}`
                            }
                        </MenuItem>)
                    }
                </Select>
            </>
        }
        <Stack spacing={2} direction="row" sx={{ mb: 1 }} alignItems="center">
            <VolumeDownIcon />
            <Slider aria-label="Volume" value={volume} onChange={(_, x) => {
                setVolume(x as number);
                debouncedSetVolume(x as number);
            }} max={65536} />
            <VolumeUpIcon />
            <ToggleButton value={mute} selected={mute} onChange={() => toggleMute(!mute)}>
                <VolumeOffIcon />
            </ToggleButton>
        </Stack>

    </>
}