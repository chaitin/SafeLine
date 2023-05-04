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
    <Box sx={{ display: { xs: "block", sm: "flex" }, height: "100%" }}>
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
        {list.map((group) => (
          <Box sx={{ pl: "40px", lineHeight: "32px" }} key={group.title}>
            <Box sx={{ color: "text.auxiliary", mb: "6px", mt: "20px" }}>
              {group.title}
            </Box>
            {group.list.map((nav) => (
              <Box
                key={nav.title}
                component={Link}
                href={`/posts/${nav.id}`}
                sx={{
                  fontSize: "16px",
                  display: "block",
                  textDecoration: "none",
                  color: asPath.startsWith(`/posts/${nav.id}`)
                    ? "primary.main"
                    : "inherit",
                  fontWeight: asPath.startsWith(`/posts/${nav.id}`) ? 700 : 400,
                  "&:hover": {
                    color: "primary.main",
                  },
                }}
              >
                {nav.title}
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
            {list.map((group) => (
              <li key={`section-${group.title}`}>
                <ul>
                  <ListSubheader sx={{ color: "text.auxiliary" }}>
                    {group.title}
                  </ListSubheader>
                  {group.list.map((nav) => (
                    <ListItem
                      key={nav.title}
                      onClick={() => {
                        router.push(`/posts/${nav.id}`);
                        handleClick();
                      }}
                      sx={{
                        color: asPath.startsWith(`/posts/${nav.id}`)
                          ? "primary.main"
                          : "text.primary",
                        fontWeight: asPath.startsWith(`/posts/${nav.id}`)
                          ? 700
                          : 400,
                      }}
                    >
                      <ListItemText primary={nav.title} />
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
          width: { xs: "100%", sm: "calc(100% - 240px)" },
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
