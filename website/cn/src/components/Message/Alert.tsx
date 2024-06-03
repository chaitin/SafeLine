import React, { type FC, useRef, useEffect } from 'react'

import { Alert as MAlert, type AlertColor } from '@mui/material'

interface AlertProps {
  duration?: number
  onClose?(key: React.Key): void
  noticeKey: React.Key
  content?: React.ReactNode
  severity: AlertColor
}

const Alert: FC<AlertProps> = (props) => {
  const { duration, severity, content, noticeKey, onClose } = props
  const closeTimer = useRef<number | null>(null)

  const startCloseTimer = () => {
    if (duration) {
      closeTimer.current = window.setTimeout(() => {
        close()
      }, duration * 1000)
    }
  }

  const clearCloseTimer = () => {
    if (closeTimer.current) {
      clearTimeout(closeTimer.current)
      closeTimer.current = null
    }
  }

  const close = (e?: React.MouseEvent<HTMLAnchorElement>) => {
    if (e) {
      e.stopPropagation()
    }
    clearCloseTimer()
    if (onClose) {
      onClose(noticeKey)
    }
  }

  useEffect(() => {
    startCloseTimer()
    return () => {
      clearCloseTimer()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <MAlert severity={severity} sx={{ mb: '10px' }}>
      {content}
    </MAlert>
  )
}

export default Alert
