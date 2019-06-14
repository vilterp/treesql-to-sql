import * as React from "react";
import { FormEvent } from "react";
import { APICallState, State } from "./apiCallState";

export type UpdateFn<State> = (s: State) => State;
export type UpdaterFn<State> = (updater: UpdateFn<State>) => void;

// TODO(vilterp): some way of focusing on mount

interface FormProps<State, Response> {
  render: (props: FormRenderProps<State, Response>) => React.ReactNode;
  submit: (state: State) => Promise<Response>;
  initialState: State;
  onSucceed?: (
    r: Response,
    onUpdate: (updater: UpdateFn<State>) => void,
  ) => void;
}

interface FormState<State, Response> {
  formState: State;
  apiCallState: APICallState<Response, string>;
}

interface FormRenderProps<State, Response> {
  state: State;
  update: UpdaterFn<State>;
  apiCallState: APICallState<Response, string>;
}

export default class Form<S, R> extends React.Component<
  FormProps<S, R>,
  FormState<S, R>
  > {
  constructor(props: FormProps<S, R>) {
    super(props);
    this.state = {
      formState: props.initialState,
      apiCallState: { tag: State.NOT_TRIED_YET },
    };
  }

  handleFormUpdate = (updater: UpdateFn<S>) => {
    this.setState(prevState => ({
      ...prevState,
      formState: updater(prevState.formState),
    }));
  };

  handleSubmit = (evt: FormEvent) => {
    evt.preventDefault();
    this.setState(prevState => ({
      ...prevState,
      apiCallState: {
        tag: State.LOADING,
      },
    }));
    this.props.submit(this.state.formState).then(
      res => {
        this.setState(prevState => ({
          ...prevState,
          apiCallState: { tag: State.SUCCEEDED, response: res },
        }));
        if (this.props.onSucceed) {
          this.props.onSucceed(res, this.handleFormUpdate);
        }
      },
      err => {
        this.setState(prevState => ({
          // TODO(vilterp): genericize error stringification
          ...prevState,
          apiCallState: { tag: State.FAILED, error: err },
        }));
      },
    );
  };

  render() {
    return (
      <form method="post" onSubmit={this.handleSubmit}>
        {this.props.render({
          state: this.state.formState,
          update: this.handleFormUpdate,
          apiCallState: this.state.apiCallState,
        })}
      </form>
    );
  }
}

// SubmitButton is a small wrapper around a Button component, intended to be used for
// a submit button on a form. It's meant to encapsulate these two common behaviors:
// - button is greyed out when the form is in an invalid state (the `invalid` prop)
// - button is greyed out with different text ("e.g. Loading..."; the `loadingText` prop)
//   while the form is submitting.
export function SubmitButton<R, E>(props: {
  callState: APICallState<R, E>;
  text: string;
  loadingText: string;
  invalid?: boolean;
  className?: string;
}) {
  return (
    <button
      className={props.className}
      disabled={props.callState.tag === State.LOADING || props.invalid}
    >
      {props.callState.tag === State.LOADING ? props.loadingText : props.text}
    </button>
  );
}
