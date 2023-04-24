import "@/styles/globals.css";
import "@/styles/markdown.css";
import type { AppProps } from "next/app";
import { ThemeProvider } from "@/components";
import MainLayout from "@/layout/MainLayout";
import type { ReactElement, ReactNode } from "react";
import type { NextPage } from "next";
import Script from "next/script";

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const getLayout = Component.getLayout || ((page) => page);

  return (
    <>
      {process.env.NODE_ENV === "production" && (
        <Script
          type="text/javascript"
          src="https://v1.cnzz.com/z_stat.php?id=1281262430&web_id=1281262430"
        />
      )}
      <ThemeProvider>
        <MainLayout>{getLayout(<Component {...pageProps} />)}</MainLayout>
      </ThemeProvider>
    </>
  );
}
