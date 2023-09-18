import { IconCircleCheck, IconReload } from "@tabler/icons-preact"
import { Record } from "pocketbase"
import { useState } from "preact/hooks"
import { client } from "../client"
import { useQueryClient } from "react-query"

type Props = {
    project: Record
    task: Task
}

export default function CommentComposer(props: Props){
    const query = useQueryClient()
    const [body, setBody] = useState("")
    const [working, setWorking] = useState(false)

    const send = async () => {
        setWorking(true);
        await client.send(`/api/projects/${props.project.slug}/issues/${props.task.issue}/comments`, {
            method: "POST",
            body: JSON.stringify({ body })
        })
        await query.invalidateQueries({ queryKey: [`comments-${props.project.id}-${props.task.issue}`] })
        setBody("");
        setWorking(false);
    }

    const height = (e: any) => {
        e.target.style.height = 'inherit';
        e.target.style.height = `${e.target.scrollHeight}px`;
    }

    const validate = () => body.length > 0

    return <div class="flex flex-col bg-gray-100 rounded-lg overflow-hidden">
        <div class="flex">
            <textarea 
                value={body} 
                class="pr-4 py-2 text-lg max-h-32 ml-4 flex-1 bg-transparent outline-none" 
                placeholder="Type your message" 
                rows={1} 
                onInput={((e: any) => setBody(e.target.value))}
                onKeyUp={(e) => height(e)}
            ></textarea>
            <button 
                class={`font-semibold py-1 px-4 rounded-tr-lg rounded-br-lg ${validate() ? 'text-black' : 'text-gray-400'}`} 
                disabled={!validate()} 
                onClick={send}
            >{working ? <IconReload class="spin"/> : <IconCircleCheck/>}</button>
        </div>
    </div> 
}