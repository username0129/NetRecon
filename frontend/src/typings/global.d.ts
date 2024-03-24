export {}; // 作为一个模块使用

declare global {
    /**
     * 系统设置
     */
    interface AppSettings {
        /** 系统标题 */
        title: string;
        /** 系统版本 */
        version: string;
        /** 语言( zh-cn | en) */
        language: string;
    }
}

