import React, { useEffect, useState } from "react";
import {
  AppBar,
  Drawer,
  Grid,
  Toolbar,
  Typography,
  Button,
  Box,
  Container,
  Link,
  List,
  ListItem,
  ListItemText,
  Stack,
  IconButton,
  MenuItem,
  Select,
} from "@mui/material";
import Image from "next/image";
import dynamic from "next/dynamic";
import Icon from "@/components/Icon";
import CloseIcon from "@mui/icons-material/Close";
import usePopupState, {
  bindPopover,
  bindHover,
} from "@/components/Popover/usePopupState";

const navs = [
  { to: "/docs", label: "帮助文档", target: "_blank" },
  { to: "/community", label: "开发计划", target: "_self" },
  { to: "/version", label: "付费版本", target: "_self" },
];

const menus = [
  ...navs,
  {
    to: "https://github.com/chaitin/SafeLine",
    label: "GitHub",
    target: "_blank",
  },
  {
    to: "https://demo.waf-ce.chaitin.cn:9443/dashboard",
    label: "演示 Demo",
    target: "_blank",
  },
];

const HoverPopover = dynamic(
  () => import("@/components/Popover/HoverPopover"),
  {
    ssr: false,
  }
);

export default function NavBar() {
  const [isSticky, setIsSticky] = useState(false);
  const [open, setOpen] = useState(false);
  const [langOpen, setLangOpen] = useState(false);

  const handleOpen = () => {
    setLangOpen(true);
  };
  const handleClose = () => {
    setLangOpen(false);
  };
  const handleChange = () => {
    window.open("https://waf.chaitin.com/");
  };
  const popoverState = usePopupState({
    popupId: "wechat-qrcode-popover",
  });

  const popoverState1 = usePopupState({
    popupId: "bounty-popover",
  });

  useEffect(() => {
    const handleScroll = () => {
      const scrollTop = window.pageYOffset;
      setIsSticky(scrollTop > 0);
    };
    handleScroll();
    window.addEventListener("scroll", handleScroll);

    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);
  const langRender = () => (
    <Select
      value={"cn"}
      label=""
      onChange={handleChange}
      size="small"
      open={langOpen}
      onMouseEnter={handleOpen}
      onClose={handleClose}
      onOpen={handleOpen}
      MenuProps={{
        PaperProps: {
          onMouseLeave: handleClose,
        },
      }}
      renderValue={() => {
        return (
          <Stack
            direction="row"
            alignItems="center"
            spacing={1}
            sx={{ ml: "auto" }}
          >
            <Box display={{ xs: "none", md: "flex" }}>
              <svg
                className="icon_svg"
                style={{ width: "16px", height: "16px" }}
              >
                <use xlinkHref="#icon-diqiuyangshi1" />
              </svg>
            </Box>
            <Box>中文</Box>
          </Stack>
        );
      }}
      sx={{
        border: "none",
        fontFamily: "GilroyBold",
        px: "4px",
        ml: { sx: "auto", md: 1 },
        "& fieldset": {
          border: langOpen
            ? "2px solid rgba(15,198,194,0.1)!important"
            : "none",
        },
        "& .MuiSelect-select": { px: "12px" },
      }}
      className="lang_select"
    >
      <MenuItem
        value="en"
        sx={{
          "&.Mui-selected": {
            bgcolor: "rgba(15,198,194,0.1)!important",
          },
          "&:hover": { bgcolor: "rgba(15,198,194,0.1)" },
        }}
      >
        English
      </MenuItem>
      <MenuItem
        value="cn"
        sx={{
          "&.Mui-selected": {
            bgcolor: "rgba(15,198,194,0.1)!important",
          },
          "&:hover": { bgcolor: "rgba(15,198,194,0.1)" },
        }}
      >
        简体中文
      </MenuItem>
    </Select>
  );
  return (
    <>
      <AppBar
        position="fixed"
        color="transparent"
        sx={{
          boxShadow: "none",
          ...(isSticky
            ? {
                backdropFilter: "blur(8px)",
                background: "rgba(255,255,255,0.8)",
              }
            : {}),
          alignItems: "center",
        }}
      >
        <Container maxWidth="lg" sx={{ mx: 0 }}>
          <Toolbar
            sx={{ backgroundColor: "transparent", py: 1, ml: 2, px: { sm: 0 } }}
          >
            <Grid container justifyContent="space-around">
              <Grid item xs={10} md={6} display="flex">
                <Box display="flex" alignItems="center" sx={{ width: "100%" }}>
                  <SafelineTitle />
                  {langRender()}
                  <Box display={{ xs: "none", md: "flex" }} alignItems="center">
                    {navs.map((nav, index) => (
                      <Box component="span" key={index} mr={3.5}>
                        <Link
                          key={index}
                          href={nav.to}
                          sx={{ color: "common.black" }}
                          target={nav.target}
                        >
                          {nav.label}
                        </Link>
                      </Box>
                    ))}
                  </Box>
                </Box>
              </Grid>
              <Grid item xs={2} md={6} display="flex" justifyContent="flex-end">
                <Box
                  sx={{
                    fontSize: "16px",
                    display: { xs: "none", md: "flex" },
                    alignItems: "center",
                  }}
                >
                  <Link
                    href="https://discord.gg/wyshSVuvxC"
                    target="_blank"
                    sx={{
                      color: "common.black",
                      display: "flex",
                      "&:hover": {
                        color: "primary.main",
                      },
                    }}
                    mr={3.5}
                  >
                    <Icon type="icon-discord1" sx={{ mr: 1 }} />
                    Discord
                  </Link>
                  <Link
                    href="https://github.com/chaitin/SafeLine"
                    target="_blank"
                    sx={{
                      color: "common.black",
                      display: "flex",
                      "&:hover": {
                        color: "primary.main",
                      },
                    }}
                    mr={3.5}
                  >
                    <Icon type="icon-github-fill" sx={{ mr: 1 }} />
                    GitHub
                  </Link>
                  <Box mr={3.5}>
                    <Typography
                      {...bindHover(popoverState as any)}
                      variant="body1"
                      sx={{
                        color: popoverState.isOpen
                          ? "primary.main"
                          : "common.black",
                        cursor: "pointer",
                        "&:hover": {
                          color: "primary.main",
                          backgroundColor: "transparent",
                        },
                        transition: "unset",
                      }}
                    >
                      微信讨论组
                    </Typography>
                    <HoverPopover
                      {...bindPopover(popoverState)}
                      anchorOrigin={{
                        vertical: 42,
                        horizontal: "left",
                      }}
                      transformOrigin={{
                        vertical: "top",
                        horizontal: "left",
                      }}
                      marginThreshold={16}
                    >
                      <Image
                        src="/images/wechat-230825.png"
                        alt="wechat"
                        width={259}
                        height={259}
                      />
                    </HoverPopover>
                  </Box>
                  <Button
                    variant="contained"
                    target="_blank"
                    sx={{
                      width: { xs: "100%", sm: "auto" },
                      textTransform: "none",
                    }}
                    href="https://demo.waf-ce.chaitin.cn:9443/dashboard"
                  >
                    Demo
                  </Button>
                </Box>
                <Stack justifyContent="center">
                  <Icon
                    type="icon-caidan"
                    color="common.black"
                    sx={{
                      display: { xs: "flex", md: "none" },
                      fontSize: "26px",
                    }}
                    onClick={() => setOpen(true)}
                  />
                </Stack>
              </Grid>
            </Grid>
          </Toolbar>
        </Container>
      </AppBar>
      <Drawer
        anchor="right"
        sx={{ width: "100%" }}
        variant="temporary"
        open={open}
        PaperProps={{
          style: {
            width: "100%",
          },
        }}
        onClose={() => setOpen(false)}
      >
        <Stack
          direction="row"
          justifyContent="space-between"
          pl={4}
          pr={0.5}
          py={1}
          sx={{ boxShadow: "rgba(0, 0, 0, 0.1) 0px 1px 2px 0px" }}
        >
          <Box>
            <SafelineTitle />
          </Box>
          <IconButton onClick={() => setOpen(false)}>
            <CloseIcon />
          </IconButton>
        </Stack>
        <List>
          {menus.map((menu) => (
            <Link key={menu.label} href={menu.to} target={menu.target}>
              <ListItem>
                <ListItemText primary={menu.label} />
              </ListItem>
            </Link>
          ))}
        </List>
        <Box ml={2}>
          <Image
            src="/images/wechat-230825.png"
            alt="wechat"
            width={160}
            height={160}
          />
        </Box>
      </Drawer>
    </>
  );
}

export const SafelineTitle: React.FC = () => {
  return (
    <Link href="/">
      <Grid
        container
        flexDirection="row"
        display="flex"
        spacing={2}
        sx={{ marginTop: "0px", flexWrap: "nowrap" }}
      >
        <Box
          width={{ xs: "30px", md: "24px" }}
          height={{ xs: "33px", md: "26px" }}
          position="relative"
        >
          <Image
            src="/images/safeline.svg"
            alt="SafeLine Logo"
            layout="responsive"
            width={30}
            height={33}
          />
        </Box>
        <Typography
          variant="h6"
          sx={{
            ml: { xs: "4px", md: 1 },
            mr: { xs: 0, md: 0 },
            fontSize: { xs: "20px", md: "16px" },
            display: "flex",
            alignItems: "center",
            fontFamily: "AlimamaShuHeiTi-Bold",
            fontWeight: 500,
          }}
        >
          雷池 SafeLine
        </Typography>
      </Grid>
    </Link>
  );
};
