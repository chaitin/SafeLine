import { Grid, Box } from "@mui/material";

import Link from "next/link";

const LINK_LIST = [
  {
    name: "北京长亭科技有限公司",
    url: "https://chaitin.cn/",
  },
  {
    name: "CTStack 安全社区",
    url: "https://stack.chaitin.cn/",
  },
  {
    name: "长亭百川云平台",
    url: "https://rivers.chaitin.cn/",
  },
  {
    name: "长亭 GitHub 主页",
    url: "https://github.com/chaitin",
  },
  {
    name: "长亭 B 站主页",
    url: "https://space.bilibili.com/521870525",
  },
  {
    name: "长亭合作伙伴论坛",
    url: "https://bbs.chaitin.cn/",
  },
];

const FriendlyLinks = () => {
  return (
    <Grid container>
      {LINK_LIST.map((item) => (
        <Grid item xs={6} sm={4} md={3} key={item.name}>
          <Box
            component={Link}
            target="_blank"
            href={item.url}
            sx={{
              color: "rgba(0,0,0,0.7)",
              fontSize: "14px",
              lineHeight: "24px",
              "&:hover": {
                color: "primary.main",
              },
            }}
          >
            {item.name}
          </Box>
        </Grid>
      ))}
    </Grid>
  );
};

export default FriendlyLinks;
