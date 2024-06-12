import Image from "next/image";
import { contactInfo } from "./main-menu/main-menu";

export default function FooterSection({}) {
  const links = [].map((link) => ({
    ...link,
    url: link.url[0] === "#" ? `/${link.url}` : link.url,
  }));

  return (
    <footer className="footer pt-120">
      <div className="container">
        <div className="row">
          <div className="col-xl-3 col-lg-4 col-md-6 col-sm-10">
            <div className="footer-widget">
              <div className="logo">
                <img
                  src="/images/logo.png"
                  alt="Logo"
                  width={190}
                  height={46}
                  style={{
                    maxWidth: "100%",
                    height: "auto",
                  }}
                />
              </div>
              <p className="desc">The Best WAF for Webmaster</p>
              <ul className="social-links">
                {contactInfo.map((i) => (
                  <li key={i.icon}>
                    <a href={i.link} target="_blank">
                      <svg className="icon_svg" width="24px">
                        <use xlinkHref={'#'+i.icon} />
                      </svg>
                    </a>
                  </li>
                ))}

                <li>
                  <a href="https://github.com/chaitin/SafeLine" target="_blank">
                    <svg className="icon_svg" width="24px">
                      <use xlinkHref="#icon-github-fill" />
                    </svg>
                  </a>
                </li>
              </ul>
            </div>
          </div>
          <div className="col-xl-5 col-lg-4 col-md-12 col-sm-12 offset-xl-1">
            <div className="footer-widget">
              <h3>About Us</h3>
              <ul className="links">
                <li>
                  <a
                    href="https://docs.waf.chaitin.com/"
                    target="_blank"
                    style={{ textDecoration: "none", color: "unset" }}
                  >
                    Docs
                  </a>
                </li>
              </ul>
            </div>
          </div>

          <div className="col-xl-3 col-lg-4 col-md-6">
            <div className="footer-widget">
              <p>
                SafeLine is a simple, lightweight, locally deployable WAF that
                protects your website from network attacks that including OWASP
                attacks ...
              </p>
              {/* <form action="#">
                <input type="email" placeholder="Email" />
                <button className="main-btn btn-hover">Subscribe</button>
              </form> */}
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
