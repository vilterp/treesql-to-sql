/// <reference types="react-scripts" />

declare module "react-inspector" {
  interface ObjectInspectorProps {
    data: any;
  }
  export class ObjectInspector extends React.Component<ObjectInspectorProps> {}
}
