import { useEffect } from "react";

export default function Hero({ scrollAnchorId }) {
  useEffect(() => {
    import("tiny-slider").then(({ tns }) => {
      tns({
        container: ".testimonial-active",
        autoplay: true,
        autoplayTimeout: 3000,
        autoplayButtonOutput: false,
        mouseDrag: true,
        gutter: 0,
        nav: false,
        navPosition: "bottom",
        controls: false,
        autoplayHoverPause: true,
        items: 1,
      });
      const testimonialDom = document.getElementsByClassName(
        "testimonial-active-wrapper"
      )?.[0];
      testimonialDom.style.display = "block";
    });
  }, []);
  return (
    <section id={scrollAnchorId} className="hero-section">
      <div className="container">
        <div className="row align-items-center">
          <div className="col-xl-6 col-lg-6 col-md-10">
            <div className="hero-content" style={{ color: "#000" }}>
              <h1
                style={{
                  fontFamily: "GilroyBold",
                  fontSize: 64,
                  marginBottom: "3px",
                }}
              >
                SafeLine
              </h1>
              <h2
                style={{
                  fontFamily: "GilroyBold",
                  fontSize: 28,
                  margin: "4px auto 14px auto",
                }}
              >
                The Best Free WAF For Webmaster
              </h2>
              <p style={{ letterSpacing: "1px", color: "rgba(0, 0, 0, 0.7)" }}>
                SafeLine is a simple, lightweight, locally deployable WAF that
                protects your website from network attacks that including OWASP
                attacks, zero-day attacks, web crawlers, vulnerability scanning,
                vulnerability exploit, http flood and so on.
              </p>
              <a
                href="https://docs.waf.chaitin.com/en/tutorials/install"
                className="main-btn btn-hover"
                target="_blank"
              >
                Get Started
              </a>
            </div>
          </div>

          <div className="col-xxl-6 col-xl-6 col-lg-6 testimonial-section">
            <div className="hero-image text-center text-lg-end">
              <div
                className="testimonial-active-wrapper"
                style={{
                  display: "none",
                  borderRadius: "12px",
                  overflow: "hidden",
                  boxShadow: "0 30px 50px rgba(145,158,171,0.3);",
                }}
              >
                <div className="testimonial-active">
                  {Array.from({ length: 4 }).map((_, index) => (
                    <img
                      key={index}
                      src={`/images/screenshot-${index}.png`}
                      className="single-testimonial"
                    />
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
