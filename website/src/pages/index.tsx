import React, { useEffect, useRef } from "react";
import { Carousel } from "react-responsive-carousel";
import Features from "@site/src/components/Features";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import { Box, Grid, Button, Typography, Container, Stack } from "@mui/material";
import Layout from "@theme/Layout";
import Title from "@site/src/components/Title";
import Version from "@site/src/components/Version";
import { getSetupCount } from "@site/src/api/home";
import "react-responsive-carousel/lib/styles/carousel.min.css";
import ThemeProvider from "@site/src/components/Theme";
import Head from "@docusaurus/Head";

const IMAGE_LIST = [
  {
    name: "可视化仪表盘",
    url: "/images/album/0.png",
  },
  {
    name: "登录页",
    url: "/images/album/5.png",
  },
  {
    name: "攻击检测列表",
    url: "/images/album/1.png",
  },
  {
    name: "攻击检测详情",
    url: "/images/album/2.png",
  },
  {
    name: "防护站点列表",
    url: "/images/album/3.png",
  },
  {
    name: "自定义规则列表",
    url: "/images/album/3.png",
  },
  {
    name: "攻击阻断页面",
    url: "/images/album/block.png",
  },
];

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();

  const totalRef = useRef(null);

  const initTotal = async (n: number) => {
    const countUpModule = await import("countup.js");
    const anim = new countUpModule.CountUp(totalRef.current!, Math.max(0, n), {
      duration: 2,
    });
    anim.start();
  };

  useEffect(() => {
    getSetupCount().then((d) => {
      initTotal(d.total);
    });
  });

  return (
    <Layout title="" description="长亭雷池 Web 应用防火墙 | 长亭雷池 WAF">
      <Head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta
          name="keywords"
          content="WAF,雷池,长亭,社区版,免费,开源,网站防护"
        ></meta>
      </Head>
      <ThemeProvider>
        <Box sx={{ backgroundColor: "#F8F9FC" }}>
          <Box
            sx={{ pt: 16, pb: 18, backgroundColor: "#0F1935", color: "#fff" }}
          >
            <Container maxWidth="lg">
              <Box
                sx={{ display: { xs: "block", sm: "flex" }, flexWrap: "wrap" }}
              >
                <Grid container sx={{ flex: 1 }}>
                  <Grid item xs={12}>
                    <Typography variant="h3" sx={{ pb: 3 }}>
                      雷池 Web 应用防火墙
                    </Typography>
                  </Grid>
                  <Grid item xs={12}>
                    <Typography variant="subtitle1" sx={{ pb: 3 }}>
                      耗时近 10 年，长亭科技倾情打造，核心检测能力由
                      <Box component="span" color="primary.main" sx={{ px: 1 }}>
                        智能语义分析算法
                      </Box>
                      驱动，专为社区而生，不让黑客越雷池半步。
                    </Typography>
                  </Grid>
                  <Grid container item xs={12} spacing={2}>
                    <Button
                      variant="contained"
                      target="_blank"
                      sx={{
                        width: { xs: "100%", sm: "auto" },
                        ml: { xs: 2, sm: 2 },
                        mb: { xs: 2, sm: 0 },
                        "&:hover": {
                          backgroundColor: "rgb(10, 138, 135)",
                          color: "white",
                        },
                      }}
                      href="https://stack.chaitin.com/tool/detail?id=717"
                    >
                      免费使用
                    </Button>
                    <Button
                      variant="contained"
                      target="_blank"
                      sx={{
                        textTransform: "none",
                        backgroundColor: "#fff",
                        color: "#000",
                        ml: { xs: 2, sm: 2 },
                        mb: { xs: 2, sm: 0 },
                        width: { xs: "100%", sm: "auto" },
                        "&:hover": {
                          fontWeight: "500",
                          backgroundColor: "white",
                        },
                      }}
                      href="https://github.com/chaitin/safeline"
                      startIcon={
                        <img src="/images/github.png" width={16} height={16} />
                      }
                    >
                      GitHub
                    </Button>
                    <Button
                      variant="contained"
                      sx={{
                        backgroundColor: "#fff",
                        color: "#000",
                        ml: { xs: 2, sm: 2 },
                        width: { xs: "100%", sm: "auto" },
                        "&:hover": {
                          backgroundColor: "white",
                        },
                      }}
                      href="/#groupchat"
                      startIcon={
                        <img
                          src="/images/wechat-logo.png"
                          alt="Logo"
                          width={16}
                          height={16}
                        />
                      }
                    >
                      讨论组
                    </Button>
                  </Grid>
                </Grid>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: { xs: "center", sm: "right" },
                    pt: { xs: 3, sm: 0 },
                    ml: { xs: 0, sm: 3 },
                  }}
                >
                  <img
                    src="/images/403.svg"
                    alt="Logo"
                    width={196}
                    height={196}
                  />
                </Box>
              </Box>
            </Container>
          </Box>
          <Container sx={{ mt: -10, color: "#000", pb: 3 }}>
            <Features />
          </Container>
          <Container>
            <Stack sx={{ pt: 15 }} spacing={3} alignItems="center">
              <Title title="装机量" />
              <Typography
                sx={{
                  color: "primary.main",
                  fontSize: "96px",
                  letterSpacing: "10px",
                }}
                ref={totalRef}
              >
                -
              </Typography>
            </Stack>
          </Container>
          <Container
            sx={{
              color: "#000",
              ".carousel .control-dots .dot": {
                backgroundColor: "#000",
              },
              ".carousel .control-prev.control-arrow": {
                padding: "20px",
                borderRadius: "12px 0 0 12px",
              },
              ".carousel .control-next.control-arrow": {
                padding: "20px",
                borderRadius: "0 12px 12px 0",
              },
              ".carousel .control-prev.control-arrow:before": {
                borderRightColor: "rgba(0,0,0,0.5)",
              },
              ".carousel .control-next.control-arrow:before": {
                borderLeftColor: "rgba(0,0,0,0.5)",
              },
              ".carousel .slide .legend": {
                width: "30%",
                marginLeft: "-15%",
              },
            }}
          >
            <Stack sx={{ pt: 15 }} spacing={6} alignItems="center">
              <Title title="产品展示" />
              <Box sx={{ boxShadow: "0px 0px 6px #0fc6c2" }}>
                <Carousel
                  interval={2000}
                  infiniteLoop
                  autoPlay
                  showStatus={false}
                  showThumbs={false}
                >
                  {IMAGE_LIST.map((item) => (
                    <Box
                      key={item.url}
                      sx={{
                        borderRadius: "12px",
                        overflow: "hidden",
                      }}
                    >
                      <Box component="img" src={item.url} alt={item.name} />
                      <Box
                        className="legend"
                        sx={{
                          opacity: "0.40 !important",
                          py: "4px !important",
                          borderRadius: "4px !important",
                        }}
                      >
                        <Typography variant="h6" sx={{ fontSize: "14px" }}>
                          {item.name}
                        </Typography>
                      </Box>
                    </Box>
                  ))}
                </Carousel>
              </Box>
            </Stack>
          </Container>
          <Container sx={{ color: "#000", pb: 3 }}>
            <Stack sx={{ pt: 15 }} spacing={3} alignItems="center">
              <Title title="版本对比" />
              <Version />
            </Stack>

            <Stack
              sx={{ pt: 15 }}
              id="groupchat"
              spacing={6}
              alignItems="center"
            >
              <Title title="加入讨论组" />
              <img
                src="/images/wechat-light.png"
                alt="wechat"
                width={300}
                height={300}
              />
            </Stack>
          </Container>
        </Box>
      </ThemeProvider>
    </Layout>
  );
}
