export { getSetupCount };

function getSetupCount() {
  return fetch("/api/count").then((res) => res.json());
}
