import { type GetStaticProps } from "next";
import { useRef } from "react";
import { Box } from "@mui/material";
import Head from "next/head";
import SideLayout from "@/layout/SideLayout";
import { MarkdownNavbar } from "@/components";
import { MDXRemote, MDXRemoteProps } from "next-mdx-remote";

import {
  getPostsIds,
  getPostsGroup,
  getPostData,
  type GroupItem,
} from "@/utils/posts";

interface PostsProps {
  group: GroupItem[];
  postData: MDXRemoteProps;
  content: string;
  headTitle: string;
}

const Posts: React.FC<PostsProps> = ({
  group,
  postData,
  content,
  headTitle,
}) => {
  const domRef = useRef<HTMLElement>();

  return (
    <SideLayout list={group}>
      <Head>
        <title>{headTitle}</title>
      </Head>
      <Box
        sx={{
          display: "flex",
          height: "calc(100vh - 168px)",
          overflow: "auto",
          justifyContent: "space-between",
          "&::-webkit-scrollbar": {
            display: "none",
          },
        }}
      >
        <Box
          ref={domRef}
          id="markdown-body"
          className="markdown-body"
          sx={{ width: { xs: "100%", sm: "77%" } }}
        >
          <MDXRemote {...postData} />
        </Box>
        <Box
          sx={{
            display: { xs: "none", sm: "block" },
            width: "20%",
            position: "sticky",
            top: 0,
          }}
        >
          <MarkdownNavbar source={content} container={domRef.current!} />
        </Box>
      </Box>
    </SideLayout>
  );
};

export async function getStaticPaths() {
  const paths = getPostsIds();
  return {
    paths,
    fallback: false,
  };
}

export const getStaticProps: GetStaticProps<PostsProps> = async ({
  params,
}) => {
  const { postData, content } = await getPostData(params?.id as string);
  const group = getPostsGroup();

  let idToTextMap: Record<string, string> = {};
  group.forEach((g) => {
    g.list.forEach((nav) => {
      idToTextMap[nav.id] = nav.title;
    });
  });

  return {
    props: {
      group,
      content,
      postData,
      headTitle: idToTextMap[params?.id as string] + " - 长亭雷池 WAF 社区版",
    },
  };
};

export default Posts;
