export const dark = {
  primary: {
    main: '#7267EF',
    contrastText: '#fff',
  },
  secondary: {
    lighter: '#D6E4FF',
    light: '#84A9FF',
    main: '#2196F3',
    dark: '#1939B7',
    darker: '#091A7A',
    contrastText: '#fff',
  },
  info: {
    lighter: '#D0F2FF',
    light: '#74CAFF',
    main: '#1890FF',
    dark: '#0C53B7',
    darker: '#04297A',
    contrastText: '#fff',
  },

  success: {
    lighter: '#E9FCD4',
    light: '#AAF27F',
    main: '#02BFA5',
    dark: '#229A16',
    darker: '#08660D',
    contrastText: 'rgba(0,0,0,0.7)',
  },

  warning: {
    lighter: '#FFF7CD',
    light: '#FFE16A',
    main: '#FFBF00',
    dark: '#B78103',
    darker: '#7A4F01',
    contrastText: 'rgba(0,0,0,0.7)',
  },
  neutral: {
    main: '#232C59',
    contrastText: 'rgba(255, 255, 255, 0.60)',
  },
  error: {
    lighter: '#FFE7D9',
    light: '#FFA48D',
    main: '#FF1844',
    dark: '#B72136',
    darker: '#7A0C2E',
    contrastText: '#fff',
  },
  text: {
    primary: '#fff',
    secondary: 'rgba(255,255,255,0.7)',
    auxiliary: 'rgba(255,255,255,0.5)',
    slave: 'rgba(255,255,255,0.05)',
    disabled: 'rgba(255,255,255,0.15)',
    inversePrimary: '#000',
    inverseAuxiliary: 'rgba(255,255,255,0.5)',
    inverseDisabled: 'rgba(255,255,255,0.15)',
  },
  divider: '#38405D',
  background: {
    paper0: '#111936',
    paper: '#1A223F',
    paper2: '#212A46',
    default: 'rgba(255,255,255,0.6)',
    disabled: '#2E375F',
    chip: '#232D4F',
    circle: '#3B476A',
  },
  common: {},
  shadowColor: '#0B1233',
  table: {
    head: {
      backgroundColor: '#232D4F',
      color: 'rgba(255,255,255,0.7)',
    },
    row: {
      hoverColor: '#212A46',
    },
    cell: {
      borderColor: '#232D4F',
    },
  },

  charts: {
    color: ['#7267EF', '#02BFA5'],
  },
} as const

export const light = {
  primary: {
    main: '#7635DC',
    contrastText: '#fff',
  },
  secondary: {
    lighter: '#D6E4FF',
    light: '#84A9FF',
    main: '#3366FF',
    dark: '#1939B7',
    darker: '#091A7A',
    contrastText: '#fff',
  },
  info: {
    lighter: '#D0F2FF',
    light: '#74CAFF',
    main: '#1890FF',
    dark: '#0C53B7',
    darker: '#04297A',
    contrastText: '#fff',
  },

  success: {
    lighter: '#E9FCD4',
    light: '#AAF27F',
    main: '#02BFA5',
    mainShadow: '#02BFA5',
    dark: '#229A16',
    darker: '#08660D',
    contrastText: 'rgba(0,0,0,0.7)',
  },

  warning: {
    lighter: '#FFF7CD',
    light: '#FFE16A',
    main: '#FFBF00',
    dark: '#B78103',
    darker: '#7A4F01',
    contrastText: 'rgba(0,0,0,0.7)',
  },
  neutral: {
    main: '#F4F6F8',
    contrastText: 'rgba(0, 0, 0, 0.60)',
  },
  error: {
    lighter: '#FFE7D9',
    light: '#FFA48D',
    main: '#FF1844',
    dark: '#B72136',
    darker: '#7A0C2E',
    contrastText: '#fff',
  },
  divider: '#E3E8EF',
  text: {
    primary: '#000',
    secondary: 'rgba(0,0,0,0.7)',
    auxiliary: 'rgba(0,0,0,0.5)',
    slave: 'rgba(0,0,0,0.05)',
    disabled: 'rgba(0,0,0,0.15)',
    inversePrimary: '#fff',
    inverseAuxiliary: 'rgba(255,255,255,0.5)',
    inverseDisabled: 'rgba(255,255,255,0.15)',
  },
  background: {
    paper0: '#fff',
    paper: '#fff',
    paper2: '#f6f8fa',
    default: '#fff',
    chip: '#F4F6F8',
    circle: '#E6E8EC',
  },

  shadowColor: 'rgba(145,158,171,0.2)',
  table: {
    head: {
      backgroundColor: '#F4F6F8',
      color: 'rgba(0,0,0,0.7)',
    },
    row: {
      hoverColor: '#F9FAFB',
    },
    cell: {
      borderColor: '#F3F4F5',
    },
  },
  charts: {
    color: ['#673AB7', '#02BFA5'],
  },
}

export type Color = typeof dark
