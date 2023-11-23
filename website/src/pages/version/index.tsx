import React from "react";
import Head from 'next/head';
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
      <Box mb={26}>
        <Box
          sx={{
            width: "100%",
            height: "380px",
            backgroundImage: "url(/images/version-banner.png)",
            backgroundSize: "cover",
            position: 'relative',
            backgroundPosition: 'center center',
            backgroundRepeat: 'no-repeat'
          }}
        >
          <Container>
            <Box pt={23}>
              <Typography variant="h2" sx={{ fontFamily: "AlimamaShuHeiTi-Bold" }}>大小网站皆宜，免费即可开始</Typography>
            </Box>
          </Container>
        </Box>
        <Container>
          <Stack sx={{ pt: 20 }} spacing={3} alignItems="center">
            <Version />
          </Stack>
        </Container>
      </Box>
    </main>
  );
}
