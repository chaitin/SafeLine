export {
    getSetupCount,
    getDiscussions,
    getIssues,
  };
  
  const BASE_API = "/api";

  const BASE_URL = "http://10.10.4.142:8080/api";

  function getSetupCount() {
    return fetch(BASE_API + "/count").then((res) => res.json());
  }

  function getDiscussions(query: string, isClient  = true) {
    return fetch((isClient ? BASE_API : BASE_URL) + "/repos/discussions?q=" + query).then((res) => res.json());
  }

  function getIssues(query: string, isClient  = true) {
    return fetch((isClient ? BASE_API : BASE_URL) + "/repos/issues?q=" + query).then((res) => res.json());
  }
  