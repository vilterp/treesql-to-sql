import * as React from "react";
import { Alert, AlertType } from "./alert";
import { APICallState, State } from "./apiCallState";

import "./load.scss";

interface LoadProps<T> {
  load: () => Promise<T>;
  render: (props: LoadRenderProps<T>) => React.ReactNode;
}

interface LoadState<T> {
  apiCallState: APICallState<T, string>;
}

interface LoadRenderProps<T> {
  data: T;
  reload: () => void;
}

class Load<T> extends React.Component<LoadProps<T>, LoadState<T>> {
  state: LoadState<T> = {
    apiCallState: { tag: State.NOT_TRIED_YET },
  };

  componentDidMount() {
    this.reload();
  }

  reload = () => {
    this.setState({
      apiCallState: { tag: State.LOADING },
    });
    this.props.load().then(
      res => {
        this.setState({
          apiCallState: { tag: State.SUCCEEDED, response: res },
        });
      },
      err => {
        this.setState({
          apiCallState: { tag: State.FAILED, error: err },
        });
      },
    );
  };

  render() {
    switch (this.state.apiCallState.tag) {
      case State.LOADING:
      case State.NOT_TRIED_YET:
        return (
          <div className="loading">
            Loading...
          </div>
        );
      case State.FAILED:
        return (
          <Alert
            type={AlertType.DANGER}
            title="Load failed"
            message={this.state.apiCallState.error}
          />
        );
      case State.SUCCEEDED:
        return this.props.render({
          data: this.state.apiCallState.response,
          reload: this.reload,
        });
    }
  }
}

export default Load;
