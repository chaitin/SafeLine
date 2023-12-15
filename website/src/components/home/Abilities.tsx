import React, { useState } from "react";
import { Grid, Box, Typography, Container, Link } from "@mui/material";
import Image from "next/image";
import Icon from "@/components/Icon";

const ABILITY_LIST = [
  {
    title: "人机验证",
    href: "/docs/about/challenge",
    img: "/images/ability/ability_verification.png",
  },
  {
    title: "百川网站监控联动",
    href: "/docs/practice/monitor",
    img: "/images/ability/ability_rivers.png",
  },
  {
    title: "APISIX 插件集成",
    href: "/docs/about/apisix",
    img: "/images/ability/ability_apisix.svg",
  },
  {
    title: "长亭社区恶意 IP 情报",
    href: "/docs/about/IpIntelligence",
    img: "/images/ability/ability_maliciousip.svg",
  },
  {
    title: "申请免费证书",
    href: "",
    img: "/images/ability/ability_cert.svg",
  },
  {
    title: "站点资源一览",
    href: "",
    img: "/images/ability/ability_asset.png",
  },
  {
    title: "CC 攻击防护",
    href: "",
    img: "/images/ability/ability_CC.svg",
  },
  {
    title: "一键强制 HTTPS",
    href: "",
    img: "/images/ability/ability_HTTPS.svg",
  },
];

const DEFAULT_URL = "/images/ability/ability_verification.png";

const Abilities = () => {
  const [hoveredUrl, setHoveredUrl] = useState(DEFAULT_URL);

  const handleIconHover = (url: string) => {
    setHoveredUrl(url);
  };

  return (
    <Box
      position={"relative"}
      sx={{
        background: "#111227",
        color: "common.white",
        pt: { xs: 8, md: 18 },
        pb: { xs: 8, md: 27 },
        px: 2,
      }}
    >
      <Container maxWidth="lg">
        <Grid container alignItems="center">
          <Grid item xs={12} md={6}>
            <Typography variant="h2" mb={4.5} textAlign={{ xs: "center", md: "left" }} fontSize={{ xs: "32px", md: "48px" }}>
              多维能力拓展
            </Typography>
            <Grid container spacing={2}>
              {ABILITY_LIST.map((ability) => (
                <AbilityItem
                  key={ability.title}
                  ability={ability}
                  hoveredUrl={hoveredUrl}
                  handleIconHover={handleIconHover}
                />
              ))}
            </Grid>
          </Grid>
          <Grid item xs={0} md={6} display={{ xs: "none", md: "block" }}>
            <Box sx={{ width: { xs: "100%", sm: "100%" } }}>
              {ABILITY_LIST.map((ability) => (
                <Image
                  key={ability.title}
                  src={ability.img}
                  alt={ability.title}
                  layout={hoveredUrl === ability.img ? "responsive" : "fixed"}
                  width={0}
                  height={100}
                />
              ))}
            </Box>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

export default Abilities;

interface ItemProps {
  ability: any
  hoveredUrl?: string;
  handleIconHover: Function;
}

const AbilityItem: React.FC<ItemProps> = ({
  ability,
  hoveredUrl,
  handleIconHover,
}) => {
  return (
    <Grid item xs={12} sm={6}>
      <Box
        sx={{
          background: "#1B243E",
          boxShadow: "0px 4px 10px 0px rgba(3,13,23,0.6)",
          borderRadius: "12px",
          width: { xs: "100%", lg: "274px" },
        }}
        onMouseEnter={() => handleIconHover(ability.img)}
        onMouseLeave={() => {}}
        onClick={() => handleIconHover(ability.img)}
      >
        {ability.href ? (
          <Link href={ability.href} target="_blank" rel={ability.title}>
            <Typography
              variant="h6"
              px={3}
              py={3}
              sx={{
                fontSize: "20px",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                color: "common.white",
                "&:hover": {
                  backgroundColor: "primary.main",
                  boxShadow: "0px 4px 10px 0px rgba(3,13,23,0.6)",
                  borderRadius: "12px",
                  color: "common.white",
                },
              }}
            >
              {ability.title}
              <Icon type="icon-youjiantouxian" />
            </Typography>
          </Link>
        ) : (
          <Typography
            variant="h6"
            px={3}
            py={3}
            sx={{
              fontSize: "20px",
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            {ability.title}
          </Typography>
        )}
      </Box>
    </Grid>
  );
};
