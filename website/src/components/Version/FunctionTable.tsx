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
  Tooltip,
} from "@mui/material";
import React, { useState } from "react";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";
import Icon from "@site/src/components/Icon";

const Support = () => {
  return (
    <Icon type="icon-duihao" color="success.main" sx={{ margin: "auto" }} />
  );
};
const NotSupport = () => {
  return <Icon type="icon-chahao" color="error.main" sx={{ margin: "auto" }} />;
};

const Unlimited = () => {
  return <Box color="success.main">无限制</Box>;
};

const FunctionTable = () => {
  const [open, setOpen] = useState<string[]>([
    "基本功能",
    "高级防护能力",
    "部署形态",
    "流量接入方式",
  ]);

  const cells = [
    {
      title: "基本功能",
      data: [
        {
          name: "可视化 DashBoard",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "自定义黑白名单",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "自定义防护策略",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "可防护站点数量",
          experience: <Unlimited />,
          basics: <Unlimited />,
        },
        {
          name: "可支撑流量大小",
          experience: <Unlimited />,
          basics: <Unlimited />,
        },
      ],
    },
    {
      title: "高级防护能力",
      data: [
        {
          name: "智能语义分析引擎",
          tip: "",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "智能业务建模",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "动态防护",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "API 防护",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "Bot 管理",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "拟态防护",
          experience: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "部署形态",
      data: [
        {
          name: "软件形态",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "硬件形态",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "分布式集群形态",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "云原生 K8S 集群形态",
          experience: <NotSupport />,
          basics: <Support />,
        },
      ],
    },
    {
      title: "流量接入方式",
      data: [
        {
          name: "反向代理接入",
          experience: <Support />,
          basics: <Support />,
        },
        {
          name: "SDK 接入",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "透明代理接入",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "透明桥接接入",
          experience: <NotSupport />,
          basics: <Support />,
        },
        {
          name: "旁路镜像接入",
          experience: <NotSupport />,
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
              px: "12px",
            },
          }}
          style={{ marginTop: "48px" }}
        >
          <TableHead>
            <TableRow>
              <TableCell sx={{ width: "33%" }} />
              <TableCell align="center" sx={{ width: "33%", fontSize: "16px" }}>
                <Box
                  sx={(theme) => ({
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    width: "100%",
                    height: "40px",
                    borderRadius: "4px",
                    color: theme.palette.primary.main,
                    backgroundColor: alpha(theme.palette.primary.main, 0.2),
                  })}
                >
                  社区版
                </Box>
              </TableCell>
              <TableCell align="left" sx={{ width: "33%", fontSize: "16px" }}>
                <Box
                  sx={(theme) => ({
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    width: "100%",
                    color: theme.palette.primary.main,
                    backgroundColor: alpha(theme.palette.primary.main, 0.1),
                    height: "40px",
                    borderRadius: "4px",
                  })}
                >
                  企业版
                </Box>
              </TableCell>
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
              pl: "20px",
            }}
          >
            <ListItemText
              primary={data.title}
              sx={{ ".MuiTypography-root": { fontSize: "16px" } }}
            />
            <ListItemIcon sx={{ color: "#000" }}>
              {open?.includes(data?.title) ? <ExpandLess /> : <ExpandMore />}
            </ListItemIcon>
          </ListItem>
          <Collapse
            in={open?.includes(data?.title)}
            timeout="auto"
            unmountOnExit
            sx={{ width: "100%" }}
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
                    pl: "20px !important",
                    pr: "8px !important",
                    py: "12px !important",
                    backgroundColor: "transparent !important",
                    color: "#000",
                    borderRight: "1px solid",
                    borderColor: "rgba(0,0,0,.04)",
                    "&:last-of-type": {
                      borderRight: "none",
                    },
                  },
                }}
              >
                <TableBody
                  sx={{
                    backgroundColor: alpha("#EFF1F8", 0.6),
                  }}
                >
                  {data.data.map((item) => (
                    <TableRow key={item.name}>
                      <TableCell sx={{ width: "33%" }}>
                        <Box sx={{ display: "flex", alignItems: "center" }}>
                          {item.name}
                          <Tooltip title={item.tip} arrow placement="right">
                            <Box component="span" sx={{ ml: 1 }}>
                              <Icon type="icon-bangzhu1" />
                            </Box>
                          </Tooltip>
                        </Box>
                      </TableCell>
                      <TableCell sx={{ width: "33%" }} align="center">
                        {item.experience}
                      </TableCell>
                      <TableCell sx={{ width: "33%" }} align="center">
                        {item.basics}
                      </TableCell>
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
