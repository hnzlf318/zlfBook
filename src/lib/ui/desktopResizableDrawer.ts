import { ref, onBeforeUnmount } from 'vue';

export interface ResizableNavigationDrawerWidthOptions {
    storageKey: string;
    defaultWidth: number;
    minWidth: number;
    /**
     * Max width in px. If not provided, only `maxWidthRatio` (if any) is used.
     */
    maxWidth?: number;
    /**
     * Max width as a ratio of viewport width (0-1). Example: 0.6 means max 60% of viewport width.
     */
    maxWidthRatio?: number;
}

export interface ResizableNavigationDrawerWidthResult {
    drawerWidth: ReturnType<typeof ref<number>>;
    onResizePointerDown: (e: PointerEvent) => void;
}

function readNumberFromLocalStorage(key: string): number | null {
    try {
        const raw = localStorage.getItem(key);
        if (!raw) return null;
        const n = parseFloat(raw);
        return Number.isFinite(n) ? n : null;
    } catch {
        return null;
    }
}

function writeNumberToLocalStorage(key: string, value: number): void {
    try {
        localStorage.setItem(key, String(value));
    } catch {
        // ignore
    }
}

function clamp(n: number, min: number, max: number): number {
    return Math.max(min, Math.min(max, n));
}

function computeMaxWidth(options: ResizableNavigationDrawerWidthOptions): number {
    const ratioMax = options.maxWidthRatio && options.maxWidthRatio > 0 && options.maxWidthRatio <= 1
        ? Math.floor(window.innerWidth * options.maxWidthRatio)
        : Number.POSITIVE_INFINITY;
    const pxMax = options.maxWidth && options.maxWidth > 0 ? options.maxWidth : Number.POSITIVE_INFINITY;
    return Math.max(options.minWidth, Math.min(ratioMax, pxMax));
}

export function useResizableNavigationDrawerWidth(options: ResizableNavigationDrawerWidthOptions): ResizableNavigationDrawerWidthResult {
    const initial = readNumberFromLocalStorage(options.storageKey);
    const drawerWidth = ref<number>(Number.isFinite(initial as number) && (initial as number) > 0 ? (initial as number) : options.defaultWidth);

    let dragging = false;
    let startX = 0;
    let startWidth = 0;
    let rafId = 0;

    function normalizeAndPersist(): void {
        const maxW = computeMaxWidth(options);
        const next = clamp(drawerWidth.value, options.minWidth, maxW);
        drawerWidth.value = next;
        writeNumberToLocalStorage(options.storageKey, next);
    }

    function stopDragging(): void {
        if (!dragging) return;
        dragging = false;
        document.body.classList.remove('ez-resizing');
        window.removeEventListener('pointermove', onPointerMove);
        window.removeEventListener('pointerup', stopDragging);
        window.removeEventListener('pointercancel', stopDragging);
        normalizeAndPersist();
    }

    function onPointerMove(e: PointerEvent): void {
        if (!dragging) return;
        const dx = e.clientX - startX;
        const maxW = computeMaxWidth(options);
        const next = clamp(startWidth + dx, options.minWidth, maxW);

        // Throttle UI updates to animation frames
        if (rafId) cancelAnimationFrame(rafId);
        rafId = requestAnimationFrame(() => {
            drawerWidth.value = next;
        });
    }

    function onResizePointerDown(e: PointerEvent): void {
        // Left button only
        if (e.button !== 0) return;
        dragging = true;
        startX = e.clientX;
        startWidth = drawerWidth.value;
        document.body.classList.add('ez-resizing');

        window.addEventListener('pointermove', onPointerMove);
        window.addEventListener('pointerup', stopDragging);
        window.addEventListener('pointercancel', stopDragging);

        e.preventDefault();
    }

    function onWindowResize(): void {
        // Ensure width stays within bounds after viewport resize
        const maxW = computeMaxWidth(options);
        drawerWidth.value = clamp(drawerWidth.value, options.minWidth, maxW);
    }

    window.addEventListener('resize', onWindowResize);

    // Initial clamp + persist once to ensure valid value
    normalizeAndPersist();

    onBeforeUnmount(() => {
        if (rafId) cancelAnimationFrame(rafId);
        stopDragging();
        window.removeEventListener('resize', onWindowResize);
    });

    return {
        drawerWidth,
        onResizePointerDown
    };
}


