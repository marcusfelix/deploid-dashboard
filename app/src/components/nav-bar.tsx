import { IconApps, IconChecklist, IconCode, IconLock, IconWallet } from "@tabler/icons-preact";
import { Link } from 'preact-router/match';
import { Record } from "pocketbase";

type Props = {
    project: Record,
    action: string
}

export default function NavBar(props: Props){

    return <div class="w-64 flex flex-col gap-6 py-4 shrink-0">
        <div class="px-3">
            <h1 class="text-xl font-semibold">{props.project.name}</h1>
            <p class="text-gray-400">{props.project.slug}</p>
        </div>
        <div class="flex flex-col gap-1">
            <Link href={`/${props.project.slug}/`} class={`flex px-3 py-2 rounded-lg hover:bg-zinc-100 items-center gap-2 ${props.action == props.project.slug ? 'text-black bg-zinc-100' : 'text-gray-400'} hover:text-black`}>
                <IconCode size={24}/>
                <span class="text-lg">Project</span>
            </Link>
            <Link href={`/${props.project.slug}/tasks`} class={`flex px-3 py-2 rounded-lg hover:bg-zinc-100 items-center gap-2 ${props.action == "tasks" ? 'text-black bg-zinc-100' : 'text-gray-400'} hover:text-black`}>
                <IconChecklist size={24}/>
                <span class="text-lg">Tasks</span>
            </Link>
            <Link href={`/${props.project.slug}/releases`} class={`flex px-3 py-2 rounded-lg hover:bg-zinc-100 items-center gap-2 ${props.action == "releases" ? 'text-black bg-zinc-100' : 'text-gray-400'} hover:text-black`}>
                <IconApps size={24}/>
                <span class="text-lg">Releases</span>
            </Link>
            <Link href={`/${props.project.slug}/secrets`} class={`flex px-3 py-2 rounded-lg hover:bg-zinc-100 items-center gap-2 ${props.action == "secrets" ? 'text-black bg-zinc-100' : 'text-gray-400'} hover:text-black`}>
                <IconLock size={24}/>
                <span class="text-lg">Secrets</span>
            </Link>
            <div class="h-[2px] w-10 m-4 bg-gray-300"></div>
            <Link href={`/${props.project.slug}/subscription`} class={`flex px-3 py-2 rounded-lg hover:bg-zinc-100 items-center gap-2 ${props.action == "subscription" ? 'text-black bg-zinc-100' : 'text-gray-400'} hover:text-black`}>
                <IconWallet size={24}/>
                <span class="text-lg">Subscription</span>
            </Link>
        </div>
    </div>
}