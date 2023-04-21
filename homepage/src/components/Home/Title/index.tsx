import React from "react";
import { Typography, SxProps } from "@mui/material";
import Image from "next/image";

interface TitleProps {
  title: string;
  sx?: SxProps;
}

const Title: React.FC<TitleProps> = ({ title, sx }) => {
  return (
    <Typography
      variant="h4"
      sx={{ fontSize: "32px", fontWeight: 500, position: "relative", ...sx }}
    >
      <Image
        src="/images/class.png"
        alt="class"
        width={56}
        height={56}
        priority
        style={{ position: "absolute", top: -20, left: -24 }}
      />
      {title}
    </Typography>
  );
};

export default Title;
