import React from "react";
import Image from 'next/image'
import { Typography, SxProps, Grid, Link } from "@mui/material";

interface TitleProps {
  title: string;
  sx?: SxProps;
}

const Title: React.FC<TitleProps> = ({ title, sx }) => {
  return (
    <Link href="/">
      <Grid container flexDirection="row" display="flex" alignItems="center" sx={{ marginTop: 0 }}>
        <Image
          src="/images/safeline.svg"
          alt="SafeLine Logo"
          width={40}
          height={43}
        />
        <Typography
          variant="h4"
          sx={{ color: "common.white", fontFamily: "AlimamaShuHeiTi-Bold", ...sx }}
        >
          {title}
        </Typography>
      </Grid>
    </Link>
  );
};

export default Title;
