import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Card, Grid, List, ListItem, Stack, Typography } from "@mui/material";

import {faSpotify as SpotifyIcon} from '@fortawesome/free-brands-svg-icons';
import { useState } from "react";

type Props = {
    artist: string;
    title: string;

    onChange?: (url: string) => void;
};

type ApiSong = {
    artist: string;
    song: string;
    cover: string;
}

export default function SearchSpotify({ artist, title, onChange }: Props) {
    const [loading, setLoading] = useState<boolean>(false);
    const [results, setResults] = useState<ApiSong[]|null>(null);

    const [selectedSong, setSelectedSong] = useState<ApiSong|null>(null);

    const searchSpotify = async () => {
        setLoading(true);

        let resp = await fetch(`/api/modules/karaoke/search/spotify?q=${encodeURI(artist + ' ' + title)}`);
        resp = await resp.json();
        //@ts-ignore
        setResults(resp as ApiSong[]);

        setLoading(false);
    };

    return <Stack alignItems="stretch" justifyContent="center" gap={2}>
        <Button
            variant="contained"
            startIcon={<FontAwesomeIcon icon={SpotifyIcon} />}
            style={{ backgroundColor: '#1ED760' }}
            onClick={searchSpotify}
            disabled={loading || artist.length == 0 || title.length == 0}
        >
            Search
        </Button>

        {
            results && results?.length > 0 && <Card elevation={2}>
                <Stack 
                    direction="column"
                    height="350px"
                    overflow="scroll"
                    alignContent="stretch"
                    p={1}
                >
                    {
                        results && results.map(x => <Stack
                                key={x.artist+x.song+x.cover}
                                style={selectedSong === x ? {backgroundColor: 'red'} : {}}
                                direction="row"
                                gap={2}
                                p={1}
                                onClick={() => {
                                    if (selectedSong?.cover === x.cover) {
                                        setSelectedSong(null);
                                    } else {
                                        setSelectedSong(x);
                                        if (onChange) {
                                            onChange(x.cover);
                                        }
                                    }
                                }}
                            >
                                <img
                                    src={x.cover}
                                    alt={`Album cover of ${x.song}`}
                                    width="150px"
                                    style={{
                                        display: 'block',
                                        objectFit: 'contain',
                                    }}
                                />
                                <Stack direction="column">
                                    <Typography variant="body1" textOverflow="clip">{x.song}</Typography>
                                    <Typography variant="body1">{x.artist}</Typography>
                                </Stack>
                            </Stack>
                        )
                    }
                </Stack>
            </Card>
        }
    </Stack>
}