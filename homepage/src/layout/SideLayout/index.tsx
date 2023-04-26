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
  ListItemIcon,
  ListItemText,
  ListSubheader,
  ListItem,
} from "@mui/material";
import { type GroupItem } from "@/utils/posts";
interface SideLayoutProps {
  list: GroupItem[];
  children?: React.ReactNode;
}

const SideLayout: FC<SideLayoutProps> = ({ children, list }) => {
  const router = useRouter();
  const { asPath } = router;
  const [open, setOpen] = useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <Box sx={{ display: { xs: "block", sm: "flex" }, height:'100%' }}>
      <Box
        sx={{
          position: "fixed",
          top: 0,
          left: 0,
          width: 240,
          pt: "80px",
          height: "100%",
          backgroundColor: "#f8f9fc",
          // boxShadow: "inset 0px 0px 16px 0px rgba(0, 145, 255, 1)",
          display: { xs: "none", sm: "block" },
          borderRight: "1px solid hsla(210, 18%, 87%, 1)",
        }}
      >
        {list.map((menu) => (
          <Box sx={{ pl: "40px", lineHeight: "32px" }} key={menu.category}>
            <Box sx={{ color: "text.auxiliary", mb: "6px", mt: "20px" }}>
              {menu.category}
            </Box>
            {menu.list.map((item) => (
              <Box
                key={item.title}
                component={Link}
                href={`/posts/${item.id}`}
                sx={{
                  fontSize: "16px",
                  display: "block",
                  textDecoration: "none",
                  color: asPath.startsWith(`/posts/${item.id}`)
                    ? "primary.main"
                    : "inherit",
                  fontWeight: asPath.startsWith(`/posts/${item.id}`)
                    ? 700
                    : 400,
                  "&:hover": {
                    color: "primary.main",
                  },
                }}
              >
                {item.title}
              </Box>
            ))}
          </Box>
        ))}
      </Box>

      <AppBar
        sx={{
          top: 72,
          backgroundColor: "#fff",
          display: {
            xs: "block",
            sm: "none",
            boxShadow: "0 12px 25px -12px rgba(93,99,112, 0.2)",
          },
        }}
      >
        <ListItem onClick={handleClick}>
          <ListItemIcon>{open ? <ExpandLess /> : <ExpandMore />}</ListItemIcon>
          <ListItemText primary="目录" sx={{ color: "#000" }} />
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
            {list.map((menu) => (
              <li key={`section-${menu.category}`}>
                <ul>
                  <ListSubheader sx={{ color: "text.auxiliary" }}>
                    {menu.category}
                  </ListSubheader>
                  {menu.list.map((item) => (
                    <ListItem
                      key={item.title}
                      onClick={() => {
                        router.push(`/posts/${item.id}`);
                        handleClick();
                      }}
                      sx={{
                        color: asPath.startsWith(`/posts/${item.id}`)
                          ? "primary.main"
                          : "text.primary",
                        fontWeight: asPath.startsWith(`/posts/${item.id}`)
                          ? 700
                          : 400,
                      }}
                    >
                      <ListItemText primary={item.title} />
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
          pt: { xs: 18, sm: 13 },
          pb: 8,
          px: 3,
          flexGrow: 1,
          height: "100%",
          marginLeft: { xs: 0, sm: "240px" },
          backgroundColor: "#F8F9FC",
        }}
      >
        {children}
      </Box>
    </Box>
  );
};

export default SideLayout;
