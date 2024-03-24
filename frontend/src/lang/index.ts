import {useAppStore} from "@/store/modules/app";

import enLocale from "./package/en";
import zhCnLocale from "./package/zh_cn";
import {createI18n} from "vue-i18n";


const appStore = useAppStore();

const {language} = storeToRefs(appStore)

const messages = {
    "zh-cn": {
        ...zhCnLocale,
    },
    en: {
        ...enLocale,
    },
};


const i18n = createI18n({
    locale: language.value,
    legacy: false,
    messages: messages,
    globalInjection: true,
});

export default i18n;
