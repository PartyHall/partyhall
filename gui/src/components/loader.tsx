import { ReactNode } from "react";
import '../assets/css/loader.scss';

export default function Loader({ children, loading }: { children: ReactNode, loading: boolean }) {
    return <>
        {loading && <div className="lds-ripple"><div></div><div></div></div>}
        {!loading && children}
    </>
}

export function TextLoader({ children, loading, text }: { children?: ReactNode, loading: boolean, text?: string | null }) {
    if (!loading) {
        return <>{children}</>
    }

    return <div className="loaderParent">
        <div className="lds-ripple"><div></div><div></div></div>
        {text && <span>{text}</span>}
    </div>
}