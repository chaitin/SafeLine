import { useState, useEffect } from "react";
import { Message, Modal } from "@/components";
import { Box, TextField, Typography, Button } from "@mui/material";

function Consultation() {
  const [text, setText] = useState("");
  const [wrongPhoneNumber, setWrongPhoneNumber] = useState(false);
  const [consultOpen, setConsultOpen] = useState(false);

  const consultHandler = () => {
    const valid = /^1[3-9]\d{9}$/.test(text);
    setWrongPhoneNumber(!valid);
    if (!valid) {
      Message.error("手机号格式不正确");
      return;
    }
    fetch("https://leads.chaitin.net/api/trial", {
      // fetch('http://116.62.230.26:8999/api/trial', { // 测试用地址
      method: "POST",
      mode: "cors",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        phone: text,
        platform_source: "product-official-site",
        product_source: "safeline-ce",
        source_detail: "来自雷池社区版官网",
      }),
    })
      .then((d) => d.json())
      .then((d) => {
        if (d.code == 0) {
          Message.success("提交成功");
        } else {
          Message.error("提交失败");
        }
        setConsultOpen(false);
      });
  };

  const textHandler = (v: string) => {
    setText(v);
  };

  useEffect(() => {
    if (!consultOpen) {
      setText("");
    }
  }, [consultOpen]);

  return (
    <>
      <Button
        variant="outlined"
        sx={{
          width: 200,
          my: 2,
          "&:hover": {
            borderColor: "primary.main",
            backgroundColor: "primary.main",
            color: "#fff",
          },
        }}
        onClick={() => setConsultOpen(true)}
      >
        立即咨询
      </Button>
      <Modal
        open={consultOpen}
        onCancel={() => setConsultOpen(false)}
        title="咨询企业版"
        okText="提交"
        onOk={consultHandler}
        sx={{ width: 600 }}
      >
        <Box sx={{}}>
          <div style={{ margin: "10px 0 15px" }}>
            <TextField
              error={wrongPhoneNumber}
              fullWidth
              size="small"
              label="手机号"
              helperText={
                wrongPhoneNumber ? "手机号格式不正确" : "请输入您的手机号"
              }
              variant="outlined"
              value={text}
              onChange={(e) => textHandler(e.target.value)}
            />
          </div>
          <Typography>
            我们将在工作时间 2 小时内联系您，您的手机号不会用于其他目的
          </Typography>
        </Box>
      </Modal>
    </>
  );
}

export default Consultation;
