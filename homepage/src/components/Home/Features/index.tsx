import { Grid, Box, Typography } from "@mui/material";
import Image from "next/image";

const FEATURE_LIST = [
  {
    title: "便捷",
    content: (
      <>
        <Typography variant="body1" sx={{ mb: "10px" }}>
          采用容器化部署，一条命令即可完成安装，0 成本上手
        </Typography>
        <Typography variant="body1">
          安全配置开箱即用，无需人工维护，可实现安全躺平式管理
        </Typography>
      </>
    ),
  },
  {
    title: "安全",
    content: (
      <>
        <Typography variant="body1" sx={{ mb: "10px" }}>
          首创业内领先的智能语义分析算法，精准检测、低误报、难绕过
        </Typography>
        <Typography variant="body1">
          语义分析算法无规则，面对未知特征的 0day 攻击不再手足无措
        </Typography>
      </>
    ),
  },
  {
    title: "高性能",
    content: (
      <>
        <Typography variant="body1" sx={{ mb: "10px" }}>
          无规则引擎，线性安全检测算法，平均请求检测延迟在 1 毫秒级别
        </Typography>
        <Typography variant="body1">
          并发能力强，单核轻松检测 2000+
          TPS，只要硬件足够强，可支撑的流量规模无上限
        </Typography>
      </>
    ),
  },
  {
    title: "高可用",
    content: (
      <>
        <Typography variant="body1" sx={{ mb: "10px" }}>
          流量处理引擎基于 Nginx 开发，性能与稳定性均可得到保障
        </Typography>
        <Typography variant="body1">
          内置完善的健康检查机制，服务可用性高达 99.99%
        </Typography>
      </>
    ),
  },
];

const Features = () => {
  return (
    <Grid
      container
      xs={12}
      spacing={5}
      sx={{
        borderRadius: "12px",
        backgroundColor: "#fff",
        color: "#000",
        boxShadow: "0 12px 25px -12px rgba(93,99,112, 0.2)",
        pb: 5,
        pr: 5,
        ml: 0,
      }}
    >
      {FEATURE_LIST.map((feature) => (
        <Grid item xs={12} sm={6} key={feature.title}>
          <Typography
            variant="h5"
            sx={{ fontWeight: 500, fontSize: "18px", display: "flex" }}
          >
            <Image
              src="/images/feature.svg"
              alt="feature"
              width={40}
              height={40}
              priority
            />
            <Box component="span" sx={{ mt: "5px" }}>
              {feature.title}
            </Box>
          </Typography>
          <Box sx={{ color: "rgba(0,0,0,.7)", whiteSpace: "pre-line", mt: 3 }}>
            {feature.content}
          </Box>
        </Grid>
      ))}
    </Grid>
  );
};

export default Features;
