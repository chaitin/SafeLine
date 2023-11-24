export {
    getSetupCount,
    getDiscussions,
    getIssues,
  };
  
  const BASE_API = "/api";

  function getAPIUrl() {
    return (process.env.TARGET || '') + BASE_API;
  }

  function getSetupCount() {
    return fetch("https://waf-ce.chaitin.cn/api/count").then((res) => res.json());
  }

  function getDiscussions(query: string) {
    return fetch(getAPIUrl() + "/repos/discussions?q=" + query).then((res) => res.json());
  }

  function getIssues(query: string) {
    return fetch(getAPIUrl() + "/repos/issues?q=" + query).then((res) => res.json());
  }
  