declare module 'cdgraphics' {
    class CDGContext {
        readonly WIDTH: number;
        readonly HEIGHT: number;
        readonly DISPLAY_WIDTH: number;
        readonly DISPLAY_HEIGHT: number;
        readonly DISPLAY_BOUNDS: [number, number, number, number];
        readonly TILE_WIDTH: number;
        readonly TILE_HEIGHT: number;

        hOffset: number;
        vOffset: number;
        keyColor: number | null;
        bgColor: number | null;
        borderColor: number | null;
        clut: [number, number, number][];
        pixels: Uint8ClampedArray;
        buffer: Uint8ClampedArray;
        imageData: ImageData;
        backgroundRGBA: [number, number, number, number];
        contentBounds: [number, number, number, number];

        constructor();
        init(): void;
        setCLUTEntry(index: number, r: number, g: number, b: number): void;
        renderFrame(options?: { forceKey?: boolean }): void;
    }

    interface CDGInstruction {
        execute(ctx: CDGContext): void;
    }

    class CDGParser {
        readonly COMMAND_MASK: number;
        readonly CDG_COMMAND: number;
        readonly BY_TYPE: Record<
            number,
            new (bytes: Uint8Array) => CDGInstruction
        >;

        bytes: Uint8Array;
        numPackets: number;
        pc: number;

        constructor(buffer: ArrayBuffer);
        parseThrough(
            sec: number
        ): (CDGInstruction & { isRestarting?: boolean })[];
        parse(packet: Uint8Array): CDGInstruction | false;
    }

    interface RenderResult {
        imageData: ImageData;
        isChanged: boolean;
        backgroundRGBA: [number, number, number, number];
        contentBounds: [number, number, number, number];
    }

    class CDGPlayer {
        ctx: CDGContext;
        private parser: CDGParser | null;
        private forceKey: boolean | null;

        constructor();
        load(buffer: ArrayBuffer): void;
        render(time: number, opts?: { forceKey?: boolean }): RenderResult;
    }

    export = CDGPlayer;
}
