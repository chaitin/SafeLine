import Head from "next/head";
import { useEffect } from "react";
import CssBaseline from '@mui/material/CssBaseline';

import "@/css/lineicons.css";
import "bootstrap/dist/css/bootstrap.css";

import "@/css/main.css";
import "@/css/tiny-slider.min.css";

function MyApp({ Component }) {
  useEffect(() => {
    import("bootstrap/dist/js/bootstrap.js");
  }, []);
  return (
    <>
      <Head>
        <meta charSet="utf-8" />
        <meta httpEquiv="x-ua-compatible" content="ie=edge" />
        <title>SafeLine | the Best WAF for Webmaster</title>
        <meta
          name="description"
          content="SafeLine is a simple, lightweight, locally deployable WAF that protects your website from network attacks that including OWASP attacks, zero-day attacks, web crawlers, vulnerability scanning, vulnerability exploit, http flood and so on."
        />
        <meta
          name="keywords"
          content="waf,safeline,free,open source,sql injection,xss"
        ></meta>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <script
          src="https://at.alicdn.com/t/c/font_4031246_8o6dgtp7wm5.js?spm=a313x.manage_type_myprojects.i1.13.5f303a813w1hKj&file=font_4031246_8o6dgtp7wm5.js"
          defer
        />
        <script
          async
          src="https://www.googletagmanager.com/gtag/js?id=G-Z48W47MR7B"
        ></script>
        <script async src="/ga.js"></script>

        <link rel="shortcut icon" type="image/x-icon" href="/favicon.png" />
      </Head>
      <CssBaseline />
      <Component />
    </>
  );
}

MyApp.getInitialProps = async (appContext) => {
  let mainMenu = [];

  return { ...{}, mainMenu };
};

export default MyApp;
