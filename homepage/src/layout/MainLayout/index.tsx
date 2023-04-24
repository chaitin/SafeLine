import { type FC } from "react";
import Head from "next/head";
import { Box } from "@mui/material";

import Header from "./Header";

interface IProps {
  redirect?: string;
  children?: React.ReactNode;
}
const MainLayout: FC<IProps> = ({ children }) => {
  return (
    <>
      <Head>
        <title>长亭雷池 WAF</title>
        <meta
          name="description"
          content="长亭雷池 Web 应用防火墙  | 长亭雷池 WAF"
        />
        <meta name="keywords" content="WAF,雷池,长亭,社区版" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/images/safeline.png" />
      </Head>

      <Box
        sx={{
          position: "relative",
          flex: 1,
          display: "flex",
          flexDirection: "column",
          height: "100vh",
        }}
      >
        <Header />
        <Box component="main" sx={{ flexGrow: 1 }}>
          {children}
        </Box>
      </Box>
    </>
  );
};

export default MainLayout;
