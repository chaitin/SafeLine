import React from "react";
import { useState } from "react";
import hljs from "highlight.js";

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
import Paper from "@mui/material/Paper";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import TableCell from "@mui/material/TableCell";
import Title from "@site/src/components/Title";
import SampleCount from "./SampleCount";
import SamplesForm from "./SamplesForm";
import Message from "@site/src/components/Message";
import { getSampleDetail } from "@site/src/api";
import { sizeLength } from "@site/src/components/utils";

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

      <Accordion sx={{ borderRadius: "4px" }}>
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
              <TableRow sx={{ borderBottom: "0" }}>
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
                  <TableRow key={index}>
                    <TableCell sx={{ borderLeft: "0", borderRight: "0" }}>
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
                    <TableCell sx={{ borderLeft: "0", borderRight: "0" }}>
                      {sizeLength(row.size)}
                    </TableCell>
                    <TableCell sx={{ borderLeft: "0", borderRight: "0" }}>
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
                    <TableCell sx={{ borderLeft: "0", borderRight: "0" }}>
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
