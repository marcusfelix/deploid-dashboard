import { Record } from "pocketbase";
import { client } from "../client";
import { IconAlertHexagon, IconReload } from "@tabler/icons-preact";
import { useState } from "preact/hooks";

export default function Activate() {
    const [working, setWorking] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const user = client.authStore.model as Record

    const activate = async () => {
        setWorking(true);
        try {
            await client.collection("users").requestVerification(user.email);
            window.alert(`An e-mail was sent to ${user.email} with instructions to verify your account`)
        } catch(e: any){
            setError(e.message)
        }
        setWorking(false);
    }

    const refresh = async () => {
        setWorking(true);
        try {
            const data = await client.collection("users").authRefresh();
            if(!data.record.verified){
                window.alert(`Your account is not verified yet`)
            } else {
                client.authStore.save(data.token, data.record)
            }
        } catch(e: any){
            setError(e.message)
        }
        setWorking(false);
    }

    return <div class="flex items-center justify-center h-screen w-screen">
        <div class="max-w-lg px-6 py-10 overflow-x-auto">
            <div class="flex flex-col gap-4 justify-center">
                <img src="/logo.svg" class="w-10 h-10" />
                <h1 class="font-bold text-4xl">Activate your account</h1>
                <p class="text-xl text-gray-400">Please, activate your account to continue</p>
                <div class="h-6"></div>
                {error ? <p class="flex gap-4 text-xl text-red-500"><IconAlertHexagon /> {error}</p> : null}
                <button
                    class={`h-14 px-4 text-xl bg-black text-white rounded-lg cursor-pointer`}
                    onClick={refresh}
                    disabled={working}
                >{working ? <IconReload class="spin mx-auto" /> : "Refresh"}</button>
                <button
                    class="p-4 text-xl rounded-lg"
                    onClick={activate}
                    disabled={working}
                >Resend activation e-mail</button>
            </div>
        </div>
    </div>
}