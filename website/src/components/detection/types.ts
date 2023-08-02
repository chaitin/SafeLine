export type RecordSamplesType = Array<{
  id: string;
  size: number;
  isAttack: boolean;
  summary: string;
}>;

export type SampleDetailType = {
  id: string;
  size: number;
  isAttack: boolean;
  raw: string;
  summary: string;
};

export type ResultRowsType = Array<{
  engine: string;
  version: string;
  detectionRate: number;
  failedRate: number;
  accuracy: number;
  cost: string;
}>;
