import React, { type FC } from "react";

import ErrorIcon from "@mui/icons-material/Error";
import { Box } from "@mui/material";

import { ThemeProvider } from "@/components";

import Modal, { type ModalProps } from "./Modal";

export interface ConfirmDialogProps extends ModalProps {
  content?: React.ReactNode;
  width?: string;
}

const ConfirmDialog: FC<ConfirmDialogProps> = (props) => {
  const { title = "提示", content, width = "480px", ...rest } = props;
  return (
    <ThemeProvider>
      <Modal
        title={
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              lineHeight: "22px",
              color: "text.main",
              fontWeight: 500,
            }}
          >
            <ErrorIcon
              sx={{ color: "#FFBF00", mr: "16px", fontSize: "24px" }}
            />
            {title}
          </Box>
        }
        closable={false}
        {...rest}
        sx={{ width }}
      >
        <Box sx={{ color: "text.main", pl: "40px" }}>{content}</Box>
      </Modal>
    </ThemeProvider>
  );
};

export default ConfirmDialog;
