import PlaceHolderCover from '../assets/placeholder_cover.webp';

type Props = {
    hasImage?: boolean;
    src: string;
    alt: string;
    className?: string;
};

export default function Image({ hasImage, className, src, alt }: Props) {
    return <img className={className} src={!hasImage && hasImage !== undefined ? PlaceHolderCover : src} alt={alt} />;
}
