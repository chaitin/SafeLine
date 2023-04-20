import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Box,
  Button,
} from "@mui/material";
import Link from "next/link";
const Support = () => {
  return <Box color="success.main">支持</Box>;
};
const NotSupport = () => {
  return <Box color="error.main">不支持</Box>;
};

const SoonSupport = () => {
  return <Box color="warning.main">即将支持</Box>;
};

const Unlimited = () => {
  return <Box color="success.main">无限制</Box>;
};

const Version = () => {
  const bodyCell = [
    {
      name: "智能语义分析引擎",
      experience: <Support />,
      basics: <Support />,
    },
    {
      name: "可防护站点数量",
      experience: <Unlimited />,
      basics: <Unlimited />,
    },
    {
      name: "自定义白名单",
      experience: <Support />,
      basics: <Support />,
    },
    {
      name: "自定义黑名单",
      experience: <Support />,
      basics: <Support />,
    },
    {
      name: "反向代理接入",
      experience: <Support />,
      basics: <Support />,
    },
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
      name: "分布式集群部署",
      experience: <NotSupport />,
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
    {
      name: "",
      experience: (
        <Button
          variant="contained"
          component={Link}
          href="https://stack.chaitin.com/tool/detail?id=717"
          target="_blank"
        >
          免费使用
        </Button>
      ),
      basics: (
        <Button
          variant="contained"
          color="neutral"
          component={Link}
          href="https://stack.chaitin.com/tool/detail?id=717"
          target="_blank"
          sx={{ backgroundColor: "#fff", color: "#000" }}
        >
          立即咨询
        </Button>
      ),
    },
  ];
  return (
    <TableContainer>
      <Table
        sx={{
          ".MuiTableCell-root": {
            pl: "8px !important",
            pr: "8px !important",
            backgroundColor: "background.paper",
            borderRight: "1px solid",
            borderColor: "divider",
            "&:first-of-type": {
              borderLeft: "1px solid",
              borderColor: "divider",
            },
          },
          ".MuiTableHead-root": {
            ".MuiTableCell-root": {
              height: "60px",
              fontSize: "16px",
              fontWeight: 700,
              color: "text.primary",
              borderTop: "1px solid",
              borderColor: "divider",
            },
          },
        }}
      >
        <TableHead>
          <TableRow>
            <TableCell />
            <TableCell align="left">
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  width: "100%",
                }}
              >
                社区版{" "}
                <Box
                  sx={{
                    display: { xs: "none", sm: "block" },
                    color: "primary.main",
                  }}
                >
                  免费
                </Box>
              </Box>
            </TableCell>
            <TableCell align="left">
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  width: "100%",
                }}
              >
                企业版
                <Box
                  sx={{
                    display: { xs: "none", sm: "block" },
                  }}
                >
                  独立报价
                </Box>
              </Box>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {bodyCell.map((item) => (
            <TableRow key={item.name}>
              <TableCell>{item.name}</TableCell>
              <TableCell align="left">{item.experience}</TableCell>
              <TableCell align="left">{item.basics}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default Version;
