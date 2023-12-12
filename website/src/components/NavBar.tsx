import React, { useEffect, useState } from 'react';
import { AppBar, Drawer, Grid, Toolbar, Typography, Button, Box, Container, Link, List, ListItem, ListItemText, Stack, IconButton } from '@mui/material';
import Image from 'next/image';
import dynamic from 'next/dynamic';
import Icon from "@/components/Icon";
import CloseIcon from '@mui/icons-material/Close';
import usePopupState, { bindPopover, bindHover } from '@/components/Popover/usePopupState'

const navs = [
  { to: "/docs", label: "帮助文档", target: "_blank" },
  { to: "/community", label: "开发计划", target: "_self" },
  { to: "/version", label: "付费版本", target: "_self" },
];

const menus = [
  ...navs,
  { to: "https://github.com/chaitin/SafeLine", label: "GitHub", target: "_blank" },
  { to: "https://demo.waf-ce.chaitin.cn:9443/dashboard", label: "演示 Demo", target: "_blank" },
];

const HoverPopover = dynamic(() => import('@/components/Popover/HoverPopover'), {
  ssr: false,
});

export default function NavBar() {
  const [isSticky, setIsSticky] = useState(false);
  const [open, setOpen] = useState(false);

  const popoverState = usePopupState({
    popupId: "wechat-qrcode-popover"
  })

  useEffect(() => {
    const handleScroll = () => {
      const scrollTop = window.pageYOffset;
      setIsSticky(scrollTop > 0);
    };
    handleScroll()
    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  return (
    <>
      <AppBar position='fixed' color='transparent'
        sx={{
          boxShadow: 'none',
          ...(isSticky ? { backdropFilter: 'blur(8px)', background: 'rgba(255,255,255,0.8)' }  : {}),
          alignItems: 'center'
        }}
      >
        <Container maxWidth="lg" sx={{ mx: 0 }}>
          <Toolbar sx={{ backgroundColor: 'transparent', py: 1, ml: 2, px: { sm: 0 } }}>
            <Grid container justifyContent="space-around">
              <Grid item xs={10} md={6} display="flex">
                <Box display="flex" alignItems="center">
                  <SafelineTitle />
                  <Box display={{ xs: 'none', md: 'flex' }} alignItems="center">
                    {navs.map((nav, index) => (
                      <Box component="span" key={index} mr={3.5}>
                        <Link key={index} href={nav.to} sx={{ color: "common.black" }} target={nav.target}>{nav.label}</Link>
                      </Box>
                    ))}
                  </Box>
                </Box>
              </Grid>
              <Grid item xs={2} md={6} display="flex" justifyContent="flex-end">
                <Box sx={{ fontSize: "16px", display: { xs: "none", md: "flex" }, alignItems: "center" }}>
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
                      variant='body1'
                      sx={{
                        color: popoverState.isOpen ? 'primary.main' : 'common.black',
                        cursor: "pointer",
                        '&:hover': {
                          color: "primary.main",
                          backgroundColor: "transparent",
                        },
                        transition: 'unset',
                      }}
                    >
                      微信讨论组
                    </Typography>
                    <HoverPopover
                      {...bindPopover(popoverState)}
                      anchorOrigin={{
                        vertical: 42,
                        horizontal: 'left',
                      }}
                      transformOrigin={{
                        vertical: 'top',
                        horizontal: 'left',
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
                  <Link href="https://demo.waf-ce.chaitin.cn:9443/dashboard" sx={{ color: "common.black" }} mr={3.5} target="_blank">演示 Demo</Link>
                  <Button
                    variant="contained"
                    target="_blank"
                    sx={{ width: { xs: "100%", sm: "auto" } }}
                    href="/docs/guide/install"
                  >
                    立即安装
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
        anchor='right'
        sx={{ width: "100%" }}
        variant="temporary"
        open={open}
        PaperProps={{
          style: {
            width: '100%',
          },
        }}
        onClose={() => setOpen(false)}
      >
        <Stack direction="row" justifyContent="space-between" pl={4} pr={0.5} py={1} sx={{ boxShadow: "rgba(0, 0, 0, 0.1) 0px 1px 2px 0px" }}>
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
      <Grid container flexDirection="row" display="flex" spacing={2} sx={{ marginTop: '0px', minWidth: "192px" }}>
        <Box width={{ xs: "40px", md: "24px" }} height={{ xs: "43px", md: "26px" }} position="relative">
          <Image
            src="/images/safeline.svg"
            alt="SafeLine Logo"
            layout="responsive"
            width={40}
            height={43}
          />
        </Box>
        <Typography
          variant="h6"
          sx={{
            ml: { xs: 2, md: 1 },
            mr: { xs: 0, md: 7 },
            fontSize: { xs: "24px", md: "16px" },
            display: 'flex',
            alignItems: 'center',
            fontFamily: "AlimamaShuHeiTi-Bold",
          }}
        >
          雷池 SafeLine
        </Typography>
      </Grid>
    </Link>
  );
};
