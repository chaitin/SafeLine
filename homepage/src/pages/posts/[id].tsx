import { type GetStaticProps } from "next";
import { Box } from "@mui/material";
import SideLayout from "@/layout/SideLayout";

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
}

const Posts: React.FC<PostsProps> = ({ group, postData }) => {
  return (
    <SideLayout list={group}>
      <Box id="markdown-body" className="markdown-body">
        <MDXRemote {...postData} />
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
  const postData = await getPostData(params?.id as string);

  const group = getPostsGroup();
  return {
    props: {
      group,
      postData,
    },
  };
};

export default Posts;
