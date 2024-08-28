export interface Plugin {
    name: string;
    execute: () => void;
}