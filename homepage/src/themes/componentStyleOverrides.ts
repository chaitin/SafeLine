import { type Color } from './color'
export default function componentStyleOverrides(color: Color) {
  return {
    MuiButton: {
      styleOverrides: {
        root: ({ ownerState, theme }: any) => {
          return {
            boxShadow: 'none',
            '&:hover': {
              boxShadow: 'none',
              ...(ownerState.color === 'neutral' && {
                color: color.primary.main,
                fontWeight: 700,
              }),
            },
          }
        },
      },
    },
    MuiLoadingButton: {
      styleOverrides: {
        root: ({ ownerState, theme }: any) => {
          return {
            // position: 'relative',
            // ...(ownerState.loading && {
            //   backgroundImage: 'linear-gradient(263deg, #3892EE 0%, #2E7AE9 100%)',
            //   boxShadow: '0px 2px 8px 0px rgba(133,190,245,0.5)',
            //   '&.MuiLoadingButton-loading': {
            //     color: '#fff',
            //   },
            //   '&:before': {
            //     content: "''",
            //     position: 'absolute',
            //     top: '-1px',
            //     right: '-1px',
            //     bottom: '-1px',
            //     left: '-1px',
            //     zIndex: 1,
            //     background: '#fff',
            //     opacity: 0.35,
            //     transition: 'opacity .2s',
            //     pointerEvents: 'none',
            //   },
            // }),
          }
        },
      },
    },

    MuiListItemButton: {
      styleOverrides: {
        root: {
          // paddingTop: '10px',
          // paddingBottom: '10px',
          // borderRadius: '8px',
          // '&.Mui-selected': {
          //   color: theme.menuSelected,
          //   backgroundColor: theme.menuSelectedBack,
          //   '&:hover': {
          //     backgroundColor: theme.menuSelectedBack,
          //   },
          //   '& .MuiListItemIcon-root': {
          //     color: theme.menuSelected,
          //   },
          // },
          // '&:hover': {
          //   backgroundColor: 'transparent',
          //   color: theme.menuSelected,
          //   '& .MuiListItemIcon-root': {
          //     color: theme.menuSelected,
          //   },
          // },
        },
      },
    },
    MuiListItemIcon: {
      styleOverrides: {
        root: {
          // minWidth: '36px',
        },
      },
    },
    MuiFormControlLabel: {
      styleOverrides: {
        root: {
          // marginRight: '24px',
        },
      },
    },
    MuiInputBase: {
      styleOverrides: {
        root: {
          // '.MuiInputBase-colorPrimary&:hover:not(.Mui-disabled):before': {
          //   borderColor: theme.colors?.primaryMain,
          // },
        },
      },
    },

    MuiFormControl: {
      styleOverrides: {
        root: {
          '.MuiFormLabel-asterisk': {
            color: color.error.main,
          },
        },
      },
    },
    MuiTableRow: {
      styleOverrides: {
        root: {
          '&:hover': {
            '.MuiTableCell-root': {
              // backgroundColor: color.table.row.hoverColor,
            },
          },
        },
      },
    },
    MuiTableBody: {
      styleOverrides: {
        root: {
          '.MuiTableRow-root:hover': {
            '.MuiTableCell-root': {
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
          fontSize: '14px',
          paddingTop: '24px',
          paddingBottom: '24px',
          borderColor: color.table.cell.borderColor,
          paddingLeft: 0,
          '&:first-of-type': {
            paddingLeft: '32px',
          },
        },
        head: {
          backgroundColor: color.table.head.backgroundColor,
          color: color.table.head.color,
          fontSize: '12px',
          height: '24px',
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
            backgroundImage: 'none',
          }
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: '4px'
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
          '&:hover .MuiOutlinedInput-notchedOutline': {
            borderColor: color.primary.main,
          },
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: color.divider,
          },
        },
      },
    },
  }
}
