import { IconCircleCheck, IconReload } from "@tabler/icons-preact";
import { useState } from "preact/hooks";
import { client } from "../client";
import { Record } from "pocketbase";
import { useQueryClient } from "react-query";

type Props = {
    project: Record,
}

export default function TaskComposer(props: Props){
    const [working, setWorking] = useState(false)
    const [title, setTitle] = useState("")
    const [body, setBody] = useState("")

    const query = useQueryClient()

    const send = async () => {
        setWorking(true);
        
        
        await client.send(`/api/projects/${props.project.slug}/issues`, {
            method: "POST",
            body: JSON.stringify({
                title: title,
                body: body,
            })
        })

        setBody("");
        setTitle("");
        setWorking(false);

        // Invalidade tasks
        query.invalidateQueries({ queryKey: [`tasks-${props.project.id}`] })
    }

    const height = (e: any) => {
        e.target.style.height = 'inherit';
        e.target.style.height = `${e.target.scrollHeight}px`;
    }

    const validate = () => body.length > 0 && title.length > 3

    return <div class="flex py-1 bg-gray-100 rounded-lg overflow-hidden">
        <div class="flex flex-col flex-1 pl-4">
            <input value={title} placeholder="Your new task" class="pr-4 py-2 text-xl bg-transparent outline-none" onInput={(e: any) => setTitle(e.target.value)}/>
            {title.length > 3 ? <textarea 
                value={body} 
                class="pr-4 py-2 text-lg max-h-32 bg-transparent outline-none" 
                placeholder="Elaborate on your task" 
                rows={1} 
                onInput={((e: any) => setBody(e.target.value))}
                onKeyUp={(e) => height(e)}
            ></textarea> : null}
        </div>
        <button class={`font-semibold py-1 px-4 rounded-tr-lg rounded-br-lg ${validate() ? 'text-black' : 'text-gray-400'}`} disabled={!validate()} onClick={send}>{working ? <IconReload class="spin"/> : <IconCircleCheck/>}</button>
    </div> 
}