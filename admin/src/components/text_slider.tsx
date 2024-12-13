import { Flex, Slider } from "antd";
import { ReactNode } from "react";

type Props = {
    leftText?: ReactNode;
    rightText?: ReactNode;
    min?: number;
    max?: number;
    value: number;
    onChange?: (val: number) => void;
    disabled?: boolean;
    className?: string;
}

export default function TextSlider({ leftText, rightText, min, max, value, onChange, disabled, className }: Props) {
    return <Flex align="center">
        {leftText}
        <Slider
            min={min ?? 0}
            max={max ?? 100}
            value={value}
            onChange={onChange}
            className={`${className} flex1`}
            disabled={disabled !== undefined ? disabled : false}
        />
        {rightText}
    </Flex>
}