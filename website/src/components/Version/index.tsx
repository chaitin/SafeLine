import React from "react";
import { alpha, Box, Button, Typography } from "@mui/material";
import FunctionTable from "./FunctionTable";
import Consultation from "./Consultation";

const VERSION_LIST = [
  {
    name: "社区版",
    name_bg: "/images/community-version.png",
    apply_desc: "适合小型个人网站或业余爱好项目",
    fee: "免费",
    fee_desc: "",
    desc: "无限站点，无限规则",
    operation: (
      <Button
        variant="contained"
        target="_blank"
        sx={{
          width: 146,
          my: 4,
          boxShadow: "0px 15px 25px 0px rgba(15,198,194,0.3)",
        }}
        href="https://waf-ce.chaitin.cn/posts/guide_install"
      >
        立即安装
      </Button>
    ),
    functions: [
      "语义分析防护",
      "自定义规则",
      "访问频率限制",
      "人机验证",
      "社区 IP 情报",
    ],
  },
  // {
  //   name: "专业版",
  //   name_bg: "/images/professional-version.png",
  //   apply_desc: "适合专业网站与关键业务",
  //   fee: "¥88",
  //   fee_desc: "/月 按年订阅",
  //   desc: "¥118 按月订阅",
  //   operation: (
  //     <Button
  //       variant="contained"
  //       target="_blank"
  //       sx={{
  //         width: 146,
  //         my: 4,
  //         boxShadow: "0px 15px 25px 0px rgba(15,198,194,0.3)",
  //       }}
  //       href="https://stack.chaitin.com/tool/detail?id=717"
  //     >
  //       立即购买
  //     </Button>
  //   ),
  //   functions: [
  //     "所有社区版能力",
  //     "自定义阻断页面",
  //     "自定义相应动作",
  //     "根据地区封禁",
  //     "......",
  //   ],
  // },
  {
    name: "企业版",
    name_bg: "/images/enterprise-version.png",
    apply_desc: "适合小微到中大型企业",
    fee: "按需定制",
    fee_desc: "",
    desc: "",
    operation: (
      <Consultation />
    ),
    functions: [
      "所有社区版能力",
      "高级统计分析与报告",
      "智能业务建模",
      "高级 Bot 防护、API 防护",
      "......",
    ],
  },
]

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
      >
        {VERSION_LIST.map((item, index) => (
          <Box
            key={item.name}
            sx={{
              width: { xs: "100%", sm: "42%" },
              flexShrink: 0,
              height: { xs: "auto" },
              borderRadius: "12px",
              mb: { xs: 2, sm: 2, md: 2 },
              mr: index < VERSION_LIST.length - 1 ?  { xs: 0, sm: 2, md: 4 } : { xs: 0, sm: 0, md: 0 },
              position: "relative",
              bottom: 0,
              transition: "all 0.5s ease", 
              "&:hover": {
                boxShadow: "0px 30px 40px 0px rgba(145,158,171,0.11)",
                bottom: 16,
              },
            }}
          >
            <Box
              sx={{
                px: 3,
                py: 2,
                borderRadius: '12px 12px 0px 0px',
                backgroundSize: "cover",
                backgroundPosition: "right bottom",
                textAlign: "center",
                backgroundImage: `url(${item.name_bg})`,
              }}
            >
              <Typography variant="h4" sx={{ fontWeight: 500, color: "common.white" }}>
                {item.name}
              </Typography>
            </Box>
            <Box
              sx={{
                py: 3,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                border: "1px solid #E3E8EF",
                borderRadius: '0px 0px 12px 12px ',
                borderTop: "none",
              }}
            >
              <Description content={item.apply_desc} />
              <Typography variant="h3" sx={{ mt: 3, lineHeight: "46px" }}>
                {item.fee}
                {item.fee_desc && (
                  <Typography component="span" variant="subtitle2" sx={{ color: alpha("#000", 0.5), height: "20px" }}>
                    {item.fee_desc}
                  </Typography>
                )}
              </Typography>
              {/* <Description content={item.desc} /> */}
              <Box>{item.operation}</Box>
              <FunctionItems items={item.functions} />
            </Box>
          </Box>
        ))}
      </Box>
      <Typography
        variant="h4"
        sx={{ fontSize: "48px", fontWeight: 600, mt: "180px !important" }}
      >
        细节对比
      </Typography>
      <FunctionTable />
    </>
  );
};

export default Version;


const FunctionItems: React.FC<{ items: any[] }> = ({ items }) => {
  return (
    <Box>
      {items.map((f) => (
        <Box
          key={f}
          sx={{
            py: 2,
            position: "relative",
            pl: 2,
            "&:before": {
              content: "' '",
              position: "absolute",
              left: 0,
              top: "22px",
              width: 4,
              height: 4,
              borderRadius: "50%",
              backgroundColor: "success.main",
            },
          }}
        >
          {f}
        </Box>
      ))}
    </Box>
  );
};

const Description: React.FC<{ content: string }> = ({ content }) => {
  return (
    <Typography variant="subtitle2" sx={{ color: alpha("#000", 0.5), height: "20px" }}>
      {content}
    </Typography>
  );
};


