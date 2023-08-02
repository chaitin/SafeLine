import React from "react";
import { Typography, Box, Button, alpha } from "@mui/material";
import Title from "@site/src/components/Title";
import FunctionTable from "./FunctionTable";
import Consultation from "./Consultation";

const FREE_FUNCTION = [
  "智能语义分析检测",
  "反向代理接入",
  "自定义黑白名单",
  "精细化引擎调节",
  "可视化安全分析",
];
const ENTERPRISE_FUNCTION = [
  "智能语义分析检测",
  "串行、旁路均可接入",
  "集群式可扩展部署",
  "CC 攻击防护",
  "业务 API 智能建模",
  "Bot 管理，恶意 Bot 防护",
  "专业技术支持服务",
  "漏洞应急服务",
];

const Version = () => {
  return (
    <>
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          width: "100%",
          flexWrap: "wrap",
        }}
        style={{ marginTop: "48px" }}
      >
        <Box
          sx={{
            width: { xs: "100%", sm: 306 },
            flexShrink: 0,
            height: { xs: "auto", sm: 440 },
            px: 3,
            py: 2,
            mb: { xs: 2, sm: 0 },
            mr: { xs: 0, sm: "20px" },
            border: "1px solid #eee",
            borderRadius: "12px",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            backgroundImage: "linear-gradient(to top,#C3CFE2,#EFF1F8)",
            "&:hover": {
              boxShadow:
                "0 36px 70px -10px rgba(61,64,76,.15), 0 18px 20px -10px rgba(61,64,76,.05)",
            },
            // color: "#fff",
          }}
        >
          <Typography variant="h5" sx={{ fontWeight: 500, fontSize: "18px" }}>
            社区版
          </Typography>
          <Button
            variant="contained"
            // component={}
            target="_blank"
            sx={{
              width: 200,
              my: 2,
            }}
            href="https://stack.chaitin.com/tool/detail?id=717"
          >
            免费使用
          </Button>
          <Box>
            {FREE_FUNCTION.map((f) => (
              <Box
                key={f}
                sx={{
                  py: 1,
                  position: "relative",
                  pl: 2,
                  "&:before": {
                    content: "' '",
                    position: "absolute",
                    left: 0,
                    top: "16px",
                    width: 6,
                    height: 6,
                    borderRadius: "50%",
                    backgroundColor: "success.main",
                  },
                }}
              >
                {f}
              </Box>
            ))}
          </Box>
        </Box>

        <Box
          sx={(theme) => ({
            width: { xs: "100%", sm: 306 },
            flexShrink: 0,
            height: { xs: "auto", sm: 440 },
            px: 3,
            py: 2,
            border: "1px solid",
            borderColor: alpha(theme.palette.primary.main, 0.5),
            borderRadius: "12px",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            "&:hover": {
              boxShadow:
                "0 36px 70px -10px rgba(61,64,76,.15), 0 18px 20px -10px rgba(61,64,76,.05)",
            },
          })}
        >
          <Typography variant="h5" sx={{ fontWeight: 500, fontSize: "18px" }}>
            企业版
          </Typography>
          <Consultation />
          <Box>
            {ENTERPRISE_FUNCTION.map((f) => (
              <Box
                key={f}
                sx={{
                  py: 1,
                  position: "relative",
                  pl: 2,
                  "&:before": {
                    content: "' '",
                    position: "absolute",
                    left: 0,
                    top: "16px",
                    width: 6,
                    height: 6,
                    borderRadius: "50%",
                    backgroundColor: "success.main",
                  },
                }}
              >
                {f}
              </Box>
            ))}
          </Box>
        </Box>
      </Box>
      <Title title="详细参数" sx={{ mt: "120px !important" }} />

      <FunctionTable />
    </>
  );
};

export default Version;
