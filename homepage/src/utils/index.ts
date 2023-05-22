export * from "./render";

export { sampleLength, sizeLength, sampleSummary };

function sampleLength(s: string) {
  const l = new Blob([s]).size;
  return l > 1024 * 2 ? Math.round(l / 1024) + "KB" : l + "B";
}

function sizeLength(l: number) {
  return l > 1024 * 2 ? Math.round(l / 1024) + "KB" : l + "B";
}

function sampleSummary(s: string) {
  return s.split("\n").slice(0, 2).join(" ").slice(0, 60);
}
