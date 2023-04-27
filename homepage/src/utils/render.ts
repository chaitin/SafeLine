import ReactDOM, { type Root } from 'react-dom/client'

const MARK = '__ct_react_root__'

type ContainerType = (Element | DocumentFragment) & {
  [MARK]?: Root
}

export function render(node: React.ReactElement, container: ContainerType) {
  const root = container[MARK] || ReactDOM.createRoot(container)

  root.render(node)

  container[MARK] = root
}

export async function unmount(container: ContainerType) {
  return Promise.resolve().then(() => {
    container[MARK]?.unmount()
    delete container[MARK]
  })
}
