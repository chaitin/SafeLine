import { useState } from "react";
import {
  Box,
  Button,
  RadioGroup,
  StepLabel,
  TextField,
  FormControlLabel,
  Radio,
} from "@mui/material";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@mui/material";

import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import SampleCount from "@/components/detection/SampleCount";

import {
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
} from "@mui/material";
import { sampleLength, sampleSummary } from "@/utils";

export default SampleSteps;

interface SampleStepsProps {
  // onDetect: (samples: Array<{ sample: string; isAttack: boolean }>) => void;
  onDetect: (value: { sample: string; isAttack: boolean }) => void;
}

function SampleSteps({ onDetect }: SampleStepsProps) {
  const [activeStep, setActiveStep] = useState(0);
  const [completed, setCompleted] = useState([false, false, false]);
  const [sampleText, setSampleText] = useState("");
  const [sampleIsAttack, setSampleIsAttack] = useState(false);
  const [sampleTextError, setSampleTextError] = useState("");
  const [count, setCount] = useState({
    total: 0,
    normal: 0,
    attack: 0,
  });

  const nextStep = () => {
    setActiveStep(activeStep + 1);
    completed[activeStep] = true;
    setCompleted([...completed]);
  };

  const handleCommit = () => {
    if (sampleText == "") {
      setSampleTextError("样本内容不能为空");
      return;
    }
    nextStep();
  };

  const handleMark = (v: Array<{ isAttack: boolean }>) => {
    const isAttack = v[0].isAttack;
    setSampleIsAttack(isAttack);
    setCount({ total: 1, normal: isAttack ? 0 : 1, attack: isAttack ? 1 : 0 });
    nextStep();
  };

  const handleDetect = () => {
    onDetect({
      sample: sampleText,
      isAttack: sampleIsAttack,
    });
  };

  const handleStep = (n: number) => {
    if(completed[n]) setActiveStep(n)
  }

  return (
    <Box>
      <Stepper nonLinear activeStep={activeStep} sx={{ marginBottom: 2 }}>
        <Step completed={completed[0]} onClick={() => handleStep(0)}>
          <StepLabel>自定义样本</StepLabel>
        </Step>
        <Step completed={completed[1]} onClick={() => handleStep(1)}>
          <StepLabel>标记样本</StepLabel>
        </Step>
        <Step completed={completed[2]} onClick={() => handleStep(2)}>
          <StepLabel>开始测试</StepLabel>
        </Step>
      </Stepper>
      {activeStep == 0 && (
        <Box sx={{ p: 1 }}>
          <TextField
            multiline
            fullWidth
            minRows={4}
            error={sampleTextError != ""}
            helperText={sampleTextError}
            value={sampleText}
            onChange={(e) => setSampleText(e.target.value)}
            placeholder={`GET /path/api HTTP/1.1
Host: example.com`}
            sx={{ mb: 2 }}
          />
          <Button fullWidth variant="contained" onClick={handleCommit}>
            提交样本
          </Button>
        </Box>
      )}
      {activeStep == 1 && (
        <SampleMarkable onChange={handleMark} samples={[sampleText]} />
      )}
      {activeStep == 2 && (
        <Box>
          <SampleCount
            total={count.total}
            normal={count.normal}
            attack={count.attack}
          />
          <Button
            sx={{ mt: 2 }}
            fullWidth
            variant="contained"
            onClick={handleDetect}
          >
            开始检测
          </Button>
        </Box>
      )}
    </Box>
  );
}

function SampleMarkable({
  samples,
  onChange,
}: {
  samples: Array<string>;
  onChange: (value: Array<{ sample: string; isAttack: boolean }>) => void;
}) {
  const [sampleDetail, setSampleDetail] = useState("");
  const [rows, setRows] = useState(() => {
    const rows: Array<{
      isAttack: boolean;
      summary: string;
      size: string;
      raw: string;
    }> = [];

    samples.forEach((i) => {
      rows.push({
        isAttack: false,
        summary: sampleSummary(i),
        raw: i,
        size: sampleLength(i),
      });
    });

    return rows;
  });

  const handle = () => {
    onChange(rows.map((i) => ({ sample: i.raw, isAttack: i.isAttack })));
  };

  const handleType = (n: number) => (v: string) => {
    rows[n].isAttack = v == "attack";
    setRows([...rows]);
  };

  const handleSampleDetail = (text: string) => () => {
    setSampleDetail(text);
  };

  const handleClose = () => {
    setSampleDetail("");
  };

  return (
    <Box>
      <Table sx={{ mb: 2 }}>
        <TableHead>
          <TableRow>
            <TableCell width={350}>样本类型</TableCell>
            <TableCell width={150}>样本大小</TableCell>
            <TableCell>摘要</TableCell>
            <TableCell width={100}></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row, index) => (
            <TableRow
              key={index}
              sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
            >
              <TableCell>
                <RadioGroup
                  value={row.isAttack ? "attack" : "normal"}
                  onChange={(e) => handleType(index)(e.target.value)}
                  sx={{ flexDirection: "row" }}
                >
                  <FormControlLabel
                    value="attack"
                    control={<Radio sx={{ color: "divider" }} />}
                    label="攻击样本"
                  />
                  <FormControlLabel
                    value="normal"
                    control={<Radio sx={{ color: "divider" }} />}
                    label="普通样本"
                  />
                </RadioGroup>
              </TableCell>
              <TableCell>{row.size}</TableCell>
              <TableCell>{row.summary}</TableCell>
              <TableCell>
                <Button onClick={handleSampleDetail(row.raw)}>详情</Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Button onClick={handle} fullWidth variant="contained">
        标记样本
      </Button>

      <Dialog open={sampleDetail != ""} onClick={handleClose}>
        <DialogTitle>样本详情</DialogTitle>
        <DialogContent>
          <DialogContentText>{sampleDetail}</DialogContentText>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose}>关闭</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
