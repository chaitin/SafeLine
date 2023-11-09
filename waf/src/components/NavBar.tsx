import React, { useEffect, useState } from 'react';
import { AppBar, Grid, Toolbar, Typography, Button, Box, Container, Link, Stack } from '@mui/material';
import Image from 'next/image';
import Icon from "@/components/Icon";
import usePopupState, { bindPopover, bindHover } from '@/components/Popover/usePopupState'
import HoverPopover from '@/components/Popover/HoverPopover'

const navs = [
  { to: "https://waf-ce.chaitin.cn/docs", label: "帮助文档", target: "_blank" },
  { to: "/community", label: "社区", target: "_self" },
  { to: "/version", label: "版本对比", target: "_self" },
];

export default function NavBar() {
  const [isSticky, setIsSticky] = useState(false);

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
            <Grid item xs={12} md={6} display="flex">
              <Box display="flex" alignItems="center">
                <Link href="/">
                  <Grid container flexDirection="row" display="flex" spacing={2} sx={{ marginTop: '0px' }}>
                    <Image
                      src="/images/safeline.svg"
                      alt="Waf Logo"
                      width={24}
                      height={26}
                    />
                    <Typography
                      variant="h6"
                      sx={{
                        ml: 1,
                        mr: 7,
                        display: 'flex',
                        alignItems: 'center',
                        color: "common.black",
                        fontFamily: "AlimamaShuHeiTi-Bold",
                      }}
                    >
                      雷池 SafeLine
                    </Typography>
                  </Grid>
                </Link>
                <Box display={{ xs: 'none', md: 'flex' }} alignItems="center">
                  {navs.map((nav, index) => (
                    <Box component="span" key={index} mr={3.5}>
                      <Link key={index} href={nav.to} sx={{ color: "common.black" }} target={nav.target}>{nav.label}</Link>
                    </Box>
                  ))}
                </Box>
              </Box>
            </Grid>
            <Grid item xs={0} md={6} display={{ xs: 'none', md: 'flex' }} justifyContent={'flex-end'}>
              <Box sx={{ fontSize: "16px", display: "flex", alignItems: "center" }}>
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
                <Link href="https://demo.waf-ce.chaitin.cn:9443/dashboard" sx={{ color: "common.black" }} mr={3.5} target="_blank">演示 demo</Link>
                <Button
                  variant="contained"
                  target="_blank"
                  sx={{ width: { xs: "100%", sm: "auto" } }}
                  href="https://waf-ce.chaitin.cn/docs/guide/install"
                >
                  立即安装
                </Button>
              </Box>
            </Grid>
          </Grid>
        </Toolbar>
      </Container>
    </AppBar>
  );
}
