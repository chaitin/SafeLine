import React from "react";
import Head from 'next/head';
import Image from 'next/image';
import { Box, Container, Stack, Typography } from "@mui/material";
import Version from "@/components/Version";

export default function VersionView() {

  return (
    <main title="版本对比 ｜ 雷池 WAF 社区版">
      <Head>
        <title>版本对比 | 雷池 WAF 社区版</title>
        <meta name="keywords" content="WAF,雷池,社区版,免费,版本对比,企业版,智能语义分析检测"></meta>
        <meta name="description" content="雷池 WAF 社区版，大小网站皆宜，免费开启使用"></meta>
      </Head>
      <Box mb={{ xs: 10, md: 26 }}>
        <Box
          sx={{
            width: "100%",
            height: { xs: "866px", md: "380px" },
            position: "relative",
            backgroundImage: { xs: "url(/images/version-banner-mobile.png)", md: "url(/images/version-banner.png)" },
            backgroundSize: "cover",
            backgroundPosition: 'center center',
            backgroundRepeat: 'no-repeat'
          }}
        >
          <Container className="relative">
            <Box pt={{ xs: 21, md: 23 }}>
              <Typography
                variant="h2"
                sx={{
                  fontFamily: "AlimamaShuHeiTi-Bold",
                  textAlign: { xs: "center", md: "left" },
                  fontSize: { xs: "32px", md: "48px" },
                }}
              >大小网站皆宜，免费即可开始</Typography>
            </Box>
          </Container>
        </Box>
        <Container sx={{ position: "relative", bottom: { xs: 124, md: 0 }, mb: { xs: "-124px", md: 0 } }}>
          <Stack pt={{ xs: 0, md: 20 }} spacing={3} alignItems="center">
            <Version />
          </Stack>
        </Container>
      </Box>
    </main>
  );
}
