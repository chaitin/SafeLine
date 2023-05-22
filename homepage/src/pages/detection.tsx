import { useEffect, useState } from "react";
import Container from "@mui/material/Container";
import Result from "@/components/detection/Result";
import { useRouter } from "next/router";
import "highlight.js/styles/a11y-light.css";
import { getSampleSet, getSampleSetResult } from "@/api/detection";
import { Message } from "@/components";
import type {
  RecordSamplesType,
  ResultRowsType,
} from "@/components/detection/types";
import SampleList from "@/components/detection/SampleList";

export default Detection;

function Detection() {
  const router = useRouter();
  const [samples, setSamples] = useState<RecordSamplesType>([]);
  const [result, setResult] = useState<ResultRowsType>([]);

  useEffect(() => {
    // useRouter 中获取 参数会有延迟，所以先判断有没有 id 参数
    const realSetId =
      new URLSearchParams(location.search).get("id") || "default";
    const setId = (router.query.id as string) || "default";
    if (setId !== realSetId) return;

    // 查询样本集合
    getSampleSet(setId).then((res) => {
      if (res.code != 0) {
        Message.error("测试集合 " + setId + ": " + res.msg);
        return;
      }
      if (!res.data.data) {
        Message.error("测试集合 " + setId + ": 获取结果为空");
        return;
      }
      setSamples(
        res.data.data?.map((i: any) => ({
          id: i.id,
          summary: i.summary,
          size: i.length,
          isAttack: i.tag == "black",
        }))
      );
    });

    // 查询样本集合结果
    getSampleSetResult(setId).then(({ data, timeout }) => {
      if (timeout) {
        Message.error("获取检测集结果超时");
        return;
      }
      setResult(
        data.map((i: any) => ({
          engine: i.engine,
          version: i.version,
          detectionRate: percent(i.recall),
          failedRate: percent(i.fdr),
          accuracy: percent(i.accuracy),
          cost: i.elapsed + "毫秒",
        }))
      );
    });
  }, [router.query.id]);

  const handleSetId = (id: string) => {
    router.push({
      pathname: router.pathname,
      query: { id },
    });
  };

  return (
    <Container sx={{ mb: 2 }}>
      <div style={{ height: "100px" }}></div>
      <SampleList value={samples} onSetIdChange={handleSetId} />
      <Result rows={result} />
    </Container>
  );
}

function percent(v: number) {
  return Math.round(v * 10000) / 100 + "%";
}
