import '../assets/css/loader.scss';
import { ReactNode } from 'react';

export default function Loader({ children, loading }: { children: ReactNode; loading: boolean }) {
    if (!loading) {
        return <>{children}</>;
    }

    return <div className="css-loader"></div>;
}
