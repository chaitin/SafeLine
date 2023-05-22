import Typography from "@mui/material/Typography";

export default SampleCount;

function SampleCount({
  total,
  normal,
  attack,
}: {
  total: number;
  normal: number;
  attack: number;
}) {
  return (
    <div style={{ display: "flex" }}>
      <Typography>共计 {total} 个 HTTP 请求样本，其中 </Typography>
      <Typography sx={{ color: "success.main" }}>
        普通样本 {normal} 个
      </Typography>
      、
      <Typography sx={{ color: "error.main" }}>攻击样本 {attack} 个</Typography>
    </div>
  );
}
