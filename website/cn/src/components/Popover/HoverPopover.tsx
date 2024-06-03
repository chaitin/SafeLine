import * as React from "react";
import { Popover, type PopoverProps } from '@mui/material'

const HoverPopover: React.ComponentType<PopoverProps> = React.forwardRef(
  function HoverPopover(props: PopoverProps, ref): any {
    return (
      <Popover
        {...props}
        ref={ref}
        style={{ pointerEvents: 'none', ...props.style }}
        PaperProps={{
          ...props.PaperProps,
          style: {
            pointerEvents: 'auto',
            ...props.PaperProps?.style,
            boxShadow: "0px 25px 50px 0px rgba(145,158,171,0.2)"
          },
        }}
      />
    )
  }
)

export default HoverPopover