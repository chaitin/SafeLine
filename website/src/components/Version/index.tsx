import React from "react";
import { alpha, Box, Button, Typography, Stack } from "@mui/material";
import dynamic from 'next/dynamic';
import FunctionTable from "./FunctionTable";
import Image from "next/image";
import Icon from "@/components/Icon";

const Consultation = dynamic(() => import('./Consultation'), {
  ssr: false,
});

const VERSION_LIST = [
  {
    name: "社区版",
    name_bg: "/images/community-version.svg",
    apply_desc: "适合小型个人网站或业余爱好项目",
    fee: "免费",
    fee_desc: "",
    desc: "",
    operation: (
      <Button
        variant="contained"
        target="_blank"
        sx={{
          width: 146,
          my: 4,
          boxShadow: "0px 15px 25px 0px rgba(15,198,194,0.3)",
        }}
        href="/docs/guide/install"
      >
        立即安装
      </Button>
    ),
    functions: [
      "Web 攻击防护",
      "爬虫防护 / 人机验证",
      "Web 访问控制 / Web 身份认证",
      "CC 攻击防护 / 频率限制",
      "黑 IP 情报防护",
      "0 Day 漏洞情报防护",
    ],
  },
  {
    name: "专业版",
    name_bg: "/images/professional-version.png",
    apply_desc: "适合专业网站与小微企业",
    fee: "¥1799",
    fee_desc: "/年",
    desc: (
      <Stack direction="row" justifyContent="center">
        <Typography variant="subtitle2" mr={1} sx={{ color: alpha("#000", 0.5), textDecoration: "line-through", lineHeight: "20px" }}>
          原价 ¥3600 /年
        </Typography>
        <Image
          src="/images/discount.svg"
          alt="限时特惠"
          width={76}
          height={20}
        />
      </Stack>
    ),
    operation: (
      <Button
        variant="contained"
        // target="_blank"
        sx={{
          width: 146,
          my: 4,
          boxShadow: "0px 15px 25px 0px rgba(15,198,194,0.3)",
        }}
        href="https://rivers.chaitin.cn/?share=85db8d21d63711ee91390242c0a8176b"
      >
        立即购买
      </Button>
    ),
    functions: [
      "所有社区版能力",
      "Web 攻击加强防护",
      "黑 IP 加强情报防护",
      "自定义拦截页面",
      "基于地理位置的访问控制",
      "上游服务器负载均衡",
    ],
  },
  {
    name: "企业版",
    name_bg: "/images/enterprise-version.svg",
    apply_desc: "适合中大型企业",
    fee: (
      <Stack direction="row" alignItems="center">
        <Icon type="icon-zixun" color="#0FC6C2" />
        <Typography variant="h3" sx={{ lineHeight: "46px", ml: 1 }}>按需定制</Typography>
      </Stack>
    ),
    fee_desc: "",
    desc: "",
    operation: (
      <Consultation />
    ),
    functions: [
      "所有专业版能力",
      "软件、硬件、云原生等交付形式",
      "反代、透明、路由、桥接、旁挂等方式部署",
      "多活、主备、Bypass 等高可用形式",
      "动态防护、拟态防护等高级防护形式",
      "分布式集群部署，检测超大规模流量",
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
              width: { xs: "100%", md: "30.7%" },
              flexShrink: 0,
              height: { xs: "auto" },
              borderRadius: "12px",
              mb: { xs: 2, sm: 2, md: 2 },
              mr: index < VERSION_LIST.length - 1 ? { xs: 0, sm: 2, md: 4 } : { xs: 0, sm: 0, md: 0 },
              position: "relative",
              bottom: 0,
              transition: "all 0.5s ease",
              "&:hover": {
                boxShadow: "0px 30px 40px 0px rgba(145,158,171,0.11)",
                transform: "translateY(-16px)"
              },
            }}
          >
            <Box
              sx={{
                px: 3,
                py: 2,
                borderRadius: '12px 12px 0px 0px',
                backgroundSize: "cover",
                backgroundPosition: "center center",
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
              <Box sx={{ height: "20px" }}>
                {item.desc}
              </Box>
              <Box>{item.operation}</Box>
              <FunctionItems items={item.functions} />
            </Box>
          </Box>
        ))}
      </Box>
      <Typography
        variant="h4"
        sx={{
          fontSize: "48px",
          fontWeight: 600,
          mt: { xs: "64px !important", md: "180px !important" },
        }}
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
            fontSize: "14px",
            textAlign: "center",
            // "&:before": {
            //   content: "' '",
            //   position: "absolute",
            //   left: 0,
            //   top: "22px",
            //   width: 4,
            //   height: 4,
            //   borderRadius: "50%",
            //   backgroundColor: "success.main",
            // },
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


