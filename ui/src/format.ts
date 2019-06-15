import { SourcePosition, Span } from "./api";

export function formatPos(pos: SourcePosition) {
  return `${pos.Line}:${pos.Col}`;
}

export function formatSpan(span: Span) {
  return `[${formatPos(span.From)} - ${formatPos(span.To)}]`;
}
