import { zhCN } from "@mui/material/locale";
import { createTheme } from "@mui/material/styles";

import * as colors from "./color";
import { type Color } from "./color";
import componentStyleOverrides from "./componentStyleOverrides";
import themePalette from "./palette";
import shadows from "./shadows";
import themeTypography from "./typography";

declare module "@mui/material/styles" {
  interface Palette {
    neutral: Palette["primary"];
  }

  // allow configuration using `createTheme`
  interface PaletteOptions {
    neutral: PaletteOptions["primary"];
  }
}
declare module "@mui/material/Button" {
  interface ButtonPropsColorOverrides {
    neutral: true;
  }
}

export const theme = (mode: "light" | "dark") => {
  const color: Color = colors[mode] as Color;
  const themeOptions = {
    palette: themePalette(color, mode),
    shadows: shadows(color),
    breakpoints: {
      values: {
        xs: 0,
        sm: 680,
        md: 900,
        lg: 1200,
        xl: 1536,
      },
    },
    mixins: {
      toolbar: {
        minHeight: "48px",
        padding: "16px",
        "@media (min-width: 600px)": {
          minHeight: "48px",
        },
      },
    },
    typography: themeTypography(color),
  };

  const themes = createTheme(themeOptions as any, zhCN);
  themes.components = componentStyleOverrides(color);

  return themes;
};

export default theme;
