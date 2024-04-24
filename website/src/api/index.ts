export {
  getSetupCount,
  getDiscussions,
  getIssues,
  getReposInfo,
  detectorPoint,
};

const BASE_API = "/api";

function getAPIUrl() {
  return (process.env.TARGET || "") + BASE_API;
}

function getSetupCount() {
  return fetch(getAPIUrl() + "/safeline/count").then((res) => res.json());
}

function getReposInfo() {
  return fetch(getAPIUrl() + "/repos/info").then((res) => res.json());
}

function getDiscussions(query: string) {
  return fetch(getAPIUrl() + "/repos/discussions?q=" + query).then((res) =>
    res.json()
  );
}

function getIssues(query: string) {
  return fetch(getAPIUrl() + "/repos/issues?q=" + query).then((res) =>
    res.json()
  );
}

// 购买 1001
// 咨询 1002
function detectorPoint(query: { source?: string; type: 1001 | 1002 }) {
  const search = new URLSearchParams(location.search);
  query.source = search.get("source") || document.referrer;
  return fetch(
    getAPIUrl() + "/api/behavior?" + new URLSearchParams(query as any), {method: 'POST'}
  ).then((res) => res.json());
}
