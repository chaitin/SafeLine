import fs from "fs";
import path from "path";
import matter from "gray-matter";

const postsDir = path.join(process.cwd(), "./src/static/md/");

const fileNames = fs.readdirSync(postsDir);

export const getPostsIds = () => {
  return fileNames.map((fileName) => ({
    params: { id: fileName.replace(/\.md$/, "") },
  }));
};

interface MatterMark {
  data: { date: string; title: string };
  content: string;
  [key: string]: unknown;
}

export const getPostsData = () => {
  return fileNames.map((fileName) => {
    // 去除文件名的md后缀，使其作为文章id使用
    const id = fileName.replace(/\.md$/, "");

    // 获取md文件路径
    const fullPath = path.join(postsDir, fileName);

    // 读取md文件内容
    const fileContents = fs.readFileSync(fullPath, "utf8");

    // 使用matter提取md文件元数据：{data:{//元数据},content:'内容'}
    const matterResult = matter(fileContents);
    return {
      id,
      ...(matterResult.data as MatterMark["data"]),
    };
  });
};
