import React from 'react'
import { Box, Button, Container, Stack, Typography } from "@mui/material";
import Image from 'next/image'

const Version = () => {
  return (
    <Box
      sx={{
        width: "100%",
        height: { xs: "205px", md: "343px" },
        backgroundImage: "url(/images/enterprise-bg.svg)",
        backgroundSize: "cover",
        backgroundPosition: "center center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <Container className="relative h-full" sx={{ px: 5 }}>
        <Stack justifyContent="center" className="h-full">
          <Typography
            variant="h4"
            sx={{
              fontWeight: 400,
              color: "common.white",
              fontSize: { xs: "24px", md: "28px" },
              fontFamily: "GilroyBold",
              letterSpacing: "2px",
            }}
          >
            Welcome to other versions of Safeline WAF
          </Typography>
          <Button
            variant="outlined"
            sx={{
              width: { xs: "100%", sm: "384px", md: "146px" },
              height: { xs: "72px", md: "50px" },
              mt: 4,
              backgroundColor: "common.white",
              fontSize: { xs: "24px", md: "16px" },
              fontFamily: 'GilroyBold',
              "&:hover": {
                color: "#0A8A87",
                backgroundColor: "common.white",
              },
            }}
            href="/version"
          >
            Version
          </Button>
        </Stack>
        <Box
          sx={{
            position: "absolute",
            right: { xs: 0, md: -96 },
            top: { xs: -32, md: -65 },
            display: { xs: "none", sm: "block" },
          }}
        >
          <Box width={{ xs: 208, md: 417 }}>
            <Image
              src="/images/shield.png"
              alt="雷池"
              layout="responsive"
              width={417}
              height={359}
            />
          </Box>
        </Box>
      </Container>
    </Box>
  );
};

export default Version;
