import React from "react";
import { Box, Grid, Button, Typography, Container, Stack } from "@mui/material";
import Image from 'next/image';
import Head from 'next/head';
import DiscussionList, { Discussion } from '@/components/community/DiscussionList';
import IssueList, { Issue } from '@/components/community/IssueList';
import { getDiscussions, getIssues } from "@/api";

type CommunityPropsType = {
  discussions: Discussion[];
  issues: Issue[];
};

export async function getServerSideProps() {
  let discussions: Discussion[] = []
  let issues: Issue[] = []

  const promises = [
    getDiscussions('').then((result) => discussions = result || []),
    getIssues('').then((result) => issues = result || []),
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
            height: "380px",
            position: 'relative',
          }}
        >
          <Image
            src="/images/community-banner.png"
            alt="开发计划"
            layout="fill"
            objectFit="cover"
            objectPosition="center"
            quality={100}
          />
          <Container className="relative">
            <Box pt={23}>
              <Typography variant="h2" sx={{ fontFamily: "AlimamaShuHeiTi-Bold" }}>勇往直前，开创无限可能的未来</Typography>
              <Typography variant="subtitle1" sx={{ opacity: 0.5 }} mt={1}>同样欢迎你的参与</Typography>
            </Box>
          </Container>
        </Box>
        <Container sx={{ pt: 6, mb: 18 }}>
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
              mt: 19,
              mb: 26,
              background: "#111227",
              borderRadius: "16px",
            }}
            className="flex flex-col justify-center"
          >
            <Grid container>
              <Grid item xs={12} md={6}>
                <Stack sx={{ color: "common.white" }}>
                  <Typography variant="h3" sx={{ fontSize: "36px" }}>绕过反馈</Typography>
                  <Typography mt={1} variant="body1" sx={{ opacity: 0.5 }}>向长亭提交雷池 XSS、SQL 绕过，您将获得现金和实物奖励</Typography>
                  <Button
                    variant="outlined"
                    target="_blank"
                    sx={{
                      width: { xs: "146px" },
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
              <Grid item xs={12} md={6} display="flex" justifyContent={{ sx: 'flex-start', md: 'flex-end' }} mt={{ xs: 2, md: 0 }}>
                <Box
                  sx={{
                    width: { xs: '100%', md: '367px' },
                    height: { xs: 'auto', md: '206px' },
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
      </Box>
    </main>
  );
}

export default Community;
