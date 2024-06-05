import { Box, MenuItem, Stack } from "@mui/material";
import Select from "@mui/material/Select";
import { useEffect, useRef, useState } from "react";

import MainMenu from "./main-menu/main-menu";

export default function HeaderSection() {
  const [isNavbarSticky, setIsNavbarSticky] = useState(false);
  const [langOpen, setLangOpen] = useState(false);
  const navbarAreaEl = useRef(null);
  const [lang, setLang] = useState("en");
  const handleOpen = () => {
    setLangOpen(true);
  };
  const handleClose = () => {
    setLangOpen(false);
  };
  const handleChange = (event) => {
    window.open("https://waf-ce.chaitin.cn/");
  };

  function fixNavBar() {
    if (navbarAreaEl.current) {
      setIsNavbarSticky(
        document.getElementById("next_container").scrollTop >
          navbarAreaEl.current.offsetTop
      );
    }
  }

  useEffect(() => {
    if (document.getElementById("next_container"))
      document
        .getElementById("next_container")
        .addEventListener("scroll", fixNavBar);

    return () => {
      document
        .getElementById("next_container")
        .removeEventListener("scroll", fixNavBar);
    };
  }, []);

  return (
    <header className="header">
      <div
        ref={navbarAreaEl}
        className={`navbar-area ${isNavbarSticky ? "sticky" : ""}`}
      >
        <div className="container">
          <div className="row align-items-center">
            <div className="col-lg-12">
              <nav className="navbar navbar-expand-lg">
                <img
                  src="/images/logo.png"
                  alt="Logo"
                  width={190}
                  height={46}
                  style={{
                    maxWidth: "100%",
                    height: "auto",
                    marginRight: "24px",
                  }}
                />

                <Select
                  value={lang}
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
                  renderValue={(v) => {
                    return (
                      <Stack direction="row" alignItems="center" spacing={1}>
                        <svg
                          className="icon_svg"
                          style={{ width: "16px", height: "16px" }}
                        >
                          <use xlinkHref="#icon-diqiuyangshi1" />
                        </svg>
                        <Box>EN</Box>
                      </Stack>
                    );
                  }}
                  sx={{
                    border: "none",
                    fontFamily: "GilroyBold",
                    px: "4px",
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
                <MainMenu />
              </nav>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
