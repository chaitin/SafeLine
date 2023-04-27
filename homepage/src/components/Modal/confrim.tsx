import { render as reactRender } from "@/utils";

import ConfirmDialog, { type ConfirmDialogProps } from "./ConfirmDialog";

export default function confirm(config: ConfirmDialogProps) {
  const container = document.createDocumentFragment();
  const { onCancel: propCancel, onOk: propOk } = config;
  const onCancel = async () => {
    await propCancel?.();
    close();
  };
  const onOk = async () => {
    await propOk?.();
    close();
  };
  let currentConfig = { ...config, open: true, onCancel, onOk } as any;
  function render(props: ConfirmDialogProps) {
    setTimeout(() => {
      reactRender(<ConfirmDialog {...props} />, container);
    });
  }

  function close() {
    currentConfig = {
      ...currentConfig,
      open: false,
    };
    render(currentConfig);
  }

  render(currentConfig);
}
