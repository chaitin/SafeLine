import React from "react";
import { Box, Container, Stack, Typography } from "@mui/material";
import Version from "@/components/Version";

export default function VersionView() {

  return (
    <Box mb={26}>
      <Box
        sx={{
          width: "100%",
          height: "380px",
          backgroundImage: "url(/images/version-banner.png)",
          backgroundSize: "cover",
          position: 'relative',
          backgroundPosition: 'center center',
          backgroundRepeat: 'no-repeat'
        }}
      >
        <Container>
          <Box pt={23}>
            <Typography variant="h2" sx={{ fontFamily: "AlimamaShuHeiTi-Bold" }}>大小网站皆宜，免费即可开始</Typography>
          </Box>
        </Container>
      </Box>
      <Container>
        <Stack sx={{ pt: 20 }} spacing={3} alignItems="center">
          <Version />
        </Stack>
      </Container>
    </Box>
  );
}
