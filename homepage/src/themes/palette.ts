import { type Color } from './color'
import { type ThemeMode } from './themeContext'
export default function themePalette(color: Color, mode: Omit<ThemeMode, 'system'>) {
  return {
    mode,
    common: { black: '#000', white: '#fff' },
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
  }
}
