
import {
  type SyntheticEvent,
  type MouseEvent,
  type FocusEvent,
  useCallback,
  useState,
  useRef,
  useEffect,
} from 'react'
import { type PopoverPosition, type PopoverReference } from '@mui/material'
import { useEvent } from './useEvent'

type ControlAriaProps = {
  'aria-controls'?: string
  'aria-describedby'?: string
  'aria-haspopup'?: true
}

export type PopupState = {
  open: (eventOrAnchorEl?: SyntheticEvent | Element | null) => void
  close: () => void
  toggle: (eventOrAnchorEl?: SyntheticEvent | Element | null) => void
  onBlur: (event: FocusEvent) => void
  onMouseLeave: (event: MouseEvent) => void
  setOpen: (
    open: boolean,
    eventOrAnchorEl?: SyntheticEvent | Element | null
  ) => void
  isOpen: boolean
  anchorEl: Element | undefined
  anchorPosition: PopoverPosition | undefined
  setAnchorEl: (anchorEl: Element | null | undefined) => any
  setAnchorElUsed: boolean
  disableAutoFocus: boolean
  popupId: string | undefined
  _openEventType: string | null | undefined
  _childPopupState: PopupState | null | undefined
  _setChildPopupState: (popupState: PopupState | null | undefined) => void
}

export type CoreState = {
  isOpen: boolean
  setAnchorElUsed: boolean
  anchorEl: Element | undefined
  anchorPosition: PopoverPosition | undefined
  hovered: boolean
  focused: boolean
  _openEventType: string | null | undefined
  _childPopupState: PopupState | null | undefined
  _deferNextOpen: boolean
  _deferNextClose: boolean
}

export const initCoreState: CoreState = {
  isOpen: false,
  setAnchorElUsed: false,
  anchorEl: undefined,
  anchorPosition: undefined,
  hovered: false,
  focused: false,
  _openEventType: null,
  _childPopupState: null,
  _deferNextOpen: false,
  _deferNextClose: false,
}

export function bindPopover({
  isOpen,
  anchorEl,
  anchorPosition,
  close,
  popupId,
  onMouseLeave,
  disableAutoFocus,
  _openEventType,
}: PopupState): {
  id?: string
  anchorEl?: Element | null
  anchorPosition?: PopoverPosition
  anchorReference: PopoverReference
  open: boolean
  onClose: () => void
  onMouseLeave: (event: MouseEvent) => void
  disableAutoFocus?: boolean
  disableEnforceFocus?: boolean
  disableRestoreFocus?: boolean
} {
  const usePopoverPosition = _openEventType === 'contextmenu'
  return {
    id: popupId,
    anchorEl,
    anchorPosition,
    anchorReference: usePopoverPosition ? 'anchorPosition' : 'anchorEl',
    open: isOpen,
    onClose: close,
    onMouseLeave,
    ...(disableAutoFocus && {
      disableAutoFocus: true,
      disableEnforceFocus: true,
      disableRestoreFocus: true,
    }),
  }
}

function controlAriaProps({
  isOpen,
  popupId,
}: PopupState): ControlAriaProps {
  return {
    'aria-haspopup': true,
    'aria-controls': isOpen && popupId != null ? popupId : undefined,
  }
}

export function bindHover(popupState: PopupState): ControlAriaProps & {
  // onTouchStart: (event: TouchEvent) => any
  onMouseOver: (event: MouseEvent) => any
  onMouseLeave: (event: MouseEvent) => any
} {
  const { open, onMouseLeave } = popupState
  return {
    ...controlAriaProps(popupState),
    // onTouchStart: open,
    onMouseOver: open,
    onMouseLeave,
  }
}

function getPopupElement(
  element: Element,
  { popupId }: PopupState
): Element | null | undefined {
  if (!popupId) return null
  const rootNode: any =
    typeof element.getRootNode === 'function' ? element.getRootNode() : document
  if (typeof rootNode.getElementById === 'function') {
    return rootNode.getElementById(popupId)
  }
  return null
}

function isElementInPopup(element: Element, popupState: PopupState): boolean {
  const { anchorEl, _childPopupState } = popupState
  return (
    isAncestor(anchorEl, element) ||
    isAncestor(getPopupElement(element, popupState), element) ||
    (_childPopupState != null && isElementInPopup(element, _childPopupState))
  )
}

function isAncestor(
  parent: Element | null | undefined,
  child: Element | null | undefined
): boolean {
  if (!parent) return false
  while (child) {
    if (child === parent) return true
    child = child.parentElement
  }
  return false
}

export default function usePopupState({
  parentPopupState,
  popupId,
  disableAutoFocus,
}: {
  parentPopupState?: PopupState | null | undefined
  popupId?: string
  disableAutoFocus?: boolean | null | undefined
}): PopupState {
  const isMounted = useRef(true)

  useEffect((): (() => void) => {
    isMounted.current = true
    return () => {
      isMounted.current = false
    }
  }, [])

  const [state, _setState] = useState(initCoreState)

  const setState = useCallback(
    (state: CoreState | ((prevState: CoreState) => CoreState)) => {
      if (isMounted.current) _setState(state)
    },
    [],
  )

  const setAnchorEl = useCallback(
    (anchorEl: Element | null | undefined) =>
      setState((state) => ({
        ...state,
        setAnchorElUsed: true,
        anchorEl: anchorEl ?? undefined,
      })),
    [],
  )

  const toggle = useEvent(
    (eventOrAnchorEl?: SyntheticEvent | Element | null) => {
      if (state.isOpen) close(eventOrAnchorEl)
      else open(eventOrAnchorEl)
      return state
    }
  )

  const open = useEvent((eventOrAnchorEl?: SyntheticEvent | Element | null) => {
    const event =
      eventOrAnchorEl instanceof Element ? undefined : eventOrAnchorEl
    const element =
      eventOrAnchorEl instanceof Element
        ? eventOrAnchorEl
        : eventOrAnchorEl?.currentTarget instanceof Element
        ? eventOrAnchorEl.currentTarget
        : undefined

    // if (event?.type === 'touchstart') {
    //   setState((state) => ({ ...state, _deferNextOpen: true }))
    //   return
    // }

    const clientX = (event as MouseEvent | undefined)?.clientX
    const clientY = (event as MouseEvent | undefined)?.clientY
    const anchorPosition =
      typeof clientX === 'number' && typeof clientY === 'number'
        ? { left: clientX, top: clientY }
        : undefined

    const doOpen = (state: CoreState): CoreState => {
      if (!eventOrAnchorEl && !state.setAnchorElUsed) {
        console.error(
          'missingEventOrAnchorEl',
          'eventOrAnchorEl should be defined if setAnchorEl is not used'
        )
      }

      if (parentPopupState) {
        if (!parentPopupState.isOpen) return state
        setTimeout(() => parentPopupState._setChildPopupState(popupState))
      }

      const newState: CoreState = {
        ...state,
        isOpen: true,
        anchorPosition,
        hovered: event?.type === 'mouseover' || state.hovered,
        focused: event?.type === 'focus' || state.focused,
        _openEventType: event?.type,
      }

      if (event?.currentTarget) {
        if (!state.setAnchorElUsed) {
          newState.anchorEl = event?.currentTarget as any
        }
      } else if (element) {
        newState.anchorEl = element
      }

      return newState
    }

    setState((state: CoreState): CoreState => {
      if (state._deferNextOpen) {
        setTimeout(() => setState(doOpen), 0)
        return { ...state, _deferNextOpen: false }
      } else {
        return doOpen(state)
      }
    })
  })


  const doClose = (state: CoreState): CoreState => {
    const { _childPopupState } = state
    setTimeout(() => {
      _childPopupState?.close()
      parentPopupState?._setChildPopupState(null)
    })
    return { ...state, isOpen: false, hovered: false, focused: false }
  }

  const close = useEvent(
    (eventOrAnchorEl?: SyntheticEvent | Element | null) => {
      const event =
        eventOrAnchorEl instanceof Element ? undefined : eventOrAnchorEl

      // if (event?.type === 'touchstart') {
      //   setState((state) => ({ ...state, _deferNextClose: true }))
      //   return
      // }

      setState((state: CoreState): CoreState => {
        if (state._deferNextClose) {
          setTimeout(() => setState(doClose), 0)
          return { ...state, _deferNextClose: false }
        } else {
          return doClose(state)
        }
      })
    }
  )

  const setOpen = useCallback(
    (
      nextOpen: boolean,
      eventOrAnchorEl?: SyntheticEvent<any> | Element | null
    ) => {
      if (nextOpen) {
        open(eventOrAnchorEl)
      } else {
        close(eventOrAnchorEl)
      }
    },
    []
  )

  const onMouseLeave = useEvent((event: MouseEvent) => {
    const { relatedTarget } = event
    setState((state: CoreState): CoreState => {
      if (
        state.hovered &&
        !(
          relatedTarget instanceof Element &&
          isElementInPopup(relatedTarget, popupState)
        )
      ) {
        if (state.focused) {
          return { ...state, hovered: false }
        } else {
          return doClose(state)
        }
      }
      return state
    })
  })

  const onBlur = useEvent((event: FocusEvent) => {
    if (!event) return
    const { relatedTarget } = event
    setState((state: CoreState): CoreState => {
      if (
        state.focused &&
        !(
          relatedTarget instanceof Element &&
          isElementInPopup(relatedTarget, popupState)
        )
      ) {
        if (state.hovered) {
          return { ...state, focused: false }
        } else {
          return doClose(state)
        }
      }
      return state
    })
  })

  const _setChildPopupState = useCallback(
    (_childPopupState: PopupState | null | undefined) =>
      setState((state) => ({ ...state, _childPopupState })),
    []
  )

  const popupState: PopupState = {
    ...state,
    setAnchorEl,
    popupId,
    open,
    close,
    toggle,
    setOpen,
    onBlur,
    onMouseLeave,
    disableAutoFocus:
      disableAutoFocus ?? Boolean(state.hovered || state.focused),
    _setChildPopupState,
  }

  return popupState
}