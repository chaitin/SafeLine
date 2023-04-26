import fs from "fs";
import path from "path";
import matter from "gray-matter";
import { serialize } from "next-mdx-remote/serialize";
// import prism from "remark-prism";
import remarkGfm from "remark-gfm";
import externalLinks from "remark-external-links";

const postsDir = path.join(process.cwd(), "./src/static/md/");

const fileNames = fs.readdirSync(postsDir);

export const getPostsIds = () => {
  return fileNames.map((fileName) => ({
    params: { id: fileName.replace(/\.md$/, "") },
  }));
};

export interface MenuItem {
  title: string;
  category: string;
  weight: number;
  id: string;
}

export interface GroupItem {
  category: string;
  list: MenuItem[];
}

const dealWithData: (data: MenuItem[]) => GroupItem[] = (data) => {
  let list: { category: string; list: MenuItem[] }[] = [];
  let cache: { [key: string]: unknown } = {};
  data.forEach((ele) => {
    if (!cache[ele.category]) {
      list.push({
        category: ele.category,
        list: [ele],
      });
      cache[ele.category] = ele;
    } else {
      const cur = list.find((item) => item.category === ele.category);
      cur?.list.push(ele);
    }
  });
  let categorys = ["上手指南", "常见问题排查", "关于雷池"];
  return list.map((item) => {
    return {
      category: item.category,
      list: item.list.sort((a, b) => a.weight - b.weight),
    };
  }).sort((a,b)=> categorys.indexOf(a.category) - categorys.indexOf(b.category));
};

export const getPostsGroup = () => {
  const dataList: MenuItem[] = fileNames.map((fileName) => {
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
      ...(matterResult.data as Omit<MenuItem, "id">),
    };
  });
  return dealWithData(dataList);
};

export async function getPostData(id: string) {
  // 文章路径
  const fullPath = path.join(postsDir, `${id}.md`);
  // const toc = require("@jsdevtools/rehype-toc");
  // 读取文章内容
  const fileContents = fs.readFileSync(fullPath, "utf8");

  // 使用matter解析markdown元数据和内容
  const matterResult = await matter(fileContents);
  let tocElement;
  const res = await serialize(matterResult.content, {
    mdxOptions: {
      remarkPlugins: [externalLinks, remarkGfm],
      // rehypePlugins: [
      //   [
      //     toc,
      //     {
      //       headings: ["h1", "h2", "h3", "h4"],
      //       customizeTOC: (tocAll) => {
      //         tocElement = tocAll;
      //         return false;
      //       },
      //     },
      //   ],
      // ],
    },
  });

  return res;
}
