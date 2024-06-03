import React from "react";
import Image from "next/image";
import {
  Box,
  Grid,
  Typography,
  Stack,
  SxProps,
  Container,
  Link,
} from "@mui/material";
import Icon from "@/components/Icon";

const LINKS = [
  {
    title: "Resource",
    items: [
      {
        label: "Technical documentation",
        to: "/docs",
      },
      {
        label: "Teaching videos",
        to: "https://www.bilibili.com/medialist/detail/ml2342694989",
      },
      {
        label: "Learning materials",
        to: "/docs",
      },
      {
        label: "Update logs",
        to: "/docs/about/changelog",
      },
    ],
  },
  {
    title: "About Us",
    items: [
      {
        label: "Chaitin",
        to: "https://www.chaitin.cn/zh/",
      },
      {
        label: "CT Stack Safe Community",
        to: "https://stack.chaitin.cn/",
      },
    ],
  },
];

export const items = [
  { to: "/community", label: "Developer", target: "_self" },
  {
    to: "https://github.com/chaitin/SafeLine/blob/main/LICENSE.md",
    label: "User Agreement",
    target: "_blank",
  },
];

export default function Footer() {
  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: "#121426",
      }}
    >
      <Container maxWidth="lg">
        <Grid
          container
          justifyContent="space-around"
          columns={24}
          sx={{ pb: 5, pt: 6 }}
          mt={0}
        >
          <Grid item xs={24} md={10}>
            <Stack id="groupchat" spacing={4} alignItems="flex-start">
              <Link href="/">
                <Grid
                  container
                  flexDirection="row"
                  display="flex"
                  alignItems="center"
                  sx={{ marginTop: 0 }}
                >
                  <Image
                    src="/images/safeline.svg"
                    alt="SafeLine Logo"
                    width={40}
                    height={43}
                  />
                  <Typography
                    variant="h4"
                    sx={{
                      color: "common.white",
                      fontFamily: "GilroyBold",
                      marginLeft: "16px",
                      fontSize: { xs: "40px", md: "28px" },
                      position: "relative",
                      bottom: 5,
                    }}
                  >
                    SafeLine WAF
                  </Typography>
                </Grid>
              </Link>
              <Box>
                {items.map((item, index) => (
                  <Box key={index} component="span" mr={5}>
                    {item.to ? (
                      <Link
                        sx={{
                          fontSize: "16px",
                          fontWeight: 600,
                          color: "common.white",
                        }}
                        href={item.to}
                        rel={item.label}
                        target={item.target}
                      >
                        {item.label}
                      </Link>
                    ) : (
                      <Typography
                        component="span"
                        sx={{
                          fontSize: "16px",
                          fontWeight: 600,
                          color: "common.white",
                        }}
                      >
                        {item.label}
                      </Typography>
                    )}
                  </Box>
                ))}
              </Box>
            </Stack>
          </Grid>
          {LINKS.map((link) => (
            <Grid item xs={8} md={5} my={{ xs: 4, md: 0 }} key={link.title}>
              <Stack id="groupchat" spacing={1} alignItems="flex-start">
                <Title title={link.title} />
                <Grid container>
                  {link.items.map((item, index) => (
                    <Grid key={index} item xs={12}>
                      <Link
                        sx={{
                          fontSize: "14px",
                          color: "common.white",
                          opacity: 0.5,
                          fontWeight: 400,
                          lineHeight: "38px",
                        }}
                        href={item.to}
                        target="_blank"
                        rel={item.label}
                      >
                        {item.label}
                      </Link>
                    </Grid>
                  ))}
                </Grid>
              </Stack>
            </Grid>
          ))}
          <Grid item xs={8} md={2} my={{ xs: 4, md: 0 }}>
            <Stack
              id="groupchat"
              spacing={1}
              alignItems="center"
              direction="row"
            >
              <Link
                href="https://discord.gg/wyshSVuvxC"
                target="_blank"
                sx={{
                  color: "common.white",
                  display: "flex",
                  "&:hover": {
                    color: "primary.main",
                  },
                }}
                mr={3.5}
              >
                <Box component='img' src='/images/logo/discord.svg' sx={{mr: 1}}></Box>
                Discord
              </Link>
            </Stack>
          </Grid>
        </Grid>
        <Grid
          container
          sx={{ pb: 1.5 }}
          justifyContent={{ xs: "center", md: "flex-start" }}
        >
          <Typography
            variant="inherit"
            sx={{
              fontSize: "12px",
              color: "rgba(255,255,255,0.26)",
              fontWeight: 400,
            }}
          >
            Copyright © 2024 Beijing Chaitin Future Technology.All rights
            reserved.
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
