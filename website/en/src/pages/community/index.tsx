import React from "react";
import { Box, Grid, Button, Typography, Container, Stack } from "@mui/material";
import Image from 'next/image';
import Head from 'next/head';
import DiscussionList, { Discussion } from '@/components/community/DiscussionList';
import IssueList, { Issues } from '@/components/community/IssueList';
import { getDiscussions, getIssues } from "@/api";

type CommunityPropsType = {
  discussions: Discussion[];
  issues: Issues;
};

export async function getServerSideProps() {
  let discussions: Discussion[] = []
  let issues: Issues = {}

  const promises = [
    getDiscussions('').then((result) => discussions = result || []),
    getIssues('').then((result) => issues = result || {}),
  ];
  try {
    await Promise.allSettled(promises)
  } finally {
    return {
      props: {
        discussions: discussions,
        issues: issues,
      },
    }
  }
}

function Community({ discussions, issues }: CommunityPropsType) {
  return (
    <main title="社区 ｜ 雷池 WAF 社区版">
      <Head>
        <title>社区 | 雷池 WAF 社区版</title>
        <meta name="keywords" content="WAF,雷池,社区版,免费,社区,反馈,discussion,roadmap,"></meta>
        <meta name="description" content="雷池 WAF 社区版，欢迎你通过社区获取更多帮助"></meta>
      </Head>
      <Box>
        <Box
          sx={{
            width: "100%",
            height: { xs: "866px", md: "380px" },
            position: 'relative',
            backgroundImage: { xs: "url(/images/community-banner-mobile.png)", md: "url(/images/community-banner.png)" },
            backgroundSize: "cover",
            backgroundPosition: 'center center',
            backgroundRepeat: 'no-repeat',
          }}
        >
          <Container className="relative">
            <Box
              pt={{ xs: 19, md: 23 }}
              textAlign={{ xs: "center", md: "left" }}
            >
              <Typography
                variant="h2"
                sx={{
                  fontFamily: "AlimamaShuHeiTi-Bold",
                  fontSize: { xs: "32px", md: "48px" }
                }}>
                  勇往直前，开创无限可能的未来
                </Typography>
              <Typography variant="subtitle1" sx={{ opacity: 0.5 }} mt={1}>同样欢迎你的参与</Typography>
            </Box>
          </Container>
        </Box>
        <Container sx={{ position: "relative", bottom: { xs: 124, md: 0 }, mb: { xs: "-124px", md: 0 } }}>
          <Container sx={{ pt: 6, mb: { xs: 10, sm: 18 }}}>
            <IssueList value={issues} />
          </Container>
          <Container>
            <DiscussionList value={discussions} />
          </Container>
          <Container>
            <Box
              px={6}
              py={6}
              sx={{
                mt: { xs: 10, sm: 19 },
                mb: { xs: 10, sm: 26 },
                background: "#111227",
                borderRadius: "16px",
              }}
              className="flex flex-col justify-center"
            >
              <Grid container>
                <Grid item xs={12} sm={6}>
                  <Stack sx={{ color: "common.white" }}>
                    <Typography variant="h3" sx={{ fontSize: "36px" }}>绕过反馈</Typography>
                    <Typography mt={1} variant="body1" sx={{ opacity: 0.5 }}>向长亭提交雷池 XSS、SQL 绕过，您将获得现金和实物奖励</Typography>
                    <Button
                      variant="outlined"
                      target="_blank"
                      sx={{
                        width: { xs: "100%", sm: "146px" },
                        height: "50px",
                        mt: { xs: 4, sm: 10 },
                        fontSize: "16px",
                        backgroundColor: "common.white",
                      }}
                      href="https://stack.chaitin.com/security-challenge/safeline/index"
                    >
                      去提交
                    </Button>
                  </Stack>
                </Grid>
                <Grid item xs={12} sm={6} display="flex" justifyContent={{ xs: 'center', md: 'flex-end' }} mt={{ xs: 2, md: 0 }}>
                  <Box
                    sx={{
                      width: '367px',
                      height: '169px',
                      mt: { xs: 2, sm: 0 },
                    }}
                  >
                    <Image
                      src="/images/feedback.png"
                      alt="XSS 挑战入口,SQL 挑战入口"
                      layout="responsive"
                      width={100}
                      height={100}
                    />
                  </Box>
                </Grid>
              </Grid>
            </Box>
          </Container>
        </Container>
      </Box>
    </main>
  );
}

export default Community;
