import { Record } from "pocketbase";
import { IconLock, IconReload } from "@tabler/icons-preact";
import { useQuery } from "react-query";
import { client } from "../client";
import { useState } from "preact/hooks";

type Props = {
    project: Record
}

export default function Secrets(props: Props){
    const [saving, setSaving] = useState(false)
    const [payload, setPayload] = useState<any>({})

    const secrets = useQuery(`deploy-secrets-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/secrets`, {
        method: "GET"
    }))

    const save = async () => {
        setSaving(true);
        const batch = Object.keys(payload).map((key: string) => {
            return client.send(`/api/projects/${props.project.slug}/secrets`, {
                method: "POST",
                body: JSON.stringify({
                    key: key,
                    value: payload[key]
                })
            })
        })
        try {
            await Promise.all(batch)
        } catch(e){
            console.log(e)
        }
        setPayload({})
        setSaving(false)
    }

    const working = () => secrets.status !== "success" || saving;

    return <div class="w-full md:flex-1 flex flex-col gap-6 shrink-0 snap-center mx-auto max-w-4xl">
        <div class="py-4">
            <h1 class="text-xl font-semibold">Secrets</h1>
        </div>
        <div class="flex flex-col gap-4">
            {(secrets.data ?? []).map((secret: any) => <div class="flex flex-col gap-2 bg-yellow-100 text-yellow-800 rounded-xl">
                <div class="flex gap-4 px-4 items-center">
                    <div title={new Date(secret.updated).toLocaleDateString() + " " + new Date(secret.updated).toLocaleTimeString()}><IconLock /></div>
                    <input type="password" placeholder={secret.name} class="w-full py-4 text-yellow-800 bg-transparent placeholder:text-yellow-800/30 outline-none" onInput={(e: any) => setPayload({...payload, ...{ [secret.name]: e.target.value }})}/>
                </div>
            </div>)}
        </div>
        <div class="flex items-center gap-4">
            <div class="flex-1 text-lg text-gray-500">Securely send your sensitive information directly to the build server</div>
            <button class="bg-black text-white font-semibold py-2 px-4 rounded-md" disabled={working()} onClick={() => save()}>{working() ? <IconReload class="spin"/> : "Save"}</button>
        </div>
    </div>
}