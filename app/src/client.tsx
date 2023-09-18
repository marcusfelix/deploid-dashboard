import PocketBase, { Admin, BaseAuthStore, Record } from "pocketbase";

class UserStore extends BaseAuthStore {

    constructor() {
        super();
        const token = localStorage.getItem("token");
        const model = localStorage.getItem("model");
        if (token && model) {
            this.save(token, JSON.parse(model));
        }
    }

    save(token: string, model: Record | Admin | null): void {
        super.save(token, model);
        localStorage.setItem("token", token);
        localStorage.setItem("model", JSON.stringify(model));
    }

    clear(): void {
        super.clear();
        localStorage.removeItem("token");
        localStorage.removeItem("model");
    }
}

export const client = new PocketBase(import.meta.env.VITE_SERVER_URL)