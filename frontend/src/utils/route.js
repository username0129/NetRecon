const viewModules = import.meta.glob('../views/**/*.vue')
const pluginModules = import.meta.glob('../plugins/**/*.vue')

export const RouteHandle = (routes) => {
  routes.forEach((item) => {
    if (item.component && typeof item.component === 'string') {
      const moduleImporter = item.component.startsWith('views/')
        ? viewModules
        : item.component.startsWith('plugins/')
          ? pluginModules
          : null
      if (moduleImporter) {
        const componentImporter = dynamicImport(moduleImporter, item.component)
        if (componentImporter) {
          item.component = componentImporter
        }
      }
    }
    if (item.children) {
      RouteHandle(item.children)
    }
  })
}

function dynamicImport(modules, component) {
  const componentPath = `../${component}`
  const moduleFn = modules[componentPath]
  if (moduleFn) {
    return () => moduleFn()
  }
  throw new Error(`Component ${component} not found`)
}
