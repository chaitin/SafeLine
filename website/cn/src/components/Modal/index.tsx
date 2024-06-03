import { type ConfirmDialogProps } from './ConfirmDialog'
import OriginModal from './Modal'
import confirm from './confrim'

type ModalStaticFunctions = Record<'confirm', (config: ConfirmDialogProps) => void>
type ModalType = typeof OriginModal

const Modal = OriginModal as ModalType & ModalStaticFunctions

Modal.confirm = confirm

export default Modal
