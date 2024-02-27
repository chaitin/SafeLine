import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Box,
  Collapse,
  ListItem,
  ListItemIcon,
  ListItemText,
  alpha,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";
import Icon from "@/components/Icon";

const Support = () => {
  return (
    <Icon type="icon-a-yiwanchengtianchong" color="#02BFA5" sx={{ margin: "auto" }} />
  );
};
const NotSupport = () => {
  return <Typography sx={{ margin: "auto" }} >-</Typography>;
};

const Illustrate = ({ text }: { text: string }) => {
  return <Box color="#02BFA5">{text}</Box>;
};

const versions = [
  { title: '社区版', key: 'experience' },
  { title: '专业版', key: 'major' },
  { title: '企业版', key: 'basics' },
]

const colors = ['lighter', 'light', 'main']

const FunctionTable = () => {
  const [open, setOpen] = useState<string[]>([
    "安全防护",
    "统计分析",
    "系统管理",
    "部署形态",
    "服务",
  ]);

  const cells = [
    {
      title: "安全防护",
      data: [
        {
          name: "智能语义分析检测",
          experience: <Support />, // 社区版
          major: <Support />, // 专业版
          basics: <Support />, // 企业版
        },
        {
          name: "简单规则特征库检测",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "频率限制 / CC 攻击防护",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "自定义黑白名单",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "爬虫防护 / 人机验证",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "基础 API 采集",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "威胁情报",
          experience: <Illustrate text="社区共享威胁 IP" />,
          major: <Illustrate text="专业版威胁 IP" />,
          basics: <Illustrate text="高级威胁情报" />,
        },
        {
          name: "商业地理位置封禁",
          experience: <NotSupport />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "应急补充规则",
          experience: <NotSupport />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "高级 Bot 防护、API 防护",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "网页防篡改",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "灵活可编程插件",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "精细化引擎调节",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "统计分析",
      data: [
        {
          name: "基础统计图表",
          tip: "",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "高级统计分析与报告",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "实时可视化大屏",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "系统管理",
      data: [
        {
          name: "自定义拦截页面",
          experience: <NotSupport />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "多设备集中管理",
          tip: "",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "多租户管理",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "安全合规与审计",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "部署形态",
      data: [
        {
          name: "反向代理接入",
          experience: <Support />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "负载均衡",
          experience: <NotSupport />,
          major: <Support />,
          basics: <Support />,
        },
        {
          name: "旁路、透明代理、路由等复杂接入方式",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "集群式可扩展部署",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "紧急模式 / Bypass",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "硬件形态交付",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "高可用 / HA",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "服务",
      data: [
        {
          name: "技术支持服务",
          experience: <Illustrate text="社区技术支持" />,
          major: <Illustrate text="社区优先技术支持" />,
          basics: <Illustrate text="专家点对点技术支持（7*24服务）" />,
        },
        {
          name: "漏洞应急服务",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "等保合规",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "定制化",
          experience: <NotSupport />,
          major: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
  ];

  const handleClick = (id: string) => {
    if (open?.includes(id)) {
      const udpateOpen = [...open];
      udpateOpen.splice(open?.indexOf(id), 1);
      setOpen([...udpateOpen]);
    } else {
      setOpen((open) => [...open, id]);
    }
  };

  return (
    <>
      <TableContainer>
        <Table
          sx={{
            ".MuiTableCell-root": {
              border: "none",
              backgroundColor: "transparent",
              pl: "8px",
              pr: "0px",
            },
          }}
        >
          <TableHead sx={{ background: "transparent" }}>
            <TableRow sx={{ border: "0" }}>
              <TableCell sx={{ width: { xs: "25%" }}}/>
              {versions.map((item, index) => (
                <TableCell key={item.title} align="center" sx={{ width: { xs: "25%" }, fontSize: "16px" }}>
                  <Box
                    sx={(theme) => ({
                      display: "flex",
                      justifyContent: "center",
                      alignItems: "center",
                      width: "100%",
                      height: "40px",
                      borderRadius: "4px",
                      color: theme.palette.common.white,
                      backgroundColor: theme.palette.primary[colors[index] as 'main' | 'light'],
                    })}
                  >
                    {item.title}
                  </Box>
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
        </Table>
      </TableContainer>
      {cells?.map((data) => (
        <React.Fragment key={`sub-table-${data.title}`}>
          <ListItem
            onClick={() => handleClick(data?.title)}
            sx={{
              backgroundColor: "#EFF1F8",
              borderRadius: "4px",
              cursor: "pointer",
              pl: 6,
            }}
          >
            <ListItemText
              primary={data.title}
              sx={{
                ".MuiTypography-root": { fontWeight: 600 },
                flex: "unset",
              }}
            />
            <ListItemIcon sx={{ color: "common.black", marginLeft: 1 }}>
              {open?.includes(data?.title) ? <ExpandLess /> : <ExpandMore />}
            </ListItemIcon>
          </ListItem>
          <Collapse
            in={open?.includes(data?.title)}
            timeout="auto"
            unmountOnExit
            sx={{ width: "100%", mt: "16px !important" }}
          >
            <TableContainer>
              <Table
                sx={{
                  ".MuiTableRow-root": {
                    "&:last-of-type": {
                      ".MuiTableCell-root": {
                        borderBottom: "none",
                      },
                    },
                  },
                  ".MuiTableCell-root": {
                    py: "12px !important",
                    backgroundColor: "transparent !important",
                    color: "#000",
                    borderRight: "1px solid",
                    borderColor: "rgba(0,0,0,.04)",
                    "&:last-of-type": {
                      borderRight: "none",
                    },
                    "&:first-of-type": {
                      paddingLeft: 6,
                    },
                  },
                }}
              >
                <TableBody
                  sx={{
                    backgroundColor: alpha("#EFF1F8", 0.2),
                  }}
                >
                  {data.data.map((item) => (
                    <TableRow key={item.name}>
                      <TableCell sx={{ width: { xs: "25%" }}}>
                        <Box>
                          <Typography variant="h6" sx={{ fontWeight: 400 }}>{item.name}</Typography>
                        </Box>
                      </TableCell>
                      {versions.map((v, index) => (
                        <TableCell key={index} sx={{ width: { xs: "25%" } }} align="center">
                         {item[v.key as 'experience' | 'major' | 'basics']}
                       </TableCell>
                      ))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Collapse>
        </React.Fragment>
      ))}
    </>
  );
};

export default FunctionTable;
