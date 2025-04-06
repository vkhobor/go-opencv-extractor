export function undefToErr<T>(promise: Promise<T>) {
    return promise.then((value) => {
        if (value === undefined) {
            throw new Error('Value is undefined');
        }
        if (value === null) {
            throw new Error('Value is null');
        }
        return value;
    });
}
