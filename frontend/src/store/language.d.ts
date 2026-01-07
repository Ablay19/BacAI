import type { Language } from "../types";
interface LanguageStore {
    language: Language;
    isRTL: boolean;
    supportedLanguages: {
        code: Language;
        name: string;
        nativeName: string;
    }[];
    setLanguage: (language: Language) => void;
    detectLanguage: () => void;
    toggleDirection: () => void;
}
export declare const useLanguageStore: import("zustand").UseBoundStore<Omit<import("zustand").StoreApi<LanguageStore>, "persist"> & {
    persist: {
        setOptions: (options: Partial<import("zustand/middleware").PersistOptions<LanguageStore, {
            language: Language;
        }>>) => void;
        clearStorage: () => void;
        rehydrate: () => Promise<void> | void;
        hasHydrated: () => boolean;
        onHydrate: (fn: (state: LanguageStore) => void) => () => void;
        onFinishHydration: (fn: (state: LanguageStore) => void) => () => void;
        getOptions: () => Partial<import("zustand/middleware").PersistOptions<LanguageStore, {
            language: Language;
        }>>;
    };
}>;
export {};
//# sourceMappingURL=language.d.ts.map