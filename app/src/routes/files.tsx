import { Record } from "pocketbase";
import { IconAlertCircle, IconFolder, IconReload } from "@tabler/icons-preact";
import { useQuery, useQueryClient } from "react-query";
import { client } from "../client";
import { useState } from "preact/hooks";
import TextEditor from "../components/text-editor";
import MarkdownViewer from "../components/markdown-viewer";
import ImageEditor from "../components/image-editor";

type Props = {
    project: Record
}

export default function Files(props: Props){
    const query = useQueryClient()
    const [working, setWorking] = useState(false)
    const [file, setFile] = useState("lib/includes/default.dart")
    const [content, setContent] = useState<string | null>(null)
    const [message, setMessage] = useState<string | null>(null)

    const manifest = useQuery(`manifest-${props.project.id}`, async () => client.send(`/api/projects/${props.project.id}/read`, {
        method: "POST",
        body: JSON.stringify({
            path: "app/manifest.json"
        })
    }).then((data) => JSON.parse(atob(data.content))), { staleTime: 1000 * 60 * 60 * 24 });

    const contents = useQuery(`file-${props.project.id}-${file}`, async () => file == "" ? null : client.send(`/api/projects/${props.project.id}/read`, {
        method: "POST",
        body: JSON.stringify({
            path: `app/${file}`
        })
    }));

    const save =  async (sha: string) => {
        setWorking(true)
        try {
            await client.send(`/api/projects/${props.project.slug}/write`, {
                method: "POST",
                body: JSON.stringify({
                    path: `app/${file}`,
                    content: content,
                    sha,
                    message
                })
            })
        } catch(e){
            console.error(e)
        }
        await query.invalidateQueries({ queryKey: [`file-${props.project.id}-${file}`] })
        setWorking(false)
        setMessage(null)
        setContent(null)
    }

    const render = () => {
        const name = (file.split("/") as any).at(-1)
        const extention = name.split(".").at(-1)
        switch(extention){
            case "md":
                return <MarkdownViewer data={contents.data.content ?? ""}/>

            case "png":
            case "jpg":
            case "jpeg":
            case "gif":
                return <ImageEditor extention={extention} data={contents.data.content ?? ""} onChange={(data: string) => setContent(data)}/>
                
            default:
                return <TextEditor value={contents.data.content ?? ""} onChange={(data: string) => setContent(data)}/>
        }
    }

    return <div class="w-full md:flex-1 flex flex-col pb-4 shrink-0 snap-center mx-auto max-w-4xl">
        <div class="py-4">
            <h1 class="text-2xl font-bold">Code</h1>
        </div>
        <div class="flex items-center gap-4 pb-2 text-xl text-gray-500 shrink-0">
            <IconFolder size={24}/>
            /
            <select defaultValue={file} onChange={(e: any) => {
                setFile(e.target.value);
                setContent(null);
            }} disabled={content != null}>
                {Object.keys(manifest.data ?? {}).map((key) => <option value={manifest.data[key]}>{key}</option>)}
            </select>
            {content != null ? <span class="text-[18px] flex gap-2 items-center text-red-500">
                <IconAlertCircle size={24}/> Pending changes 
                <button class="text-md border-[2px] border-red-500 px-2 rounded-md text-red-600" onClick={() => setContent(null)}>Discard</button>
            </span> : null}
        </div>
        <div class="h-full overflow-auto">
            {contents.data ? render() : null}
        </div>
        <div class="h-4"></div>
        <div class="p-2 flex items-center gap-4 bg-gray-100 rounded-lg">
            <input defaultValue={message ?? ""} class="flex-1 px-4 py-2 text-lg bg-transparent" placeholder="Describe your change here" onInput={(e: any) => setMessage(e.target.value)}/>
            <button class={`font-semibold py-2 px-4 rounded-md ${content ? 'bg-black text-white' : 'bg-gray-300 text-gray-500'}`} disabled={!content} onClick={() => save(contents.data.sha)}>{working ? <IconReload class="spin"/> : "Save"}</button>
        </div>
    </div>
}