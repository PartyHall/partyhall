import '../../assets/photobooth.scss';

export default function Countdown({ seconds }: { seconds: number }) {
    return (
        <div className="modal countdown">
            <h1>{seconds}</h1>
        </div>
    );
}
