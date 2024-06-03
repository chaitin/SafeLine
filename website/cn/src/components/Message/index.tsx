import type React from "react";

import { type AlertColor } from "@mui/material";

import Notification from "./Message";

type MessageStaticFunctions = Record<
  AlertColor,
  (context: React.ReactNode, duration?: number) => void
>;

const Message = {} as MessageStaticFunctions;

let notification: any = null;

if (typeof window !== "undefined") {
  // @ts-ignore
  Notification.newInstance({}, (n: any) => {
    notification = n;
  });

  const commonOpen =
    (type: AlertColor) =>
    (content: React.ReactNode, duration: number = 3) => {
      notification.notice({
        duration,
        severity: type,
        content,
      });
    };

  (["success", "warning", "info", "error"] as const).forEach((type) => {
    Message[type] = commonOpen(type);
  });
}

export default Message;
