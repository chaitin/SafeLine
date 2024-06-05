import Features from "@/components/landing-page-sections/features";
import Hero from "@/components/landing-page-sections/hero";
import Footer from "@/components/footer-section";
import TwoColumnWithImage from "@/components/landing-page-sections/two-column-with-image";
import { Box } from "@mui/material";
import HeaderSection from "@/components/header-section";

function MyApp() {
  return (
    <Box sx={{ overflow: "auto", height: "100%", width: "100%" }} id='next_container'>
      <HeaderSection />
      <Hero />
      <Features />
      {[
        {
          headline: "Easy To Use",
          subheadline:
            "Deployed by Docker, one command can complete the installation, and you can get started at 0 cost. The security configuration is ready to use, no manual maintenance is required, and safe lying management can be achieved.",
          buttonUrl: "11111",
          image: "EasyToUse.png",
          imagePosition: "right",
          scrollAnchorId: "11111",
        },
        {
          headline: "High Security Efficacy",
          subheadline:
            "The first intelligent semantic analysis algorithm in the industry, accurate detection, low false alarm, and difficult to bypass. The semantic analysis algorithm has no rules, and you are no longer at a loss when facing 0-day attacks with unknown features.",
          buttonUrl: "11111",
          image: "HighSecurityEfficacy.png",
          imagePosition: "left",
          scrollAnchorId: "11111",
        },
        {
          headline: "High Performance",
          subheadline:
            "Ruleless engine, linear security detection algorithm, average request detection delay at 1 millisecond level. Strong concurrency, single core easily detects 2000+ TPS, as long as the hardware is strong enough, there is no upper limit to the traffic scale that can be supported.",
          buttonUrl: "11111",
          image: "HighPerformance.png",
          imagePosition: "right",
          scrollAnchorId: "11111",
        },
        {
          headline: "High Availability",
          subheadline:
            "The traffic processing engine is developed based on Nginx, and both performance and stability can be guaranteed. Built-in complete health check mechanism, service availability is as high as 99.99%.",
          buttonUrl: "11111",
          image: "HighAvailability.png",
          imagePosition: "left",
          scrollAnchorId: "11111",
        },
      ].map((item) => (
        <TwoColumnWithImage key={item.headline} {...item} />
      ))}
      <Footer />
    </Box>
  );
}

export default MyApp;
