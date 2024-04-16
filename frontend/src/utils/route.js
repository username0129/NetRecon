const viewModules = import.meta.glob('../views/**/*.vue')
const pluginModules = import.meta.glob('../plugins/**/*.vue')

export const RouteHandle = (routes) => {
    routes.forEach(item => {
        if (item.component && typeof item.component === 'string') {
            if (item.component.startsWith('views/')) {
                item.component = dynamicImport(viewModules, item.component)
            } else if (item.component.startsWith('plugins/')) {
                item.component = dynamicImport(pluginModules, item.component)
            }
        }
        if (item.children) {
            RouteHandle(item.children)
        }
    })
}
function dynamicImport(modules, component) {
    const componentPath = `../${component}`;
    const moduleFn = modules[componentPath];
    if (moduleFn) {
        return moduleFn();
    }
    throw new Error(`Component ${component} not found`);
}
