import camelcaseKeys from "camelcase-keys";
;
import MissingSection from "./missing-section";
import hero from "@/components/landing-page-sections/hero";
import two_column_with_image from "@/components/landing-page-sections/two-column-with-image";
import features from "@/components/landing-page-sections/features";
import testimonials from "@/components/landing-page-sections/testimonials";

export default function LandingPageSection({ sectionData }) {
  const sectionsComponentPaths = () => ({
    hero,
    two_column_with_image,
    features,
    testimonials,
  });
  const SectionComponent = sectionsComponentPaths();

  return <SectionComponent type={type} {...camelcaseKeys(sectionData)} />;
}
