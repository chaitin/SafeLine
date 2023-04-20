import "@/styles/globals.css";
import "@/styles/markdown.css";
import type { AppProps } from "next/app";
import { ThemeProvider } from "@/components";
import MainLayout from "@/layout/MainLayout";
import type { ReactElement, ReactNode } from "react";
import type { NextPage } from "next";

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const getLayout = Component.getLayout || ((page) => page);

  return (
    <ThemeProvider>
      <MainLayout>{getLayout(<Component {...pageProps} />)}</MainLayout>
    </ThemeProvider>
  );
}
