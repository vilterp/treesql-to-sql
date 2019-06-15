export interface Resp {
  Res: string;
  SQL: string;
}

export interface Req {
  query: string;
}

interface ValidationsReq {
  Query: string;
  CursorPos: SourcePosition;
}

export interface ValidationResp {
  Completions: {
    Kind: string;
    Content: string;
  }[];
  Errors: {
    Span: Span;
    Message: string;
  };
  ParseError: string;
}

export interface Span {
  From: SourcePosition;
  To: SourcePosition;
}

export interface SourcePosition {
  Line: number;
  Col: number;
  Offset: number;
}

export interface SchemaDesc {
  Tables: { [name: string]: TableDesc };
}

export interface TableDesc {
  Columns: { [name: string]: ColumnDesc };
}

export interface ColumnDesc {}

export function runQuery(req: Req): Promise<Resp> {
  return fetch("/query", {
    method: "POST",
    body: req.query,
  }).then(res => {
    if (res.status !== 200) {
      return res.text().then(t => Promise.reject(t));
    }
    return res.json();
  });
}

export function getSchema(): Promise<SchemaDesc> {
  return fetch("/schema").then(res => {
    if (res.status !== 200) {
      return res.text().then(t => Promise.reject(t));
    }
    return res.json();
  });
}

export function validateQuery(req: ValidationsReq): Promise<ValidationResp> {
  return fetch("/validate", {
    method: "POST",
    body: JSON.stringify(req),
  }).then(res => {
    if (res.status !== 200) {
      return res.text().then(t => Promise.reject(t));
    }
    return res.json();
  });
}
