import { Flex } from 'antd';
import { ReactNode } from 'react';

export default function KeyVal({ label, children }: { label: string; children: ReactNode }) {
    return (
        <Flex justify="space-between" align="center" gap="2em">
            <span className="red">{label}: </span>
            <Flex>{children}</Flex>
        </Flex>
    );
}
