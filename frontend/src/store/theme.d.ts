import type { UserPreferences } from "../types";
type Theme = "light" | "dark";
interface ThemeStore {
    theme: Theme;
    isDark: boolean;
    setTheme: (theme: Theme) => void;
    toggleTheme: () => void;
    resetTheme: () => void;
}
export declare const useThemeStore: import("zustand").UseBoundStore<Omit<import("zustand").StoreApi<ThemeStore>, "persist"> & {
    persist: {
        setOptions: (options: Partial<import("zustand/middleware").PersistOptions<ThemeStore, {
            theme: Theme;
        }>>) => void;
        clearStorage: () => void;
        rehydrate: () => Promise<void> | void;
        hasHydrated: () => boolean;
        onHydrate: (fn: (state: ThemeStore) => void) => () => void;
        onFinishHydration: (fn: (state: ThemeStore) => void) => () => void;
        getOptions: () => Partial<import("zustand/middleware").PersistOptions<ThemeStore, {
            theme: Theme;
        }>>;
    };
}>;
interface UserPreferencesStore extends UserPreferences {
    autoSaveProgress: boolean;
    showHints: boolean;
    solutionDetail: "simple" | "detailed" | "comprehensive";
    updatePreferences: (preferences: Partial<UserPreferences>) => void;
    resetPreferences: () => void;
}
export declare const useUserPreferencesStore: import("zustand").UseBoundStore<Omit<import("zustand").StoreApi<UserPreferencesStore>, "persist"> & {
    persist: {
        setOptions: (options: Partial<import("zustand/middleware").PersistOptions<UserPreferencesStore, UserPreferencesStore>>) => void;
        clearStorage: () => void;
        rehydrate: () => Promise<void> | void;
        hasHydrated: () => boolean;
        onHydrate: (fn: (state: UserPreferencesStore) => void) => () => void;
        onFinishHydration: (fn: (state: UserPreferencesStore) => void) => () => void;
        getOptions: () => Partial<import("zustand/middleware").PersistOptions<UserPreferencesStore, UserPreferencesStore>>;
    };
}>;
export {};
//# sourceMappingURL=theme.d.ts.map