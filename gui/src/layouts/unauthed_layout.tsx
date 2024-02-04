import { Navigate, useOutlet } from "react-router-dom";
import { useApi } from "../hooks/useApi";

export default function UnauthedLayout() {
    const outlet = useOutlet();
    const { socketMode, isLoggedIn } = useApi();

    if (socketMode != 'admin') {
        return <Navigate to={"/"} />
    }

    if (!!isLoggedIn()) {
        return <Navigate to={"/admin"} />
    }

    return <>
        {outlet}
    </>
}