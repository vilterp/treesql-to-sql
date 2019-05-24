import React, { Component } from "react";
import { ObjectInspector } from "react-inspector";
import { DateTime, Duration, DurationUnit } from "luxon";

import "./App.css";

interface ConsoleState {
  query: string;
  // TODO(vilterp): unify these into an API call state thing
  loading: boolean;
  loadErr: string | null;
  response: string | null;
  responseTime: Duration | null;
}

class App extends Component<{}, ConsoleState> {
  state: ConsoleState = {
    query: `MANY clusters { id, name }`,
    loading: false,
    loadErr: null,
    response: null,
    responseTime: null,
  };

  handleUpdateQuery = (newQuery: string) => {
    this.setState(prevState => ({
      ...prevState,
      query: newQuery,
    }));
  };

  handleSubmit = () => {
    this.setState(prevState => ({
      ...prevState,
      loading: true,
    }));
    const startTime = DateTime.local();
    fetch(`http://${window.location.host}/query`, {
      method: "POST",
      body: this.state.query,
    })
      .then(res => {
        if (res.status !== 200) {
          res.text().then(text => {
            this.setState(prevState => ({
              ...prevState,
              loading: false,
              loadErr: `${text}`,
              response: null,
            }));
          });
          return;
        }

        res.json().then(jsonRes => {
          this.setState(prevState => ({
            ...prevState,
            loading: false,
            loadErr: null,
            response: jsonRes,
            responseTime: DateTime.local().diff(startTime),
          }));
        });
      })
      .catch(err => {
        this.setState(prevState => ({
          ...prevState,
          loading: false,
          loadErr: err.toString(),
        }));
      });
  };

  renderError = () => {
    if (!this.state.loadErr) {
      return null;
    }

    return <div style={{ backgroundColor: "pink" }}>{this.state.loadErr}</div>;
  };

  renderResponse = () => {
    if (this.state.loading) {
      return "Loading...";
    }

    if (!this.state.response) {
      return null;
    }

    return (
      <>
        {this.state.responseTime
          ? // jesus christ Luxon; why can't you be the same as Go's time.Duration??
            `${this.state.responseTime.as("milliseconds")}ms`
          : null}
        <ObjectInspector data={this.state.response} />
      </>
    );
  };

  render() {
    return (
      <div className="App">
        <h1>Console</h1>
        <textarea
          style={{ fontFamily: "monospace" }}
          rows={20}
          cols={50}
          value={this.state.query}
          onChange={evt => this.handleUpdateQuery(evt.target.value)}
        />
        <br />
        <button onClick={this.handleSubmit} disabled={this.state.loading}>
          Run
        </button>
        <br />
        {this.renderError()}
        {this.renderResponse()}
      </div>
    );
  }
}

export default App;
