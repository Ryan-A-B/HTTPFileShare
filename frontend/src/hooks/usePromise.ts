import React from 'react';

interface ActionPending {
    type: 'PENDING';
}

interface ActionSuccess<T> {
    type: 'SUCCESS';
    data: T;
}

interface ActionError {
    type: 'ERROR';
    error: Error;
}

type Action<T> = ActionPending | ActionSuccess<T> | ActionError;

interface StatePending {
    type: 'PENDING';
}

interface StateLoaded<T> {
    type: 'LOADED';
    data: T;
}

interface StateError {
    type: 'ERROR';
    error: Error;
}

type State<T> = StatePending | StateLoaded<T> | StateError;

function reducer<T>(state: State<T>, action: Action<T>): State<T> {
    switch (action.type) {
        case 'PENDING':
            return { type: 'PENDING' };
        case 'SUCCESS':
            return { type: 'LOADED', data: action.data };
        case 'ERROR':
            return { type: 'ERROR', error: action.error };
    }
}

export default function usePromise<T>(promise: Promise<T>): State<T> {
    const [state, dispatch] = React.useReducer(reducer, { type: 'PENDING' } as State<T>);

    React.useEffect(() => {
        dispatch({ type: 'PENDING' });
        promise.then(
            (data) => dispatch({ type: 'SUCCESS', data }),
            (error) => dispatch({ type: 'ERROR', error }),
        );
    }, [promise]);

    return state as State<T>;
}