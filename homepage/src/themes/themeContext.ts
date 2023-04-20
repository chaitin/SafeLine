import { createContext } from 'react'
export type ThemeMode = 'dark' | 'light' | 'system'

const ThemeContext = createContext({ mode: 'dark' as ThemeMode, setThemeMode: (mode: ThemeMode) => {} })

export default ThemeContext
