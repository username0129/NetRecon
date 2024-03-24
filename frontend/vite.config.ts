import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import {ElementPlusResolver} from "unplugin-vue-components/resolvers";
import Icons from "unplugin-icons/vite";
import IconsResolver from "unplugin-icons/resolver";

import * as path from 'path';
import {resolve} from "path";
import {createSvgIconsPlugin} from "vite-plugin-svg-icons";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        // 自动导入参考： https://github.com/sxzz/element-plus-best-practices/blob/main/vite.config.ts
        AutoImport({
            // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
            imports: ["vue", "@vueuse/core", "pinia", "vue-router", "vue-i18n"],
            resolvers: [
                // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox... (带样式)
                ElementPlusResolver(),
                // 自动导入图标组件
                IconsResolver({}),
            ],
            eslintrc: {
                // 是否自动生成 eslint 规则，建议生成之后设置 false
                enabled: false,
                // 指定自动导入函数 eslint 规则的文件
                filepath: "./.eslintrc-auto-import.json",
                globalsPropValue: true,
            },
            // 是否在 vue 模板中自动导入
            vueTemplate: true,
            // 指定自动导入函数TS类型声明文件路径 (false:关闭自动生成)
            // dts: false,
            dts: "src/typings/auto-imports.d.ts",
        }),
        Components({
            resolvers: [
                // 自动导入 Element Plus 组件
                ElementPlusResolver(),
                // 自动注册图标组件
                IconsResolver({
                    // element-plus图标库，其他图标库 https://icon-sets.iconify.design/
                    enabledCollections: ["ep"],
                }),
            ],
            // 指定自定义组件位置(默认:src/components)
            dirs: ["src/components", "src/**/components"],
            // 指定自动导入组件TS类型声明文件路径 (false:关闭自动生成)
            // dts: false,
            dts: "src/typings/components.d.ts",
        }),
        Icons({
            // 自动安装图标库
            autoInstall: true,
        }),
        createSvgIconsPlugin({
            // 指定需要缓存的图标文件夹
            iconDirs: [resolve(path.resolve(__dirname, './src'), "assets/icons")],
            // 指定symbolId格式
            symbolId: "icon-[dir]-[name]",
        }),
    ],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src')
        }, // 指示 Vite 如何处理特定的别名
    },
});
