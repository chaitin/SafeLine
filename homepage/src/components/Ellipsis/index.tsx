import React, { type FC, useEffect, useRef, useState } from "react";

import { Box, type SxProps, Tooltip } from "@mui/material";

interface EllipsisProps {
  children: React.ReactNode;
  sx?: SxProps;
}

const Ellipsis: FC<EllipsisProps> = (props) => {
  const { children, sx } = props;
  const content = useRef<HTMLSpanElement>(null);
  const [ellipsis, setEllipsis] = useState<boolean>(false);

  //判断文字内容是否超出div宽度
  const onResize = () => {
    if (content.current) {
      setEllipsis(
        (content.current.parentNode! as HTMLDivElement).offsetWidth <
          content.current.offsetWidth
      );
    }
  };

  useEffect(() => {
    onResize();
  }, [children]);

  return (
    <>
      <Box
        sx={{
          width: "100%",
          overflow: "hidden",
          textOverflow: "ellipsis",
          whiteSpace: "nowrap",
          ...sx,
        }}
      >
        <Tooltip
          title={ellipsis ? children : null}
          arrow
          followCursor
          componentsProps={{
            tooltip: {
              sx: {
                color: "text.primary",
                backgroundColor: "background.paper0",
                boxShadow: " 0px 0px 20px 0px rgba(0,0,0,0.1)",
                padding: "16px",
                lineHeight: "24px",
                maxWidth: "500px",
                fontSize: "14px",
              },
            },
            arrow: { sx: { color: "background.paper0" } },
          }}
        >
          <Box ref={content} component="span" sx={{ maxWidth: "100%" }}>
            {children as React.ReactElement}
          </Box>
        </Tooltip>
      </Box>
    </>
  );
};

export default Ellipsis;
