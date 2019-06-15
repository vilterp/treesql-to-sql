import * as React from "react";
import { HighlightedElement, SchemaDesc } from "./api";

export function SchemaView(props: {
  schema: SchemaDesc;
  highlighted: HighlightedElement | null;
}) {
  const schemaDesc = props.schema;
  const path = props.highlighted
    ? extractPath(props.highlighted.Path)
    : EMPTY_PATH;
  return (
    <>
      <ul>
        {Object.keys(schemaDesc.Tables)
          .sort()
          .map(tableName => (
            <li key={tableName}>
              <span
                style={{
                  textDecoration:
                    tableName === path.tableName && path.colName === null
                      ? "underline"
                      : "none",
                }}
              >
                {tableName}
              </span>
              <ul>
                {Object.keys(schemaDesc.Tables[tableName].Columns)
                  .sort()
                  .map(colName => (
                    <li
                      key={colName}
                      style={{
                        textDecoration:
                          tableName === path.tableName &&
                          colName === path.colName
                            ? "underline"
                            : "none",
                      }}
                    >
                      {colName}
                    </li>
                  ))}
              </ul>
            </li>
          ))}
      </ul>
    </>
  );
}

interface Path {
  tableName: string | null;
  colName: string | null;
}

const EMPTY_PATH = {
  colName: null,
  tableName: null,
};

function extractPath(path: string): Path {
  const segments = path.split("/");
  if (segments.length == 1) {
    return {
      tableName: segments[0],
      colName: null,
    };
  } else if (segments.length == 2) {
    return {
      tableName: segments[0],
      colName: segments[1],
    };
  }
  return EMPTY_PATH;
}
