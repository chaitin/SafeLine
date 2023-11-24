import React from 'react';
import Image from 'next/image'
import { Box, Grid, Typography, Stack, SxProps, Container, Link } from '@mui/material';
import WafTitle from './home/WafTitle'

const LINKS = [
  {
    title: "资源",
    items: [
      {
        label: "技术文档",
        to: "/docs",
      },
      {
        label: "教学视频",
        to: "https://www.bilibili.com/medialist/detail/ml2342694989",
      },
      {
        label: "学习资料",
        to: "/docs",
      },
      {
        label: "更新日志",
        to: "/docs/about/changelog",
      },
    ],
    xs: 6,
  },
  {
    title: "关于我们",
    items: [
      {
        label: "长亭科技",
        to: "https://www.chaitin.cn/zh/",
      },
      {
        label: "CT Stack 安全社区",
        to: "https://stack.chaitin.cn/",
      },
    ],
    xs: 12,
  },
];

export const items = [
  { to: "/community", label: "社区" },
  { to: "/version", label: "版本对比" },
  { to: "", label: "用户协议" },
];

export default function Footer() {
  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: '#121426',
      }}
    >
      <Container maxWidth="lg">
        <Grid container justifyContent="space-around" columns={24} sx={{ pb: 5, pt: 6 }} mt={0}>
          <Grid item xs={24} md={10}>
            <Stack
              id="groupchat"
              spacing={4}
              alignItems="flex-start"
            >
              <WafTitle title="雷池 SafeLine" sx={{ marginLeft: '16px' }} />
              <Box>
                {items.map((item, index) => (
                  <Box key={index} component="span" mr={5}>
                    {item.to ? (
                      <Link sx={{ fontSize: '16px', fontWeight: 600, color: "common.white" }} href={item.to} rel={item.label}>
                        {item.label}
                      </Link>
                    ) : (
                      <Typography component="span" sx={{ fontSize: '16px', fontWeight: 600, color: "common.white" }}>
                        {item.label}
                      </Typography>
                    )}
                  </Box>
                ))}
              </Box>
            </Stack>
          </Grid>
          {LINKS.map((link) => (
            <Grid item xs={12} md={5} my={{ xs: 4, md: 0 }} key={link.title}>
              <Stack
                id="groupchat"
                spacing={1}
                alignItems="flex-start"
              >
                <Title title={link.title} />
                <Grid container>
                  {link.items.map((item, index) => (
                    <Grid key={index} item xs={link.xs} md={12}>
                      <Link sx={{ fontSize: '14px', color: "common.white", opacity: 0.5, fontWeight: 400, lineHeight: "38px" }} href={item.to} target="_blank" rel={item.label}>
                        {item.label}
                      </Link>
                    </Grid>
                  ))}
                </Grid>
              </Stack>
            </Grid>
          ))}
          <Grid item xs={24} md={4} display={'flex'} justifyContent={{ xs: "center", lg: "flex-end" }}>
            <Stack
              id="groupchat"
              spacing={2}
              alignItems="flex-start"
            >
              <Title title="微信交流群" />
              <Image
                src="/images/wechat-230825.png"
                alt="wechat"
                width={96}
                height={96}
              />
            </Stack>
          </Grid>
        </Grid>
        <Grid container sx={{ pb: 1.5 }}>
          <Typography
            variant="inherit"
            sx={{ fontSize: "12px", color: "rgba(255,255,255,0.26)", fontWeight: 400 }}
          >
            Copyright © 2023 北京长亭科技有限公司. All rights reserved.
          </Typography>
        </Grid>
      </Container>
    </Box>
  );
}

interface TitleProps {
  title: string;
  sx?: SxProps;
}

const Title: React.FC<TitleProps> = ({ title, sx }) => {
  return (
    <Typography
      variant="h6"
      sx={{ color: "common.white", marginLeft: "8px", ...sx }}
    >
      {title}
    </Typography>
  );
};
