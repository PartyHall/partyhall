import { Button, Card, CardActions, CardContent, Grid, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { useApi } from "../hooks/auth";

export default function Login() {
    const [pwd, setPwd] = useState<string>("");
    const { connect, connecting } = useApi();

    return <div className="App">
        <Grid container spacing={0} direction="column" alignItems="center" justifyContent="center" minHeight="100%">
            <Card variant="outlined" style={{maxWidth: '20em'}}>
                <CardContent style={{display: 'flex', flexDirection: 'column', alignItems: 'center'}}>
                    <Typography sx={{ fontSize: 20 }} variant="h1" color="text.secondary" gutterBottom>PartyHall Admin</Typography>
                    <TextField sx={{pt: 2, pb: 2}} type="password" value={pwd} onChange={e => setPwd(e.target.value)} />
                </CardContent>
                <CardActions>
                    <Button style={{width: '100%'}} size="small" onClick={() => connect(pwd)} disabled={connecting}>Login</Button>
                </CardActions>
            </Card>
        </Grid>
    </div>;
}