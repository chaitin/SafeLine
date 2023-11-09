import React, { useEffect, useMemo } from "react";
import { zhCN } from "@mui/material/locale";
import { createTheme } from "@mui/material/styles";

import { ThemeProvider as MUIThemeProvider } from "@mui/material";

interface Props {
  children?: React.ReactNode;
}

export default ThemeProvider;

function ThemeProvider({ children }: Props) {
  const theme = useMemo(() => {
    return themes();
  }, []);

  useEffect(() => {
    const bodyStyle = document.body.style;
    // @ts-ignore
    bodyStyle.backgroundColor = theme.palette.background.paper0;
    bodyStyle.color = theme.palette.text.primary;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [theme]);

  return <MUIThemeProvider theme={theme}>{children}</MUIThemeProvider>;
}

function themes() {
  //   const color: Color = colors[mode] as Color;
  const themeOptions = {
    palette: themePalette(light),
    shadows: shadows(light),
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
    typography: {
      fontFamily: "inherit",
      body1: {
        fontSize: "16px",
      },
      body2: {
        fontSize: "14px",
      },
      subtitle1: {
        fontSize: "24px",
      },
      subtitle2: {
        fontSize: "12px",
        fontWeight: 400
      },
      h6: {
        fontSize: "16px",
        fontWeight: 600, 
      },
      h5: {
        fontSize: "24px",
        fontWeight: 600,
        lineHeight: "32px",
      },
      h4: {
        fontSize: "28px",
        fontWeight: 600,
      },
      h3: {
        fontSize: "38px",
        fontWeight: 600,
      },
      h2: {
        fontSize: "48px",
        fontWeight: 600,
      },
      h1: {
        fontSize: "80px",
        fontWeight: "bold",
      },
    },
  };

  const themes = createTheme(themeOptions as any, zhCN);
  themes.components = componentStyleOverrides(light);

  return themes;
}

function themePalette(color: Color) {
  return {
    mode: "light",
    common: { black: "#000", white: "#fff" },
    primary: color.primary,
    secondary: color.secondary,
    info: color.info,
    success: color.success,
    warning: color.warning,
    error: color.error,
    neutral: color.neutral,
    divider: color.divider,
    text: color.text,
    background: color.background,
    shadowColor: color.shadowColor,
    charts: color.charts,
    action: {
      selectedOpacity: 0.1,
    },
  };
}

function shadows(color: Color) {
  return [
    `0px 12px 24px -4px ${color.shadowColor},0px 0px 2px 0px ${color.shadowColor}`,
    `0px 12px 24px -4px ${color.shadowColor},0px 0px 2px 0px ${color.shadowColor}`,
    `0px 3px 1px -2px rgba(0,0,0,0.2),0px 2px 2px 0px rgba(0,0,0,0.12),0px 1px 5px 0px rgba(0,0,0,0.12)`,
    `0px 3px 3px -2px rgba(0,0,0,0.2),0px 3px 4px 0px rgba(0,0,0,0.12),0px 1px 8px 0px rgba(0,0,0,0.12)`,
    `0px 2px 4px -1px rgba(0,0,0,0.2),0px 4px 5px 0px rgba(0,0,0,0.12),0px 1px 10px 0px rgba(0,0,0,0.12)`,
    `0px 3px 5px -1px rgba(0,0,0,0.2),0px 5px 8px 0px rgba(0,0,0,0.12),0px 1px 14px 0px rgba(0,0,0,0.12)`,
    `0px 3px 5px -1px rgba(0,0,0,0.2),0px 6px 10px 0px rgba(0,0,0,0.12),0px 1px 18px 0px rgba(0,0,0,0.12)`,
    `0px 4px 5px -2px rgba(0,0,0,0.2),0px 7px 10px 1px rgba(0,0,0,0.12),0px 2px 16px 1px rgba(0,0,0,0.12)`,
    `0px 5px 5px -3px rgba(0,0,0,0.2),0px 8px 10px 1px rgba(0,0,0,0.12),0px 3px 14px 2px rgba(0,0,0,0.12)`,
    `0px 5px 6px -3px rgba(0,0,0,0.2),0px 9px 12px 1px rgba(0,0,0,0.12),0px 3px 16px 2px rgba(0,0,0,0.12)`,
    `0px 6px 6px -3px rgba(0,0,0,0.2),0px 10px 14px 1px rgba(0,0,0,0.12),0px 4px 18px 3px rgba(0,0,0,0.12)`,
    `0px 6px 7px -4px rgba(0,0,0,0.2),0px 11px 15px 1px rgba(0,0,0,0.12),0px 4px 20px 3px rgba(0,0,0,0.12)`,
    `0px 7px 8px -4px rgba(0,0,0,0.2),0px 12px 17px 2px rgba(0,0,0,0.12),0px 5px 22px 4px rgba(0,0,0,0.12)`,
    `0px 7px 8px -4px rgba(0,0,0,0.2),0px 13px 19px 2px rgba(0,0,0,0.12),0px 5px 24px 4px rgba(0,0,0,0.12)`,
    `0px 7px 9px -4px rgba(0,0,0,0.2),0px 14px 21px 2px rgba(0,0,0,0.12),0px 5px 26px 4px rgba(0,0,0,0.12)`,
    `0px 8px 9px -5px rgba(0,0,0,0.2),0px 15px 22px 2px rgba(0,0,0,0.12),0px 6px 28px 5px rgba(0,0,0,0.12)`,
    `0px 8px 10px -5px rgba(0,0,0,0.2),0px 16px 24px 2px rgba(0,0,0,0.12),0px 6px 30px 5px rgba(0,0,0,0.12)`,
    `0px 8px 11px -5px rgba(0,0,0,0.2),0px 17px 26px 2px rgba(0,0,0,0.12),0px 6px 32px 5px rgba(0,0,0,0.12)`,
    `0px 9px 11px -5px rgba(0,0,0,0.2),0px 18px 28px 2px rgba(0,0,0,0.12),0px 7px 34px 6px rgba(0,0,0,0.12)`,
    `0px 9px 12px -6px rgba(0,0,0,0.2),0px 19px 29px 2px rgba(0,0,0,0.12),0px 7px 36px 6px rgba(0,0,0,0.12)`,
    `0px 10px 13px -6px rgba(0,0,0,0.2),0px 20px 31px 3px rgba(0,0,0,0.12),0px 8px 38px 7px rgba(0,0,0,0.12)`,
    `0px 10px 13px -6px rgba(0,0,0,0.2),0px 21px 33px 3px rgba(0,0,0,0.12),0px 8px 40px 7px rgba(0,0,0,0.12)`,
    `0px 10px 14px -6px rgba(0,0,0,0.2),0px 22px 35px 3px rgba(0,0,0,0.12),0px 8px 42px 7px rgba(0,0,0,0.12)`,
    `0px 11px 14px -7px rgba(0,0,0,0.2),0px 23px 36px 3px rgba(0,0,0,0.12),0px 9px 44px 8px rgba(0,0,0,0.12)`,
    `0px 11px 15px -7px rgba(0,0,0,0.2),0px 24px 38px 3px rgba(0,0,0,0.12),0px 9px 46px 8px rgba(0,0,0,0.12)`,
  ];
}

const light = {
  primary: {
    main: "#0FC6C2",
    lighter: "#B7EDEC",
    light: "#57D7D4",
    contrastText: "#fff",
  },
  secondary: {
    lighter: "#D6E4FF",
    light: "#84A9FF",
    main: "#3366FF",
    dark: "#1939B7",
    darker: "#091A7A",
    contrastText: "#fff",
  },
  info: {
    lighter: "#D0F2FF",
    light: "#74CAFF",
    main: "#1890FF",
    dark: "#0C53B7",
    darker: "#04297A",
    contrastText: "#fff",
  },
  success: {
    lighter: "#E9FCD4",
    light: "#AAF27F",
    main: "#0FC6C2",
    mainShadow: "#02BFA5",
    dark: "#229A16",
    darker: "#08660D",
    contrastText: "rgba(0,0,0,0.7)",
  },
  warning: {
    lighter: "#FFF7CD",
    light: "#FFE16A",
    main: "#FFBF00",
    dark: "#B78103",
    darker: "#7A4F01",
    contrastText: "rgba(0,0,0,0.7)",
  },
  neutral: {
    main: "#F4F6F8",
    contrastText: "rgba(0, 0, 0, 0.60)",
  },
  error: {
    lighter: "#FFE7D9",
    light: "#FFA48D",
    main: "#FF1844",
    dark: "#B72136",
    darker: "#7A0C2E",
    contrastText: "#fff",
  },
  divider: "#E3E8EF",
  text: {
    primary: "#000",
    secondary: "rgba(0,0,0,0.7)",
    auxiliary: "rgba(0,0,0,0.5)",
    slave: "rgba(0,0,0,0.05)",
    disabled: "rgba(0,0,0,0.15)",
    inversePrimary: "#fff",
    inverseAuxiliary: "rgba(255,255,255,0.5)",
    inverseDisabled: "rgba(255,255,255,0.15)",
  },
  background: {
    paper0: "#fff",
    paper: "#fff",
    paper2: "#f6f8fa",
    default: "#fff",
    chip: "#F4F6F8",
    circle: "#E6E8EC",
  },

  shadowColor: "rgba(145,158,171,0.2)",
  table: {
    head: {
      backgroundColor: "#F4F6F8",
      color: "rgba(0,0,0,0.7)",
    },
    row: {
      hoverColor: "#F9FAFB",
    },
    cell: {
      borderColor: "#F3F4F5",
    },
  },
  charts: {
    color: ["#673AB7", "#02BFA5"],
  },
};

type Color = typeof light;

function componentStyleOverrides(color: Color) {
  return {
    MuiButton: {
      styleOverrides: {
        root: ({ ownerState, theme }: any) => {
          return {
            boxShadow: "none",
            "&:hover": {
              boxShadow: "none",
              ...(ownerState.color === "neutral" && {
                color: color.primary.main,
                fontWeight: 700,
              }),
            },
          };
        },
      },
    },
    MuiFormControl: {
      styleOverrides: {
        root: {
          ".MuiFormLabel-asterisk": {
            color: color.error.main,
          },
        },
      },
    },
    MuiTableRow: {
      styleOverrides: {
        root: {
          "&:hover": {
            ".MuiTableCell-root": {
              backgroundColor: color.table.row.hoverColor,
            },
          },
        },
      },
    },
    MuiTableBody: {
      styleOverrides: {
        root: {
          ".MuiTableRow-root:hover": {
            ".MuiTableCell-root": {
              backgroundColor: color.table.row.hoverColor,
            },
          },
        },
      },
    },

    MuiTableCell: {
      styleOverrides: {
        root: {
          background: color.background.paper,
          lineHeight: 1.5,
          fontSize: "14px",
          paddingTop: "24px",
          paddingBottom: "24px",
          borderColor: color.table.cell.borderColor,
          paddingLeft: 0,
          "&:first-of-type": {
            paddingLeft: "32px",
          },
        },
        head: {
          backgroundColor: color.table.head.backgroundColor,
          color: color.table.head.color,
          fontSize: "12px",
          height: "24px",
          paddingTop: 0,
          paddingBottom: 0,
        },
      },
    },
    MuiMenu: {
      defaultProps: {
        elevation: 0,
      },
    },
    MuiPaper: {
      defaultProps: {
        elevation: 1,
      },
      styleOverrides: {
        root: ({ ownerState }: any) => {
          return {
            ...(ownerState.elevation === 0 && {
              backgroundColor: color.background.paper0,
            }),
            ...(ownerState.elevation === 2 && {
              backgroundColor: color.background.paper2,
            }),
            backgroundImage: "none",
          };
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          // borderRadius: "4px",
          height: "auto",
        },
      },
    },
    MuiAppBar: {
      defaultProps: {
        elevation: 1,
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          "&:hover .MuiOutlinedInput-notchedOutline": {
            borderColor: color.primary.main,
          },
          "& .MuiOutlinedInput-notchedOutline": {
            borderColor: color.divider,
          },
        },
      },
    },
    MuiLink: {
      styleOverrides: {
        root: {
          color: color.primary.contrastText,
          '&:hover': {
            color: color.primary.main,
          },
        },
      },
    },
  };
}
