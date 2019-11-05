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
import { formatSpan } from "./format";
import { SchemaView } from "./SchemaView";

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
        render={({ data: schemaDesc }) => (
          <Form<AppState, Resp, ValidationResp>
            initialState={{
              cursorPos: { Line: 1, Col: 1, Offset: 1 },
              query: ``,
            }}
            submit={runQuery}
            validate={st =>
              validateQuery({
                CursorPos: st.cursorPos,
                Query: st.query,
              })
            }
            render={({ state, validationState, update, apiCallState }) => (
              <table>
                <tbody>
                  <tr style={{ verticalAlign: "top" }}>
                    <td>
                      <h2>Schema</h2>
                      <SchemaView
                        schema={schemaDesc}
                        highlighted={
                          validationState.tag === "VALIDATED"
                            ? validationState.resp.HighlightedElement
                            : null
                        }
                      />
                    </td>
                    <td style={{ width: 500 }}>
                      <div style={{ border: "1px solid black" }}>
                        <AceEditor
                          value={state.query}
                          height="200px"
                          onChange={value =>
                            update(st => ({ ...st, query: value }))
                          }
                          onCursorChange={evt => {
                            const cursor = evt.getCursor();
                            update(st => ({
                              ...st,
                              cursorPos: {
                                Line: cursor.row + 1,
                                Col: cursor.column + 1,
                                Offset: 1, // TODO(vilterp): do something about offset
                              },
                            }));
                          }}
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

                      {validationState.tag === "VALIDATED" ||
                      validationState.tag === "VALIDATING" ? (
                        <>
                          <h2>Completions</h2>
                          <ul>
                            {(
                              (validationState.resp
                                ? validationState.resp.Completions
                                : []) || []
                            ).map((comp, idx) => (
                              <li key={idx}>
                                {comp.Kind}: {comp.Content}
                              </li>
                            ))}
                          </ul>
                          <h2>Errors</h2>
                          <ul>
                            {/* TODO(vilterp): de-kludge this */}
                            {(
                              (validationState.resp
                                ? validationState.resp.Errors
                                : []) || []
                            ).map((err, idx) => (
                              <li key={idx}>
                                {formatSpan(err.Span)}: {err.Message}
                              </li>
                            ))}
                          </ul>
                          <h2>Parse Error</h2>
                          {validationState.resp &&
                          validationState.resp.ParseError
                            ? validationState.resp.ParseError
                            : null}
                        </>
                      ) : null}

                      {apiCallState.tag === State.FAILED ? (
                        <Alert
                          type={AlertType.DANGER}
                          title="Error"
                          message={apiCallState.error}
                        />
                      ) : null}
                    </td>
                    <td>
                      {apiCallState.tag === State.SUCCEEDED ? (
                        <>
                          <pre>{apiCallState.response.SQL}</pre>
                          <AceEditor
                            value={JSON.stringify(
                              JSON.parse(apiCallState.response.Res),
                              null,
                              2,
                            )}
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
