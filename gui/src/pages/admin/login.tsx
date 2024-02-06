import { Button, Card, CardActions, CardContent, Grid, Input, Switch, Typography } from "@mui/material";
import { Controller, useForm } from "react-hook-form";
import { useApi } from "../../hooks/useApi";
import { useState } from "react";

export default function Login() {
    const { login, loginAsGuest } = useApi();
    const [loginAsUser, setLoginAsUser] = useState<boolean>(false);
    const { handleSubmit, control } = useForm({
        defaultValues: {
            username: localStorage.getItem('username') || '',
            password: '',
        }
    });

    const onSubmit = (data: any) => {
        if (loginAsUser) {
            login(data.username, data.password);
        } else {
            loginAsGuest(data.username);
        }
    };

    return <Grid container spacing={0} direction="column" alignItems="center" justifyContent="center" minHeight="100%">
        <form onSubmit={handleSubmit(onSubmit)}>
            <Card variant="outlined" style={{ maxWidth: '20em' }}>
                <CardContent style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <Typography sx={{ fontSize: 20 }} variant="h1" color="text.secondary" gutterBottom>PartyHall</Typography>
                    {
                        loginAsUser && <>
                            <Controller
                                name="username"
                                control={control}
                                render={({ field }) => <Input placeholder="Username" type="username" required {...field} />}
                            />
                            <Controller
                                name="password"
                                control={control}
                                render={({ field }) => <Input placeholder="Password" type="password" required {...field} />}
                            />
                        </>
                    }
                    {
                        !loginAsUser && <>
                            <Controller
                                name="username"
                                control={control}
                                render={({ field }) => <Input placeholder="Name" type="username" required {...field} />}
                            />
                        </>
                    }
                </CardContent>
                <CardActions>
                    <Switch onChange={(_, x) => {
                        setLoginAsUser(x)
                    }} value={loginAsUser} />
                    <Button style={{ width: '100%' }} size="small" type="submit" variant="outlined">Login</Button>
                </CardActions>
            </Card>
        </form>
    </Grid>;
}