export enum State {
  NOT_TRIED_YET,
  LOADING,
  SUCCEEDED,
  FAILED,
}

interface NotTriedYet {
  tag: State.NOT_TRIED_YET;
}

interface Loading {
  tag: State.LOADING;
}

interface Succeeded<R> {
  tag: State.SUCCEEDED;
  response: R;
}

interface Failed<E> {
  tag: State.FAILED;
  error: E;
}

export type APICallState<R, E> =
  | NotTriedYet
  | Loading
  | Succeeded<R>
  | Failed<E>;
