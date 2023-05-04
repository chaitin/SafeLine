import fs from "fs";
import path from "path";
import matter from "gray-matter";
import { serialize } from "next-mdx-remote/serialize";
// import prism from "remark-prism";
import slug from "remark-slug";
import remarkGfm from "remark-gfm";
import externalLinks from "remark-external-links";

const postsDir = path.join(process.cwd(), "./src/static/md/");

const fileNames = fs.readdirSync(postsDir);

export const getPostsIds = () => {
  return fileNames.map((fileName) => ({
    params: { id: fileName.replace(/\.md$/, "") },
  }));
};

type Group = string | { title: string; order?: number };
export interface NavItem {
  title: string;
  group: Group;
  order?: number;
  id: string;
}

export interface GroupItem {
  title: string;
  order?: number;
  list: NavItem[];
}

export const sortComparer = <
  T extends { order?: number; id?: string; title: string }
>(
  a: T,
  b: T
) => {
  return (
    ("order" in a && "order" in b ? a.order! - b.order! : 0) ||
    ("id" in a && "id" in b ? a.id!.localeCompare(b.id!) : 0) ||
    (a.title ? a.title.localeCompare(b.title) : -1)
  );
};

const dealWithData: (data: NavItem[]) => GroupItem[] = (data) => {
  let list: GroupItem[] = [];
  let cache: { [key: string]: NavItem } = {};
  data.forEach((ele) => {
    let groupTitle: string, groupOrder: number | undefined;

    if (typeof ele.group === "object") {
      groupTitle = ele.group.title;
      groupOrder = ele.group.order;
    } else if (typeof ele.group === "string") {
      groupTitle = ele.group;
    }

    if (!cache[groupTitle!]) {
      let group: GroupItem = {
        title: groupTitle!,
        list: [ele],
      };
      if (groupOrder) group.order = groupOrder;
      list.push(group);
      cache[groupTitle!] = ele;
    } else {
      const cur = list.find((item) => item.title === groupTitle);
      if (groupOrder && cur) {
        cur.order = groupOrder;
      }
      cur?.list.push(ele);
    }
  });

  return list.sort(sortComparer).map((group) => ({
    ...group,
    list: group.list.sort(sortComparer),
  }));
};

export const getPostsGroup = () => {
  const dataList: NavItem[] = fileNames.map((fileName) => {
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
      ...(matterResult.data as Omit<NavItem, "id">),
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
      remarkPlugins: [externalLinks, remarkGfm, slug],
      // rehypePlugins: [
      //   [
      //     toc,
      //     {
      //       headings: ["h1", "h2", "h3", "h4"],
      //       customizeTOC: (tocAll: any) => {
      //         tocElement = tocAll;
      //         return false;
      //       },
      //     },
      //   ],
      // ],
    },
  });
  // console.log(tocElement, "tocElement--");

  return { postData: res, content: matterResult.content };
}
