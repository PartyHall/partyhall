import { Stack, TextField } from "@mui/material";
import { useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import useAsyncEffect from "use-async-effect";
import Loader from "../../../components/loader";
import Song from "./song";

// @TODO: debounce
export default function KaraokeSearch() {
    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [search, setSearch] = useState<String>("");
    const [results, setResults] = useState<KaraokeSong[]>([]);

    useAsyncEffect(async () => {
        setLoading(true);
        let resp;
        if (search.length > 0) {
            resp = await fetch(`/api/modules/karaoke/search_song?q=${search}`);
        } else {
            resp = await fetch(`/api/modules/karaoke/list_song?page=${page}`)
        }

        setResults(await resp.json());
        setLoading(false);
    }, [search]);

    return <Stack direction="column">
        <TextField placeholder="Search..." value={search} onChange={x => setSearch(x.target.value)}/>
        <Loader loading={loading}>
            <Stack direction="column" gap={1} pt={1}>
                { results.map(x => <Song key={x.id} song={x} type="SEARCH" mb={0} />) }
            </Stack>
        </Loader>
    </Stack>
}