import { useMemo, useEffect, useState } from "react";

import {
  ThemeProvider as MUIThemeProvider,
  useMediaQuery,
} from "@mui/material";
import { useLocalStorageState } from "ahooks";

import themes from "../../themes";
import ThemeContext, { type ThemeMode } from "../../themes/themeContext";

interface ThemeProviderProps {
  children?: React.ReactNode;
}

const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  // const [mode, setMode] = useLocalStorageState<ThemeMode>("themeMode", {
  //   defaultValue: "system",
  // });
  const [mode, setMode] = useState<ThemeMode>("light");
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");

  const theme = useMemo(() => {
    let newMode = mode;
    if (mode === "system") {
      newMode = prefersDarkMode ? "dark" : "light";
    }
    return themes(newMode as "dark" | "light");
  }, [mode, prefersDarkMode]);

  const themeMode = useMemo(() => {
    return {
      mode,
      setThemeMode: (mode: ThemeMode) => {
        setMode(mode);
      },
    };
  }, [mode]);

  useEffect(() => {
    const bodyStyle = document.body.style;
    const body = document.body;
    let newMode = mode;
    if (mode === "system") {
      newMode = prefersDarkMode ? "dark" : "light";
    }
    body.className = newMode;
    // @ts-ignore
    bodyStyle.backgroundColor = theme.palette.background.paper0;
    bodyStyle.color = theme.palette.text.primary;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [theme]);

  return (
    <ThemeContext.Provider value={themeMode}>
      <MUIThemeProvider theme={theme}>{children}</MUIThemeProvider>
    </ThemeContext.Provider>
  );
};

export default ThemeProvider;
