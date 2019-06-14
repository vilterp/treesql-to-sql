import React from "react";
import AceEditor from "react-ace";
import "brace";
import "brace/mode/json";
import "./App.css";
import Form, { SubmitButton } from "./util/Form";
import { State } from "./util/apiCallState";

interface Resp {
  Res: string;
  SQL: string;
}

interface Req {
  query: string;
}

function runQuery(req: Req): Promise<Resp> {
  return fetch("/query", {
    method: "POST",
    body: req.query,
  })
    .then(res => {
      if (res.status !== 200) {
        return res.text().then(t => Promise.reject(t));
      }
      return res.json();
    });
}

function App() {
  // TODO(vilterp): load schema with <Load />
  return (
    <>
      <h1>TreeSQL Console</h1>
      <Form<Req, Resp>
        initialState={{ query: `MANY posts { id, body}` }}
        submit={runQuery}
        render={({ state, update, apiCallState }) => (
          <>
            <div style={{ border: "1px solid black" }}>
              <AceEditor
                value={state.query}
                height="300px"
                onChange={value => update(st => ({ ...st, query: value }))}
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
            <SubmitButton
              callState={apiCallState}
              text="Run"
              loadingText="Running..."
            />

            {apiCallState.tag === State.SUCCEEDED ?
              <>
                <pre>{apiCallState.response.SQL}</pre>
                <AceEditor
                  value={JSON.stringify(JSON.parse(apiCallState.response.Res), null, 2)}
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
              : null
            }
          </>
        )}
      />
    </>
  )
}

export default App;
