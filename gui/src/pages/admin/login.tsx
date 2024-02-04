import { Button, Card, CardActions, CardContent, Grid, Input, Typography } from "@mui/material";
import { Controller, useForm } from "react-hook-form";
import { useApi } from "../../hooks/useApi";

export default function Login() {
    const {login} = useApi();
    const { handleSubmit, control } = useForm({
        defaultValues: {
            username: localStorage.getItem('username') || '',
            password: '',
        }
    });

    const onSubmit = (data: any) => login(data.username, data.password);

    return <Grid container spacing={0} direction="column" alignItems="center" justifyContent="center" minHeight="100%">
        <form onSubmit={handleSubmit(onSubmit)}>
            <Card variant="outlined" style={{ maxWidth: '20em' }}>
                <CardContent style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <Typography sx={{ fontSize: 20 }} variant="h1" color="text.secondary" gutterBottom>PartyHall Admin</Typography>
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

                </CardContent>
                <CardActions>
                    <Button style={{ width: '100%' }} size="small" type="submit" variant="outlined">Login</Button>
                </CardActions>
            </Card>
        </form>
    </Grid>;
}