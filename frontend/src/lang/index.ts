import {useAppStoreHook} from "@/store/modules/app";
import {createI18n} from "vue-i18n";
import enLocale from "./package/en";
import zhCnLocale from "./package/zh_cn";

const appStore = useAppStoreHook();

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
