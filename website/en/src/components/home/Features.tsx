import React from "react";
import {
  Grid,
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  Button,
  alpha,
} from "@mui/material";
import Image from "next/image";

const FEATURE_LIST = [
  {
    title:
      "Semantic analysis algorithms effectively safeguard website security",
    desc: "Pioneering semantic analysis algorithm, breaking through the limits of traditional rule algorithms, precise detection, low false positives, and difficult to bypass",
    icon: "/images/feature1-icon.png",
    content: (
      <Grid
        container
        display={"flex"}
        justifyContent={"flex-end"}
        position={"relative"}
        sx={{ mb: 5 }}
      >
        <Grid
          item
          xs={0}
          md={2}
          sx={{
            display: { xs: "none", md: "block" },
            position: "absolute",
            left: 0,
          }}
        >
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
              <ListItemText
                sx={{ textIndent: "-0.75rem" }}
                primary="· &nbsp;The first and industry-leading intelligent semantic analysis algorithm in China"
              />
            </ListItem>
            <ListItem sx={{ mb: 1 }}>
              <ListItemText
                sx={{ textIndent: "-0.75rem" }}
                primary="· &nbsp;Based on code understanding, defend against 0day attacks"
              />
            </ListItem>
            <ListItem sx={{ mb: 1 }}>
              <ListItemText
                sx={{ textIndent: "-0.75rem" }}
                primary="· &nbsp;Multidimensional Web Application Protection"
              />
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
                fontFamily: 'GilroyBold',
                fontSize: { xs: "24px", sm: "14px" },
              }}
              href="/docs/about/syntaxanalysis"
            >
              Learn more
            </Button>
          </Box>
        </Grid>
        <Grid item xs={12} md={6} mt={7}>
          <Box position={"relative"}>
            <Image
              src="/images/feature1-bg.png"
              alt=""
              layout="responsive"
              width={100}
              height={100}
              className="absolute left-1/2 top-0"
            />
            <Image
              src="/images/feature1.svg"
              alt="SQL注入,CSRF,XSS,SSRF,..."
              layout="responsive"
              width={100}
              height={100}
              className="relative"
            />
          </Box>
        </Grid>
      </Grid>
    ),
  },
  {
    title: "Simplify complexity, intelligent security",
    desc: "Easy to get started with, achieving flat management",
    icon: "/images/feature2-icon.png",
    content: (
      <Grid container position="relative" sx={{ mb: 5 }} ml={{ xs: 3, md: 0 }}>
        <Grid container>
          <Grid item xs={12} md={6} mt={7} order={{ xs: 2, md: 1 }}>
            <Box position="relative">
              <Image
                src="/images/feature2-bg.png"
                alt=""
                layout="responsive"
                width={100}
                height={100}
                className="absolute right-1/3"
                style={{ bottom: "10%" }}
              />
              <Image
                src="/images/feature2.svg"
                alt="开箱即用，轻松上手，适配多种运行环境"
                layout="responsive"
                width={100}
                height={100}
                className="relative"
                style={{ right: "43px" }}
              />
            </Box>
          </Grid>
          <Grid
            item
            xs={12}
            md={4}
            mt={{ xs: 3, md: 14 }}
            order={{ xs: 1, md: 2 }}
          >
            <List>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText
                  sx={{ textIndent: "-0.75rem" }}
                  primary="· &nbsp;One click installation, containerized management, adaptable to multiple operating environments"
                />
              </ListItem>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText
                  sx={{ textIndent: "-0.75rem" }}
                  primary="· &nbsp;Ready to use configuration out of the box, without the need for extensive adjustments to cumbersome rules"
                />
              </ListItem>
              <ListItem sx={{ mb: 1 }}>
                <ListItemText
                  sx={{ textIndent: "-0.75rem" }}
                  primary="· &nbsp;Simple operation, designed specifically for the community"
                />
              </ListItem>
            </List>
          </Grid>
        </Grid>
        <Grid
          item
          xs={0}
          md={2}
          sx={{
            display: { xs: "none", md: "block" },
            position: "absolute",
            right: 0,
          }}
        >
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
    title: "High performance, high concurrency, and high availability",
    desc: "Irregular engine, linear security detection algorithm",
    icon: "/images/feature3-icon.png",
    content: (
      <Grid
        container
        display="flex"
        justifyContent="flex-end"
        position="relative"
        sx={{ mb: 10 }}
        ml={{ xs: 3, md: 0 }}
      >
        <Grid
          item
          xs={0}
          md={2}
          sx={{
            display: { xs: "none", md: "block" },
            position: "absolute",
            left: 0,
          }}
        >
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
              <ListItemText
                sx={{ textIndent: "-0.75rem" }}
                primary="· &nbsp;Average detection delay <1 millisecond, single core easily detects 2000+TPS concurrency"
              />
            </ListItem>
            <ListItem>
              <ListItemText
                sx={{ textIndent: "-0.75rem" }}
                primary="· &nbsp;Developing a comprehensive health check mechanism based on Nginx"
              />
            </ListItem>
          </List>
        </Grid>
        <Grid item xs={12} md={6} mt={7}>
          <Box position={"relative"}>
            <Image
              src="/images/feature3-bg.png"
              alt=""
              width={1225}
              height={1310}
              className="absolute -left-1/2 -top-2/3"
            />
            <Image
              src="/images/ability/2000tps.png"
              alt="性能,服务可用性99.99%"
              layout="responsive"
              width={100}
              height={100}
              className="relative"
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
              display={{ xs: "block", md: index % 2 == 1 ? "flex" : "" }}
              alignItems="flex-end"
              flexDirection={index % 2 == 1 ? "column" : "row"}
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
                  <Box
                    component="span"
                    sx={{ mt: "5px", fontFamily: "GilroyBold" }}
                  >
                    {feature.title}
                  </Box>
                </Typography>
                <Typography
                  variant="subtitle1"
                  sx={{ color: alpha("#000", 0.7), fontWeight: 100, fontFamily: 'GilroyThin' }}
                  mt={2}
                >
                  <Box component="span">{feature.desc}</Box>
                </Typography>
              </Box>
            </Box>
            <Box sx={{ whiteSpace: "pre-line" }}>{feature.content}</Box>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Features;
