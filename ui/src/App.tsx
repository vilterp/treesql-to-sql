import React from "react";
import AceEditor from "react-ace";
import "brace";
import "brace/mode/json";
import "./App.css";
import Form, { SubmitButton } from "./util/Form";
import { State } from "./util/apiCallState";
import {
  getSchema,
  Resp,
  runQuery,
  SourcePosition,
  validateQuery,
  ValidationResp,
} from "./api";
import Load from "./util/Load";
import { Alert, AlertType } from "./util/alert";

interface AppState {
  cursorPos: SourcePosition;
  query: string;
}

function App() {
  return (
    <>
      <h1>TreeSQL Console</h1>
      <Load
        load={getSchema}
        render={schemaDesc => (
          <Form<AppState, Resp, ValidationResp>
            initialState={{
              cursorPos: { Line: 1, Col: 1, Offset: 1 },
              query: `MANY posts { id, body}`,
            }}
            submit={runQuery}
            validate={st =>
              validateQuery({
                CursorPos: st.cursorPos,
                Query: st.query,
              })
            }
            render={({ state, update, apiCallState, validationState }) => (
              <table>
                <tbody>
                  <tr>
                    <td>
                      <div style={{ border: "1px solid black" }}>
                        <AceEditor
                          value={state.query}
                          height="300px"
                          onChange={value =>
                            update(st => ({ ...st, query: value }))
                          }
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

                      {apiCallState.tag === State.FAILED ? (
                        <Alert
                          type={AlertType.DANGER}
                          title="Error"
                          message={apiCallState.error}
                        />
                      ) : null}

                      {apiCallState.tag === State.SUCCEEDED ? (
                        <>
                          <pre>{apiCallState.response.SQL}</pre>
                          <AceEditor
                            value={JSON.stringify(
                              JSON.parse(apiCallState.response.Res),
                              null,
                              2,
                            )}
                            onCursorChange={evt => {
                              const cursor = evt.getCursor();
                              update(st => ({
                                ...st,
                                cursorPos: {
                                  Line: cursor.row,
                                  Col: cursor.column + 1,
                                  Offset: 1, // TODO(vilterp): do something about offset
                                },
                              }));
                            }}
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
                      ) : null}
                    </td>
                    <td>
                      <h2>Schema</h2>
                      {JSON.stringify(schemaDesc, null, 2)}
                      <h2>Errors</h2>
                      {validationState.tag === "VALIDATED"
                        ? JSON.stringify(validationState.resp.Errors, null, 2)
                        : null}
                      <h2>Completions</h2>
                      {validationState.tag === "VALIDATED"
                        ? JSON.stringify(
                            validationState.resp.Completions,
                            null,
                            2,
                          )
                        : null}
                      <h2>Parse Error</h2>
                      {validationState.tag === "VALIDATED"
                        ? JSON.stringify(
                            validationState.resp.ParseError,
                            null,
                            2,
                          )
                        : null}
                    </td>
                  </tr>
                </tbody>
              </table>
            )}
          />
        )}
      />
    </>
  );
}

export default App;
