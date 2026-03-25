const DB_NAME = 'ezbookkeeping_transaction_picture_cache';
const DB_VERSION = 1;
const STORE_NAME = 'pictures';

type PictureCacheRecord = {
    readonly pictureId: string;
    readonly extension?: string;
    readonly blob: Blob;
};

let dbPromise: Promise<IDBDatabase> | undefined;

function openDb(): Promise<IDBDatabase> {
    if (dbPromise) {
        return dbPromise;
    }

    dbPromise = new Promise<IDBDatabase>((resolve, reject) => {
        if (!('indexedDB' in window)) {
            reject(new Error('IndexedDB is not supported in this environment'));
            return;
        }

        const request = window.indexedDB.open(DB_NAME, DB_VERSION);

        request.onupgradeneeded = () => {
            const db = request.result;

            if (!db.objectStoreNames.contains(STORE_NAME)) {
                db.createObjectStore(STORE_NAME, { keyPath: 'pictureId' });
            }
        };

        request.onsuccess = () => {
            resolve(request.result);
        };

        request.onerror = () => {
            reject(request.error);
        };
    });

    return dbPromise;
}

export async function getTransactionPictureBlobFromLocalCache(pictureId: string): Promise<Blob | undefined> {
    if (!pictureId) {
        return undefined;
    }

    const db = await openDb();

    return new Promise<Blob | undefined>((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readonly');
        const store = tx.objectStore(STORE_NAME);
        const req = store.get(pictureId);

        req.onsuccess = () => {
            const record = req.result as PictureCacheRecord | undefined;
            resolve(record?.blob);
        };

        req.onerror = () => {
            reject(req.error);
        };
    });
}

export async function setTransactionPictureBlobToLocalCache(pictureId: string, blob: Blob, extension?: string): Promise<void> {
    if (!pictureId || !blob) {
        return;
    }

    const db = await openDb();

    const record: PictureCacheRecord = {
        pictureId,
        extension,
        blob
    };

    return new Promise<void>((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite');
        const store = tx.objectStore(STORE_NAME);
        const req = store.put(record);

        req.onsuccess = () => {
            resolve();
        };

        req.onerror = () => {
            reject(req.error);
        };
    });
}

export async function deleteTransactionPictureBlobFromLocalCache(pictureId: string): Promise<void> {
    if (!pictureId) {
        return;
    }

    const db = await openDb();

    return new Promise<void>((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite');
        const store = tx.objectStore(STORE_NAME);
        const req = store.delete(pictureId);

        req.onsuccess = () => {
            resolve();
        };

        req.onerror = () => {
            reject(req.error);
        };
    });
}

