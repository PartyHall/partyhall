import { Button, Card, CardActions, CardContent, Grid, Input, Switch, Typography } from "@mui/material";
import { Controller, useForm } from "react-hook-form";
import { useApi } from "../../hooks/useApi";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useSnackbar } from "../../hooks/snackbar";

export default function Login() {
    const {t} = useTranslation();
    const {showSnackbar} = useSnackbar();
    const { login, loginAsGuest } = useApi();
    const [loginAsUser, setLoginAsUser] = useState<boolean>(false);
    const { handleSubmit, control } = useForm({
        defaultValues: {
            username: localStorage.getItem('username') || '',
            password: '',
        }
    });

    const onSubmit = (data: any) => {
        try {
            if (loginAsUser) {
                login(data.username, data.password);
            } else {
                loginAsGuest(data.username);
            }
        } catch (e) {
            showSnackbar(t('login.failed') + ': ' + e, 'error');
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
                                render={({ field }) => <Input placeholder={t('login.username')} type="username" required {...field} />}
                            />
                            <Controller
                                name="password"
                                control={control}
                                render={({ field }) => <Input placeholder="Password" type={t('login.password')} required {...field} />}
                            />
                        </>
                    }
                    {
                        !loginAsUser && <>
                            <Controller
                                name="username"
                                control={control}
                                render={({ field }) => <Input placeholder={t('login.name')} type="username" required {...field} />}
                            />
                        </>
                    }
                </CardContent>
                <CardActions>
                    <Switch onChange={(_, x) => {
                        setLoginAsUser(x)
                    }} value={loginAsUser} />
                    <Button style={{ width: '100%' }} size="small" type="submit" variant="outlined">{t('login.bt')}</Button>
                </CardActions>
            </Card>
        </form>
    </Grid>;
}