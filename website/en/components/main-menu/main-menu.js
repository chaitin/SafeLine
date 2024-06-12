import { Tooltip } from "@mui/material";
import { useEffect, useRef, useState } from "react";

const contactInfo = [
  {
    title: "Discord",
    icon: "icon-icon-17",
    link: "https://discord.gg/wyshSVuvxC",
  },
  {
    title: "x.com",
    icon: "icon-tuite3",
    link: "https://x.com/safeline_waf",
  },
  {
    title: "Telegram",
    icon: "icon-telegram1",
    link: "https://t.me/safeline_waf",
  },
];
export default function ManiMenu() {
  function highlightLinks() {
    const sections = document.querySelectorAll(".page-scroll");
    const scrollPos =
      window.pageYOffset ||
      document.documentElement.scrollTop ||
      document.body.scrollTop;

    sections.forEach((currLink) => {
      const val = currLink.getAttribute("href").slice(1);
      if (val[0] !== "#") {
        return;
      }
      const refElement = document.querySelector(val);

      if (!refElement) {
        return;
      }

      const scrollTopMinus = scrollPos + 73;

      if (
        refElement.offsetTop <= scrollTopMinus &&
        refElement.offsetTop + refElement.offsetHeight > scrollTopMinus
      ) {
        setActiveMenuLink(val);
      }
    });
  }

  useEffect(() => {
    window.addEventListener("scroll", highlightLinks);

    return () => {
      window.removeEventListener("scroll", highlightLinks);
    };
  }, []);

  const [isMenuActive, setMenuActive] = useState(false);
  const menuLinksEl = useRef(null);

  function inactivateMenu() {
    setMenuActive(false);
    if (menuLinksEl.current) {
      menuLinksEl.current.classList.remove("show");
    }
  }

  return (
    <>
      <button
        className={`navbar-toggler ${isMenuActive ? "active" : ""}`}
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarSupportedContent"
        aria-cpontrols="navbarSupportedContent"
        aria-expanded="false"
        aria-label="Toggle navigation"
        onClick={() => setMenuActive(!isMenuActive)}
      >
        <span className="toggler-icon"></span>
        <span className="toggler-icon"></span>
        <span className="toggler-icon"></span>
      </button>

      <div
        className="collapse navbar-collapse sub-menu-bar"
        ref={menuLinksEl}
        id="navbarSupportedContent"
      >
        <div className="ms-auto">
          <ul id="nav" className="navbar-nav ms-auto">
            <li className="nav-item">
              <a href="https://docs.waf.chaitin.com/" target="_blank">
                Docs
              </a>
            </li>
            <li
              className="nav-item nav-item_tooltip"
              style={{ position: "relative" }}
            >
              <a style={{ color: "rgba(0,0,0,0.2)" }}>Pricing</a>
              <div className="nav-btn_tooltip">Comming soon...</div>
            </li>
            <li
              className="nav-item"
              style={{ display: "flex", marginLeft: "17px",marginRight: '2px', }}
            >
              {contactInfo.map((i) => (
                <a
                  key={i.icon}
                  target="_blank"
                  className="nav-item_icon"
                  href={i.link}
                  style={{ padding: "4px" }}
                >
                  <Tooltip title={i.title}>
                    <svg
                      className="icon_svg"
                      style={{
                        width: "18px",
                        height: "18px",
                        marginRight: "4px",
                      }}
                    >
                      <use xlinkHref={"#" + i.icon} />
                    </svg>
                  </Tooltip>
                </a>
              ))}
            </li>

            <li className="nav-item">
              <a
                style={{ paddingLeft: 0 }}
                className="nav-item_icon"
                href="https://github.com/chaitin/SafeLine"
                target="_blank"
              >
                <svg
                  className="icon_svg"
                  style={{ width: "20px", height: "20px" }}
                >
                  <use xlinkHref="#icon-github-fill" />
                </svg>
                GitHub 10k+
                <svg
                  className="icon_svg"
                  style={{ width: "16px", height: "16px", marginLeft: "8px" }}
                >
                  <use xlinkHref="#icon-xingxing1" />
                </svg>
              </a>
            </li>
            <li className="nav-item nav-item__demo">
              <a
                href={"https://demo.waf.chaitin.com:9443/"}
                className="main-btn btn-hover"
                style={{ color: "#fff" }}
                target="_blank"
              >
                Demo
              </a>
            </li>
          </ul>
        </div>
      </div>
    </>
  );
}
