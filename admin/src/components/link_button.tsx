import { Link } from "react-router-dom";

/**
 * Workaround the fact that antd is a fucking crappy library
 * Just let me do the same that I do in Material UI ffs
 */

export default function LinkButton({ to, children }: { to: string, children: any }) {
    return <Link to={to}>
        {children}
    </Link>
}