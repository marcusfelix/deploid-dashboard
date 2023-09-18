import { Record } from "pocketbase";
import { client } from "../client";
import { useQuery } from "react-query";
import { useRef } from "preact/hooks";
import TaskTitle from "../components/task-tile";
import TaskComposer from "../components/task-composer";
import { IconInfoCircle } from "@tabler/icons-preact";

type Props = {
    project: Record,
    matches: any
}

export default function Tasks(props: Props){
    const ref = useRef<HTMLDivElement>(null)

    const query = useQuery(`tasks-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/issues`, {
        method: "GET"
    }));
    
    return <div class="h-screen w-full md:flex-1 flex flex-col shrink-0 snap-center overflow-hidden mx-auto max-w-4xl">
        <div class="py-4">
            <h1 class="text-xl font-semibold">Tasks</h1>
        </div>
        <TaskComposer project={props.project} />
        <div class="flex flex-col flex-1 pt-6 gap-12 overflow-auto">
            {(query.data ?? []).length == 0 ? <span class="flex gap-2 items-center text-zinc-400">
                <IconInfoCircle size={20} class="shrink-0"/>
                Place your tasks, issues and bugs here. To discuss ideias and features, use the "Chat bubble" instead.
            </span> : null}
            {(query.data ?? []).map((task: Task) => <TaskTitle project={props.project} task={task} key={task.issue}/>)}
            <div ref={ref}></div>
        </div>
    </div>
}