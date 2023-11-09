import { useCallback, useRef, useLayoutEffect } from 'react'

export function useEvent<Fn extends (...args: any[]) => any>(handler: Fn): Fn {
  if (typeof window === 'undefined') {
    return handler
  }
  // eslint-disable-next-line
  return useCallbackEvent(handler)
}

export function useCallbackEvent<Fn extends (...args: any[]) => any>(handler: Fn): Fn {

  const handlerRef = useRef<any>(null)

  useLayoutEffect(() => {
    handlerRef.current = handler
  })

  return useCallback((...args: any[]): any => {
    handlerRef.current?.(...args)
  }, []) as Fn
}