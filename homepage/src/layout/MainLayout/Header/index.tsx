import React, { useMemo } from "react";
import { useRouter } from "next/router";
import Image from "next/image";
import Link from "next/link";
import {
  AppBar,
  Box,
  CssBaseline,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Toolbar,
  Typography,
  Button,
  alpha,
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
interface Props {
  /**
   * Injected by the documentation to work in an iframe.
   * You won't need it on your project.
   */
  window?: () => Window;
}

const drawerWidth = 240;

export default function DrawerAppBar(props: Props) {
  const { window } = props;
  const router = useRouter();
  const { pathname, asPath } = router;
  const [mobileOpen, setMobileOpen] = React.useState(false);

  const handleDrawerToggle = () => {
    setMobileOpen((prevState) => !prevState);
  };

  const drawer = (
    <Box onClick={handleDrawerToggle} sx={{ textAlign: "center" }}>
      <Typography
        variant="h6"
        component="div"
        sx={{
          flexGrow: 1,
          display: "flex",
          alignItems: "center",
          my: 2,
          ml: 2,
        }}
        onClick={() => router.push("/")}
      >
        <Image
          src="/images/logo.png"
          alt="Logo"
          width={120}
          height={34}
          priority
        />
        <Image
          src="/images/safeline.png"
          alt="Logo"
          width={34}
          height={34}
          style={{ marginLeft: "24px" }}
          priority
        />
      </Typography>
      <Divider />
      <List>
        <ListItem disablePadding>
          <ListItemButton
            sx={{ textAlign: "center" }}
            selected={pathname === "/"}
            component={Link}
            href="/"
          >
            <ListItemText primary="主页" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton
            sx={{ textAlign: "center" }}
            selected={pathname.startsWith("/posts/")}
            component={Link}
            href="/posts/guide_introduction/"
          >
            <ListItemText primary="技术文档" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton
            sx={{ textAlign: "center" }}
            component={Link}
            href="https://space.bilibili.com/521870525"
            target="_blank"
          >
            <ListItemText primary="教学视频" />
          </ListItemButton>
        </ListItem>

        <ListItem disablePadding>
          <ListItemButton
            sx={{ textAlign: "center" }}
            component={Link}
            href="https://demo.waf-ce.chaitin.cn:9443/"
            target="_blank"
          >
            <ListItemText primary="演示环境" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );

  const container =
    window !== undefined ? () => window().document.body : undefined;

  return (
    <Box
      sx={{
        display: "flex",
      }}
    >
      <CssBaseline />
      <AppBar
        component="nav"
        sx={{
          backgroundColor: "#0F1935",
          boxShadow: "none",
          color: "text.primary",
          pr: "0 !important",
        }}
      >
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Typography
            variant="h6"
            component="div"
            sx={{
              display: "flex",
              alignItems: "center",
            }}
            onClick={() => router.push("/")}
          >
            <Image
              src="/images/logo.png"
              alt="Logo"
              width={120}
              height={34}
              priority
            />

            <Image
              src="/images/403.svg"
              alt="Logo"
              width={34}
              height={34}
              style={{ marginLeft: "24px" }}
              priority
            />
          </Typography>
          <IconButton
            color="inherit"
            onClick={handleDrawerToggle}
            sx={{ display: { sm: "none", color: "#fff" } }}
          >
            <MenuIcon />
          </IconButton>
          <Box
            sx={{
              display: { xs: "none", sm: "block" },
              ".MuiButtonBase-root": {
                ml: "16px",
              },
              "a.MuiBox-root": {
                mr: 5,
                "&:hover": {
                  color: "primary.main",
                },
              },
            }}
          >
            <Box
              sx={{
                color: pathname === "/" ? "primary.main" : "#fff",
              }}
              component={Link}
              href="/"
            >
              主页
            </Box>
            <Box
              sx={{
                color: pathname.startsWith("/posts/") ? "primary.main" : "#fff",
              }}
              component={Link}
              href="/posts/guide_introduction/"
            >
              技术文档
            </Box>
            {/* <Box
              sx={{
                color: pathname.startsWith("/detection")
                  ? "primary.main"
                  : "#fff",
              }}
              component={Link}
              href="/detection"
            >
              安全效果
            </Box> */}
            <Box
              sx={{
                color: "#fff",
              }}
              component={Link}
              href="https://space.bilibili.com/521870525"
              target="_blank"
            >
              教学视频
            </Box>
            <Button
              sx={{ ml: "0 !important" }}
              component={Link}
              variant="contained"
              href="https://demo.waf-ce.chaitin.cn:9443/"
              target="_blank"
            >
              演示环境
            </Button>
          </Box>
        </Toolbar>
      </AppBar>
      <Box component="nav">
        <Drawer
          container={container}
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: "block", sm: "none" },
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>
      </Box>
    </Box>
  );
}
