import Image from "next/image";
import Link from "next/link";
import { Box, Grid, Button, Typography, Container, Stack } from "@mui/material";
import { Carousel } from "react-responsive-carousel";
import Version from "@/components/Home/Version";
import FriendlyLinks from "@/components/Home/FriendlyLinks";
import "react-responsive-carousel/lib/styles/carousel.min.css";

const IMAGE_LIST = [
  {
    name: "攻击检测列表",
    desc: "",
    url: "/images/album/1.png",
  },
  {
    name: "攻击检测详情",
    desc: "可详细展示 攻击目标、攻击源 IP、攻击载荷、原始报文等详细信息",
    url: "/images/album/2.png",
  },
  {
    name: "防护站点列表",
    desc: "",
    url: "/images/album/3.png",
  },
  {
    name: "攻击阻断页面",
    desc: "检测到攻击时，将直接返回 403 错误，流量将被清洗，防止网站受损",
    url: "/images/album/block.png",
  },
];

export default function Home() {
  return (
    <Container maxWidth="lg">
      <Grid container sx={{ pt: 8 }}>
        <Grid item container sx={{ flex: 1 }} xs={12} sm={8}>
          <Grid item xs={12}>
            <Typography variant="h3" sx={{ pb: 3 }}>
              雷池 Web 应用防火墙{" "}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="subtitle1" sx={{ pb: 3 }}>
              耗时近 10 年，长亭科技倾情打造，核心检测能力由
              <Box component="span" color="primary.main">
                {" "}
                智能语义分析算法{" "}
              </Box>
              驱动，专为社区而生，不让黑客越雷池半步。
            </Typography>
          </Grid>
          <Grid container item xs={12} spacing={2}>
            <Button
              variant="contained"
              component={Link}
              // fullWidth
              target="_blank"
              sx={{
                width: { xs: "100%", sm: "auto" },
                ml: { xs: 2, sm: 2 },
                mb: { xs: 2, sm: 0 },
              }}
              href="https://stack.chaitin.com/tool/detail?id=717"
            >
              免费使用
            </Button>
            <Button
              variant="contained"
              color="neutral"
              component={Link}
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
                },
              }}
              href="https://github.com/chaitin/safeline"
              startIcon={
                <Image
                  src="/images/github.png"
                  alt="Logo"
                  width={16}
                  height={16}
                  priority
                />
              }
            >
              GitHub
            </Button>
            <Button
              variant="contained"
              color="neutral"
              sx={{
                backgroundColor: "#fff",
                color: "#000",
                ml: { xs: 2, sm: 2 },
                width: { xs: "100%", sm: "auto" },
              }}
              component={Link}
              href="/#groupchat"
              startIcon={
                <Image
                  src="/images/wechat-logo.png"
                  alt="Logo"
                  width={16}
                  height={16}
                  priority
                />
              }
            >
              讨论组
            </Button>
          </Grid>
        </Grid>
        <Grid
          item
          xs={12}
          sm={4}
          sx={{
            display: "flex",
            justifyContent: { xs: "center", sm: "right" },
            pt: { xs: 3, sm: 0 },
          }}
        >
          <Image
            src="/images/safeline.png"
            alt="Logo"
            width={196}
            height={196}
            style={{ marginLeft: "24px" }}
            priority
          />
        </Grid>
      </Grid>
      <Box sx={{ pt: 12 }}>
        <Carousel
          interval={3000}
          infiniteLoop
          autoPlay
          showStatus={false}
          showThumbs={false}
        >
          {IMAGE_LIST.map((item) => (
            <Box key={item.url}>
              <img src={item.url} alt={item.desc} />
              <Box className="legend">
                <Typography variant="h5">{item.name}</Typography>
                <Typography variant="h6">{item.desc}</Typography>
              </Box>
            </Box>
          ))}
        </Carousel>
      </Box>

      <Stack id="datasheet" sx={{ pt: 12 }} spacing={6} alignItems="center">
        <Typography variant="h3">版本对比</Typography>
        <Version />
      </Stack>
      <Stack sx={{ pt: 12 }} id="groupchat" spacing={6} alignItems="center">
        <Typography variant="h3">加入讨论组</Typography>
        <Image
          src="/images/wechat.png"
          alt="wechat"
          width={300}
          height={300}
          priority
        />
      </Stack>
      <Box sx={{ pt: 12 }}>
        <FriendlyLinks />
      </Box>
      <Stack
        sx={{ pt: 6, color: "text.auxiliary", textAlign: "center" }}
        justifyContent="center"
      >
        © 2023 北京长亭未来科技有限公司京 ICP 备 19035216 号京公网安备
        11010802020947 号
      </Stack>
    </Container>
  );
}
