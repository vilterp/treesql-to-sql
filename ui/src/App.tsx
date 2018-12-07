import React, { Component } from "react";

interface ConsoleState {
  query: string;
  loading: boolean;
  loadErr: string | null;
  response: string | null;
}

class App extends Component<{}, ConsoleState> {
  state: ConsoleState = {
    query: `SELECT json_agg(json_build_object('id', id, 'name', name)) FROM clusters`,
    loading: false,
    loadErr: null,
    response: null,
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
    fetch("http://localhost:9000/query", {
      method: "POST",
      body: this.state.query,
    })
      .then(res => {
        if (res.status !== 200) {
          this.setState(prevState => ({
            ...prevState,
            loading: false,
            loadErr: `Error: ${res.status}`,
            response: null,
          }));
          return;
        }

        res.json().then(jsonRes => {
          console.log("jsonRes:", jsonRes);
          this.setState(prevState => ({
            ...prevState,
            loading: false,
            loadErr: null,
            response: jsonRes,
          }));
        });
      })
      .catch(err => {
        console.error("error running query:", err);
      });
  };

  renderError = () => {
    if (!this.state.loadErr) {
      return null;
    }

    return (
      <div style={{ backgroundColor: "pink" }}>Error: {this.state.loadErr}</div>
    );
  };

  renderResponse = () => {
    if (!this.state.response) {
      return null;
    }

    return <pre>{JSON.stringify(this.state.response, null, 2)}</pre>;
  };

  render() {
    return (
      <div className="App">
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
        {this.renderError()}
        {this.renderResponse()}
      </div>
    );
  }
}

export default App;
