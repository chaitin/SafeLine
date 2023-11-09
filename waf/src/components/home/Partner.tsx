import React from 'react'
import { Box, Button, Container, Grid, Typography } from "@mui/material";
import Image from 'next/image'
import Icon from "@/components/Icon";

const PARTNER_LIST = [
  {
    title: "中国银行",
    icon: "boc.png",
  },
  {
    title: "交通银行",
    icon: "bcm.png",
  },
  {
    title: "滴滴出行",
    icon: "didi.png",
  },
  {
    title: "华为",
    icon: "huawei.png",
  },
  {
    title: "清华大学",
    icon: "thu.png",
  },
  {
    title: "bilibili",
    icon: "bilibili.png",
  },
  {
    title: "抖音",
    icon: "douyin.png",
  },
  {
    title: "爱奇艺",
    icon: "aqiyi.png",
  },
  {
    title: "顺丰快递",
    icon: "shunfeng.png",
  },
  {
    title: "中国南方航空",
    icon: "csn.png",
  },
];

const ARTICLE_LIST = [
  "入围 Gartner 《Web 应用防火墙魔力象限》",
  "入围 Forrester 《NowTech:WebApplicationFirewalls,Q42019》"
]

const Partner = () => {
  return (
    <Container>
      <Box
        sx={{
          px: 5,
          pt: 19,
        }}
        className="flex flex-col items-center"
      >
        <Typography variant="h2" mb={2}>
          优秀企业的信赖之选
        </Typography>
        {ARTICLE_LIST.map((item) => (
          <Typography
            key={item}
            variant='h5'
            sx={{ color: "#86909C", fontWeight: 500, wordBreak: "break-word" }}
            mt={2}
          >
            <Icon type="icon-jingxuan" className="mr-3 relative top-1" sx={{ fontSize: "1.3em" }} />
            {item}
          </Typography>
        ))}
        <Grid container mt={4} spacing={3} display="flex" justifyContent="center">
          {PARTNER_LIST.map((item) => (
            <Grid item key={item.title} mb={1}>
              <Box
                sx={{
                  width: { sx: "140px", md: "160px" },
                  height: { sx: "auto", md: "95px" },
                  backgroundColor: "common.white",
                  boxShadow: "0px 8px 24px -4px rgba(145,158,171,0.1)",
                  borderRadius: "4px",
                }}
                className="flex justify-center items-center"
              >
                <Box sx={{ width: "140px" }}>
                  <Image
                    src={`/images/logo/${item.icon}`}
                    alt={item.title}
                    layout="responsive"
                    width={100}
                    height={100}
                  />
                </Box>
              </Box>
            </Grid>
          ))}
        </Grid>
        <Box mt={6}>
          <Button
            variant="outlined"
            target="_blank"
            sx={{
              width: { xs: "100%", sm: "146px" },
              height: "50px",
              ml: { xs: 2, sm: 2 },
              mb: { xs: 2, sm: 0 },
            }}
            href="https://www.chaitin.cn/zh/"
          >
            了解详情
          </Button>
        </Box>
      </Box>
    </Container>
  );
};

export default Partner;
