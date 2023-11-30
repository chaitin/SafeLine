import { createRoot } from 'react-dom/client';

const MARK = "__ct_react_root__";

type ContainerType = (Element | DocumentFragment) & {
  [MARK]?: any;
};

export function render(node: React.ReactElement, container: ContainerType) {
  const root = createRoot(container[MARK] || container);
  
  root.render(node);

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

export function formatNumberWithCommas(v: number) {
  return v.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

/* 将 时间戳 转换成 "Jul 23" 的格式 **/
export function formatDate(time: number): string {
  const date = new Date(time * 1000);
  const month = date.toLocaleString('en-US', { month: 'short' });
  const day = date.getDate();
  return `${month} ${day}`;
}

export function formatStarNumber(v: number) {
  if (v < 1000) {
    return v;
  }
  return Math.round(v / 100) / 10;
}
