import * as React from "react";
import { FormEvent } from "react";
import { APICallState, State } from "./apiCallState";

export type UpdateFn<State> = (s: State) => State;
export type UpdaterFn<State> = (updater: UpdateFn<State>) => void;

// TODO(vilterp): some way of focusing on mount

interface FormProps<State, Response, ValidationResponse> {
  render: (
    props: FormRenderProps<State, Response, ValidationResponse>,
  ) => React.ReactNode;
  submit: (state: State) => Promise<Response>;
  validate?: (state: State) => Promise<ValidationResponse>;
  initialState: State;
  onSucceed?: (
    r: Response,
    onUpdate: (updater: UpdateFn<State>) => void,
  ) => void;
}

interface FormState<State, Response, ValidationResponse> {
  formState: State;
  validationState: ValidationState<ValidationResponse>; // TODO: track whether this is up to date
  apiCallState: APICallState<Response, string>;
}

type ValidationState<VResp> =
  | { tag: "NEVER_VALIDATED" }
  | { tag: "VALIDATING"; resp: VResp | null }
  | { tag: "VALIDATED"; resp: VResp };

interface FormRenderProps<State, Response, ValidationResponse> {
  state: State;
  update: UpdaterFn<State>;
  apiCallState: APICallState<Response, string>;
  validationState: ValidationState<ValidationResponse>;
}

export default class Form<S, R, V = {}> extends React.Component<
  FormProps<S, R, V>,
  FormState<S, R, V>
> {
  constructor(props: FormProps<S, R, V>) {
    super(props);
    this.state = {
      formState: props.initialState,
      validationState: { tag: "NEVER_VALIDATED" },
      apiCallState: { tag: State.NOT_TRIED_YET },
    };
  }

  componentDidMount(): void {
    this.doValidation(this.state.formState);
  }

  handleFormUpdate = (updater: UpdateFn<S>) => {
    const newState = updater(this.state.formState);
    this.setState({
      ...this.state,
      formState: newState,
    });

    this.doValidation(newState);
  };

  doValidation(newState: S) {
    if (this.props.validate) {
      this.setState({
        validationState: {
          tag: "VALIDATING",
          resp:
            this.state.validationState.tag === "VALIDATED"
              ? this.state.validationState.resp
              : null,
        },
      });
      this.props.validate(newState).then(valResp => {
        this.setState({
          validationState: { tag: "VALIDATED", resp: valResp },
        });
      });
    }
  }

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
          validationState: this.state.validationState,
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
