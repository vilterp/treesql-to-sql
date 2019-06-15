import * as React from "react";

export enum AlertType {
  PRIMARY = "primary",
  SUCCESS = "success",
  WARNING = "warning",
  DANGER = "danger",
}

export function Alert(props: {
  type: AlertType;
  title: string;
  message: string;
}) {
  return (
    <div style={{ color: "red" }}>
      {props.title}: {props.message}
    </div>
  );
}
