import { type FC, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";
import {
  Box,
  Collapse,
  AppBar,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  ListSubheader,
  ListItem,
} from "@mui/material";

interface IProps {
  children?: React.ReactNode;
}

const MENU_LIST = [
  {
    group: "上手指南",
    list: [
      {
        name: "产品介绍",
        url: "/posts/introduction",
      },
      {
        name: "快速部署",
        url: "/posts/install",
      },
      {
        name: "常见问题排查",
        url: "/posts/faq",
      },
    ],
  },
];

const SideLayout: FC<IProps> = ({ children }) => {
  const router = useRouter();
  const { asPath } = router;
  const [open, setOpen] = useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <Box sx={{ display: { xs: "block", sm: "flex" } }}>
      <Box
        sx={{
          position: "fixed",
          top: 0,
          left: 0,
          width: 240,
          pt: "120px",
          height: "100%",
          backgroundColor: "background.paper",
          display: { xs: "none", sm: "block" },
        }}
      >
        {MENU_LIST.map((menu) => (
          <Box sx={{ pl: "40px", lineHeight: "32px" }} key={menu.group}>
            <Box sx={{ color: "text.auxiliary", mb: "6px" }}>{menu.group}</Box>
            {menu.list.map((item) => (
              <Box
                key={item.name}
                component={Link}
                href={item.url}
                sx={{
                  fontSize: "16px",
                  display: "block",
                  textDecoration: "none",
                  color: asPath.startsWith(item.url)
                    ? "primary.main"
                    : "inherit",
                  fontWeight: asPath.startsWith(item.url) ? 700 : 400,
                  "&:hover": {
                    color: "primary.main",
                  },
                }}
              >
                {item.name}
              </Box>
            ))}
          </Box>
        ))}
      </Box>

      <AppBar sx={{ top: 72, display: { xs: "block", sm: "none" } }}>
        <ListItem onClick={handleClick}>
          <ListItemIcon>{open ? <ExpandLess /> : <ExpandMore />}</ListItemIcon>
          <ListItemText primary="目录" />
        </ListItem>
        <Collapse in={open} timeout="auto" unmountOnExit>
          <List
            sx={{
              width: "100%",
              bgcolor: "background.paper",
              position: "relative",
              overflow: "auto",
              maxHeight: 300,
              "& ul": { padding: 0 },
            }}
            subheader={<li />}
          >
            {MENU_LIST.map((menu) => (
              <li key={`section-${menu.group}`}>
                <ul>
                  <ListSubheader sx={{ color: "text.auxiliary" }}>
                    {menu.group}
                  </ListSubheader>
                  {menu.list.map((item) => (
                    <ListItem
                      key={item.name}
                      onClick={() => {
                        router.push(item.url);
                        handleClick();
                      }}
                      sx={{
                        color: asPath.startsWith(item.url)
                          ? "primary.main"
                          : "text.primary",
                        fontWeight: asPath.startsWith(item.url) ? 700 : 400,
                      }}
                    >
                      <ListItemText primary={item.name} />
                    </ListItem>
                  ))}
                </ul>
              </li>
            ))}
          </List>
        </Collapse>
      </AppBar>

      <Box
        sx={{
          p: { xs: 8, sm: 3 },
          flexGrow: 1,
          height: "100%",
          marginLeft: { xs: 0, sm: "240px" },
        }}
      >
        {children}
      </Box>
    </Box>
  );
};

export default SideLayout;
