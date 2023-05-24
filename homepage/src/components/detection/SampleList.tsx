import { useState } from "react";

import {
  Box,
  Button,
  Table,
  TableHead,
  TableRow,
  TableBody,
  Typography,
  Dialog,
  DialogActions,
  DialogContent,
} from "@mui/material";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import TableCell from "@mui/material/TableCell";
import Title from "@/components/Home/Title";
import SampleCount from "@/components/detection/SampleCount";
import SamplesForm from "@/components/detection/SamplesForm";
import hljs from "highlight.js";
import { Message } from "@/components";
import { getSampleDetail } from "@/api/detection";
import { sizeLength } from "@/utils";

import type { RecordSamplesType } from "./types";

export default SampleList;

interface SampleListProps {
  value: RecordSamplesType;
  onSetIdChange: (id: string) => void;
}

function SampleList({ value, onSetIdChange }: SampleListProps) {
  const [open, setOpen] = useState(false);
  const [detail, setDetail] = useState("");

  const handleDetail = (id: string) => async () => {
    const res = await getSampleDetail(id);
    if (res.code != 0) {
      Message.error(res.msg || "获取详情失败");
      return;
    }
    const text = document.createElement("textarea");
    text.innerHTML = res.data.content;
    const highlighted = hljs.highlight(text.value, {
      language: "http",
    });
    setDetail(highlighted.value);
    setOpen(true);
  };

  const handleClose = () => {
    setDetail("");
    setOpen(false);
  };

  return (
    <>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "18px",
          alignItems: "center",
        }}
      >
        <Title title="测试样本" sx={{ fontSize: "16px" }} />
        <SamplesForm onSetIdChange={onSetIdChange} />
      </Box>

      <Accordion
        sx={{ borderRadius: "4px" }}
        className="detection-samples-accordion"
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <SampleCount
            total={value.length}
            normal={value.filter((i) => !i.isAttack).length}
            attack={value.filter((i) => i.isAttack).length}
          />
        </AccordionSummary>
        <AccordionDetails>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell width={150}>样本类型</TableCell>
                <TableCell width={150}>样本大小</TableCell>
                <TableCell>摘要</TableCell>
                <TableCell width={100}></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {value.map((row, index) => {
                const text = document.createElement("textarea");
                text.innerHTML = row.summary;

                return (
                  <TableRow
                    key={index}
                    sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                  >
                    <TableCell>
                      {row.isAttack ? (
                        <Typography sx={{ color: "error.main" }}>
                          攻击样本
                        </Typography>
                      ) : (
                        <Typography sx={{ color: "success.main" }}>
                          普通样本
                        </Typography>
                      )}
                    </TableCell>
                    <TableCell>{sizeLength(row.size)}</TableCell>
                    <TableCell>
                      <Typography
                        noWrap
                        sx={{
                          width: "600px",
                          fontFamily:
                            "ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace",
                        }}
                      >
                        {text.value}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Button onClick={handleDetail(row.id)}>详情</Button>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </AccordionDetails>
      </Accordion>

      <Dialog open={open} onClose={handleClose}>
        <DialogContent sx={{ marginBottom: 0 }}>
          <Box
            component="code"
            style={{ whiteSpace: "pre-line", wordBreak: "break-all" }}
            dangerouslySetInnerHTML={{ __html: detail }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>关闭</Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
