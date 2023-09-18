import { Record } from "pocketbase";
import { client } from "../client";
import { useQuery, useQueryClient } from "react-query";
import { IconMessage } from "@tabler/icons-preact";
import CommentTile from "./comment-tile";
import { ReactMarkdown } from "react-markdown/lib/react-markdown";
import CommentComposer from "./comment-composer";

type Props = {
    task: Task,
    project: Record,
}

export default function TaskTitle(props: Props){
    const query = useQueryClient()

    const comments = useQuery(`comments-${props.project.id}-${props.task.issue}`, async () => client.send(`/api/projects/${props.project.slug}/issues/${props.task.issue}/comments`, {
        method: "GET"
    }));

    const change = async (e: any) => {
        await client.send(`/api/projects/${props.project.slug}/issues/${props.task.issue}/state`, {
            method: "POST",
            body: JSON.stringify({ state: e.target.value })
        })

        // Invalidade tasks
        query.invalidateQueries({ queryKey: [`tasks-${props.project.id}`] })
    }

    const user = () => {
        if(props.task.body.includes("－")){
            const string = props.task.body.substring(props.task.body.indexOf("－ ") + 2, props.task.body.length)
            return {
                name: string.split(" ")[0],
                login: string.split(" ")[0],
                photo: ""
            }
        }
        return props.task.user
    }

    const body = () => {
        if(props.task.body.includes("－")){
            return props.task.body.substring(0, props.task.body.indexOf("－ "),)
        }
        return props.task.body
    }
    
    return <div class={`flex flex-col gap-2`}>
        <div class="flex items-center gap-4">
            <h3 class="text-lg py-1 font-bold">{props.task.title}</h3>
            <div class="flex-1"></div>
            <select value={props.task.state} class={`${props.task.state === "open" ? "bg-green-200 text-green-700" : "bg-purple-200 text-purple-700"} font-bold text-sm px-3 py-1 rounded-full uppercase appearance-none text-center`} onChange={change}>
                <option value="open">Open</option>
                <option value="closed">Closed</option>
            </select>
        </div>
       
        <div class="flex gap-4">
            <div class={`flex gap-2`}>
                <div class="w-8 h-8 flex items-center justify-center font-bold text-sm aspect-square bg-red-100 text-red-800 rounded-full">
                    {user().photo != "" ? <img src={user().photo} class="w-8 h-8 rounded-full object-cover"/> : <span>{(user().name ?? user().login)[0]}</span>}
                </div>
            </div>
            <div class="flex flex-col gap-2 flex-1">
                <div class={`prose-lg text-gray-600`}>
                    <ReactMarkdown>{body()}</ReactMarkdown>
                </div>
                <div class="flex gap-2 items-center text-gray-400">
                    <span class="font-bold">{user().name ?? user().login}</span> • <span class="text-sm">{new Date(props.task.created).toLocaleDateString()} @ {new Date(props.task.created).toLocaleTimeString()}</span>
                    <div class="flex-1"></div>
                    {props.task.comments} <IconMessage class="hover:text-black cursor-pointer" size={20}/>
                </div>
                <div class="flex flex-col gap-6 mt-4 border-l-2 border-dashed pl-4 max-w-xl">
                    {(comments.data ?? []).map((e: Comment, i: number) => <CommentTile comment={e} key={i} />)}
                    <CommentComposer project={props.project} task={props.task}/>
                </div>
            </div>
        </div>
        
        
    </div>
}