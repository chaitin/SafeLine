import React, { type FC, useState } from 'react'

import CloseIcon from '@mui/icons-material/Close'
import { LoadingButton } from '@mui/lab'
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton } from '@mui/material'

export interface ModalProps {
  open?: boolean
  title?: React.ReactNode
  children?: React.ReactNode
  footer?: false | React.ReactNode
  okText?: React.ReactNode
  cancelText?: React.ReactNode
  showCancel?: boolean
  okColor?: 'primary' | 'error'
  cancelColor?: 'primary' | 'error'
  closable?: boolean
  onOk?(): void
  onClose?(): void
  onCancel?(): void
  sx?: any
}

const Modal: FC<ModalProps> = (props) => {
  const {
    open = false,
    title,
    children,
    footer,
    okText = '确认',
    okColor = 'primary',
    cancelColor = 'primary',
    showCancel = true,
    cancelText = '取消',
    onOk,
    onClose,
    onCancel,
    closable = true,
    sx = {},
  } = props
  const [loading, setLoading] = useState(false)

  const onConfirm = async () => {
    setLoading(true)
    try {
      await onOk?.()
    } catch (error) {}

    setLoading(false)
  }

  return (
    <Dialog
      open={open}
      onClose={onClose || onCancel}
      PaperProps={{
        elevation: 0,
      }}
      sx={{
        '.MuiDialog-paper': {
          borderRadius: '12px',
          ...sx,
        },
      }}
    >
      {(title || closable) && (
        <DialogTitle
          sx={{
            fontWeight: 600,
            fontSize: '16px',
            p: '32px',
            pb: '16px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          {title}
          {closable && (
            <IconButton
              sx={{
                color: 'text.auxiliary',
              }}
              onClick={onClose || onCancel}
            >
              <CloseIcon sx={{ fontSize: '20px' }} />
            </IconButton>
          )}
        </DialogTitle>
      )}

      <DialogContent sx={{ px: '32px' }}>{children}</DialogContent>
      {footer === false && null}
      {footer === undefined && (
        <DialogActions
          sx={{
            p: '0 32px 32px',
            '.MuiButtonBase-root': {
              px: '30px',
            },
          }}
        >
          {showCancel && (
            <Button color={cancelColor} onClick={onCancel}>
              {cancelText}
            </Button>
          )}

          <LoadingButton loading={loading} variant='contained' color={okColor} onClick={onConfirm}>
            {okText}
          </LoadingButton>
        </DialogActions>
      )}
      {footer && (
        <DialogActions
          sx={{
            p: '4px 24px 24px',
            '.MuiButtonBase-root': {
              p: '8px 30px',
            },
          }}
        >
          {footer}
        </DialogActions>
      )}
    </Dialog>
  )
}

export default Modal
