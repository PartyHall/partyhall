import { Pagination, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import useAsyncEffect from "use-async-effect";
import Loader from "../../../components/loader";
import Song from "./song";
import { Meta } from "../../../types/contextualized_response";

// @TODO: debounce
export default function KaraokeSearch() {
    const [wasSearch, setWasSearch] = useState<boolean>(false);

    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [search, setSearch] = useState<string>("");
    const [meta, setMeta] = useState<Meta|null>(null);
    const [results, setResults] = useState<KaraokeSong[]>([]);

    const loadSongs = async () => {
        let currPage = page;
        if ((wasSearch && search.length == 0) || (!wasSearch && search.length > 0)) {
            currPage = 1;
            setWasSearch(!wasSearch);
            setPage(currPage);
        }

        setLoading(true);
        let resp: any = await fetch(`/api/modules/karaoke/song?page=${currPage}` + (search.length > 0 ? `&q=${encodeURI(search)}` : ''))

        resp = await resp.json();
        setResults(resp['results']);
        setMeta(resp['meta']);

        setLoading(false);
    };

    useAsyncEffect(async () => {
        await loadSongs();
    }, [search, page]);

    return <Stack direction="column" alignItems="stretch" flex={1} gap={1}>
        <TextField placeholder="Search..." value={search} onChange={x => setSearch(x.target.value)}/>
        <Loader loading={loading}>
            <Stack direction="column" gap={1} pt={1} flex="1 1 0" style={{overflowY: 'scroll'}}>
                { results.map(x => <Song key={x.id} song={x} type="SEARCH" mb={0} />) }
            </Stack>
        </Loader>
        <Stack direction="column" alignItems="center" justifyContent="center" gap={2}>
            <Typography variant="body1">Amt of songs: {meta?.total}</Typography>
            <Pagination 
                count={meta?.last_page}
                shape="rounded"
                variant="outlined"
                page={page}
                onChange={(event: React.ChangeEvent<unknown>, value: number) => setPage(value)}
            />
        </Stack>
    </Stack>
}