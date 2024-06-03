import React from 'react'
import { Grid, Box, Typography, List, ListItem, ListItemText, Button, alpha } from "@mui/material";
import Image from 'next/image'

const FEATURE_LIST = [
  {
    title: "语义分析算法，有效保卫网站安全",
    desc: "首创语义分析算法，突破传统规则算法的极限，精准检测、低误报、难绕过",
    icon: "/images/feature1-icon.png",
    content: (
      <Grid container display={'flex'} justifyContent={'flex-end'} position={'relative'} sx={{ mb: 5 }}>
        <Grid item xs={0} md={2} sx={{ display: { xs: 'none', md: 'block' }, position: 'absolute', left: 0 }}>
          <Box>
            <Image
              src="/images/feature1-left.png"
              alt=""
              width={182}
              height={789}
            />
          </Box>
        </Grid>
        <Grid item xs={12} md={4} ml={2} mt={{ xs: 3, md: 14 }}>
          <List>
            <ListItem sx={{ mb: 1 }}>
              <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 国内首创、业内领先的智能语义分析算法" />
            </ListItem>
            <ListItem sx={{ mb: 1 }}>
              <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 基于代码的理解，防御 0day 攻击" />
            </ListItem>
            <ListItem sx={{ mb: 1 }}>
              <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 多维度 Web 应用防护" />
            </ListItem>
          </List>
          <Box>
            <Button
              variant="outlined"
              target="_blank"
              sx={{
                width: { xs: "100%", sm: "146px" },
                height: { xs: "72px", sm: "50px" },
                ml: { xs: 0, sm: 2 },
                mb: { xs: 2, sm: 0 },
                fontSize: { xs: "24px", sm: "14px" },
              }}
              href="/docs/about/syntaxanalysis"
            >
              了解详情
            </Button>
          </Box>
        </Grid>
        <Grid item xs={12} md={6} mt={7}>
          <Box
            position={'relative'}
          >
            <Image
              src="/images/feature1-bg.png"
              alt=""
              layout="responsive"
              width={100}
              height={100}
              className='absolute left-1/2 top-0'
            />
            <Image
              src="/images/feature1.svg"
              alt="SQL注入,CSRF,XSS,SSRF,..."
              layout="responsive"
              width={100}
              height={100}
              className='relative'
            />
          </Box>
        </Grid>
      </Grid>
    ),
  },
  {
    title: "化繁为简，智能安全",
    desc: "轻松上手，实现躺平式管理",
    icon: "/images/feature2-icon.png",
    content: (
      <Grid container position="relative" sx={{ mb: 5 }} ml={{ xs: 3, md: 0 }}>
        <Grid container>
          <Grid item xs={12} md={6} mt={7} order={{ xs: 2, md: 1 }}>
            <Box
              position="relative"
            >
              <Image
                src="/images/feature2-bg.png"
                alt=""
                layout="responsive"
                width={100}
                height={100}
                className='absolute right-1/3'
                style={{ bottom: "10%" }}
              />
              <Image
                src="/images/feature2.svg"
                alt="开箱即用，轻松上手，适配多种运行环境" 
                layout="responsive"
                width={100}
                height={100}
                className='relative'
                style={{ right: "43px" }}
              />
            </Box>
          </Grid>
          <Grid item xs={12} md={4} mt={{ xs: 3, md: 14 }} order={{ xs: 1, md: 2 }}>
            <List>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 一键安装，容器式管理，适配多种运行环境" />
              </ListItem>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 配置开箱即用，无需大量调整繁琐规则" />
              </ListItem>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 简洁操作，专为社区设计" />
              </ListItem>
            </List>
          </Grid>
        </Grid>
        <Grid item xs={0} md={2} sx={{ display: { xs: 'none', md: 'block' }, position: 'absolute', right: 0 }}>
          <Box>
            <Image
              src="/images/feature2-right.png"
              alt=""
              width={182}
              height={763}
            />
          </Box>
        </Grid>
      </Grid>
    ),
  },
  {
    title: "高性能、高并发、高可用性",
    desc: "无规则引擎，线性安全检测算法",
    icon: "/images/feature3-icon.png",
    content: (
      <Grid container display="flex" justifyContent="flex-end" position="relative" sx={{ mb: 10 }} ml={{ xs: 3, md: 0 }}>
        <Grid item xs={0} md={2} sx={{ display: { xs: 'none', md: 'block' }, position: 'absolute', left: 0 }}>
          <Box>
            <Image
              src="/images/feature3-left.png"
              alt=""
              width={182}
              height={649}
            />
          </Box>
        </Grid>
        <Grid item xs={12} md={4} mt={{ xs: 3, md: 14 }}>
          <List>
            <ListItem>
              <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 平均检测延迟 < 1 毫秒，单核轻松检测 2000+ TPS 并发" />
            </ListItem>
            <ListItem>
              <ListItemText sx={{ textIndent: '-0.75rem' }} primary="· 基于 Nginx 开发，完善的健康检查机制" />
            </ListItem>
          </List>
        </Grid>
        <Grid item xs={12} md={6} mt={7}>
          <Box position={'relative'}>
            <Image
              src="/images/feature3-bg.png"
              alt=""
              width={1225}
              height={1310}
              className='absolute -left-1/2 -top-2/3'
            />
            <Image
              src="/images/feature3.svg"
              alt="性能,服务可用性99.99%"
              layout="responsive"
              width={100}
              height={100}
              className='relative'
            />
          </Box>
        </Grid>
      </Grid>
    ),
  },
];

const Features = () => {
  return (
    <Box
      sx={{
        color: "common.black",
        pb: 5,
      }}
    >
      <Grid container>
        {FEATURE_LIST.map((feature, index) => (
          <Grid item xs={12} key={feature.title}>
            <Box
              display={{ xs: "block",  md: index % 2 == 1 ? 'flex' : ''}}
              alignItems="flex-end"
              flexDirection={index % 2 == 1 ? 'column' : 'row'}
            >
              <Box>
                <Image
                  src={feature.icon}
                  alt={feature.title}
                  width={72}
                  height={72}
                />
              </Box>
              <Box ml={3}>
                <Typography variant="h4">
                  <Box component="span" sx={{ mt: "5px" }}>
                    {feature.title}
                  </Box>
                </Typography>
                <Typography
                  variant="subtitle1"
                  sx={{ color: alpha("#000", 0.7), fontWeight: 100 }}
                  mt={2}
                >
                  <Box component="span">
                    {feature.desc}
                  </Box>
                </Typography>
              </Box>
            </Box>
            <Box
              sx={{ whiteSpace: "pre-line" }}
            >
              {feature.content}
            </Box>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Features;
