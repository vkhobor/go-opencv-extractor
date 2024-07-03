export function undefToErr<T>(promise: Promise<T>) {
  return promise.then((value) => {
    if (value === undefined) {
      throw new Error('Value is undefined');
    }
    return value;
  });
}
