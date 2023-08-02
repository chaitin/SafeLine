import ReactDOM, { type Root } from "react-dom";

const MARK = "__ct_react_root__";

type ContainerType = (Element | DocumentFragment) & {
  [MARK]?: Root;
};

export function render(node: React.ReactElement, container: ContainerType) {
  const root = container[MARK] || container;

  ReactDOM.render(node, root);

  container[MARK] = root;
}

export async function unmount(container: ContainerType) {
  return Promise.resolve().then(() => {
    container[MARK]?.unmount();
    delete container[MARK];
  });
}

export function sizeLength(l: number) {
  return l > 1024 * 2 ? Math.round(l / 1024) + "KB" : l + "B";
}

export function sampleLength(s: string) {
  const l = new Blob([s]).size;
  return l > 1024 * 2 ? Math.round(l / 1024) + "KB" : l + "B";
}

export function sampleSummary(s: string) {
  return s.split("\n").slice(0, 2).join(" ").slice(0, 60);
}
