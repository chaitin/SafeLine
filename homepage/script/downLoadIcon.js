const fs = require("fs");
const path = require("path");

const Axios = require("axios");

async function downloadFile(url) {
  const iconPath = path.resolve(__dirname, "../src/static/fonts/iconfont.js");
  const writer = fs.createWriteStream(iconPath);
  const response = await Axios({
    url: "https:" + url,
    method: "GET",
    responseType: "stream",
  });
  response.data.pipe(writer);
  writer.on("finish", () => console.log("下载成功"));
  writer.on("error", (err) => console.log(err));
}
let argument = process.argv.splice(2);
downloadFile(argument[0]);
