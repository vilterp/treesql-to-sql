import React, { Component } from "react";
import { DateTime, Duration } from "luxon";
import AceEditor from "react-ace";
import "brace";
import "brace/mode/json";
import "./App.css";

interface ConsoleState {
  query: string;
  // TODO(vilterp): unify these into an API call state thing
  loading: boolean;
  loadErr: string | null;
  response: Resp | null;
  responseTime: Duration | null;
}

interface Resp {
  Res: string;
  SQL: string;
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
            query: jsonRes.FormattedTreeSQL,
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

  renderEditor = () => {
    return (
      <>
        <div style={{ border: "1px solid black" }}>
          <AceEditor
            value={this.state.query}
            height="300px"
            onChange={value => this.handleUpdateQuery(value)}
            highlightActiveLine={false}
            showGutter={false}
            setOptions={{
              showLineNumbers: false,
              highlightGutterLine: false,
              tabSize: 2,
            }}
          />
        </div>
        <br />
        <button onClick={this.handleSubmit} disabled={this.state.loading}>
          Run
        </button>
      </>
    );
  };

  renderError = () => {
    if (!this.state.loadErr) {
      return null;
    }

    return (
      <div style={{ backgroundColor: "pink", fontFamily: "monospace" }}>
        {this.state.loadErr}
      </div>
    );
  };

  renderResponse = () => {
    if (this.state.loading) {
      return "Loading...";
    }

    if (!this.state.response) {
      return null;
    }

    const stringified = JSON.stringify(
      JSON.parse(this.state.response.Res),
      null,
      2,
    );

    return (
      <>
        {this.state.responseTime
          ? // jesus christ Luxon; why can't you be the same as Go's time.Duration??
            `${this.state.responseTime.as("milliseconds")}ms`
          : null}
        <pre>{this.state.response.SQL}</pre>
        {/*<JSONViewer json={JSON.parse(this.state.response.Res)} />*/}
        {/*<ObjectInspector data={JSON.parse(this.state.response.Res)} />*/}
        {/*<pre>*/}
        {/*  {JSON.stringify(JSON.parse(this.state.response.Res), null, 2)}*/}
        {/*</pre>*/}
        <AceEditor
          value={stringified}
          mode="json"
          readOnly={true}
          maxLines={Infinity}
          highlightActiveLine={false}
          setOptions={{
            showLineNumbers: false,
            highlightGutterLine: false,
          }}
        />
      </>
    );
  };

  render() {
    return (
      <div className="App">
        <h1>TreeSQL Console</h1>
        {this.renderEditor()}
        <br />
        {this.renderError()}
        {this.renderResponse()}
      </div>
    );
  }
}

export default App;
