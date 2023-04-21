import { type FC, useEffect } from "react";

import { Box, type SxProps } from "@mui/material";

interface IconProps {
  type: string;
  sx?: SxProps;
  [propName: string]: any;
}

const Icon: FC<IconProps> = ({ type, sx, ...restProps }) => {
  useEffect(() => {
    require("../../static/fonts/iconfont");
  }, []);
  return (
    // @ts-ignore
    <Box
      component="svg"
      sx={{ width: "1em", height: "1em", fill: "currentColor", ...sx }}
      aria-hidden="true"
      {...restProps}
    >
      <use xlinkHref={`#${type}`} />
    </Box>
  );
};
export default Icon;
