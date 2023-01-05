import React from 'react';

function useStorageState(storage: Storage, key: string, defaultValue: string): [string, (value: string) => void] {
    const [state, internalSetState] = React.useState<string>(() => {
        const value = storage.getItem(key);
        if (value === null) return defaultValue;
        return JSON.parse(value);
    });
    const setState = React.useCallback((value: string) => {
        storage.setItem(key, JSON.stringify(value));
        internalSetState(value);
    }, [internalSetState, storage, key]);
    return [state, setState];
}

export default useStorageState;