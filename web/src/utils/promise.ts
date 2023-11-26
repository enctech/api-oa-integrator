type Result<T, E> = { result?: T; error?: E };

export function resolvePromise<T, E = any>(
  promise: Promise<T>,
): Promise<Result<T, E>> {
  return promise.then((result) => ({ result })).catch((error) => ({ error }));
}
