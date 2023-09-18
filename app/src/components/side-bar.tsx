import { Link } from 'preact-router/match';
import { Record } from "pocketbase";
import { IconCirclePlus } from '@tabler/icons-preact';
import { client } from '../client';

type Props = {
    projects: Record[],
    slug: string
}

export default function SideBar(props: Props){
    const user = client.authStore.model as Record

    const logout = () => {
        if(window.confirm("Sure you want logout?")){
            client.authStore.clear()
        }
    }
    
    return <div class="w-12 flex flex-col gap-2 py-2 shrink-0 snap-center">
        <Link href={`/`} class="h-12 aspect-square flex items-center justify-center">
            <img src="/logo.svg" class="w-8 h-8"/>
        </Link>
        <div class="flex flex-1 flex-col overflow-y-auto scrollbar gap-2">
            {props.projects.map((project: Record) => <Link href={`/${project.slug}`} class={`h-12 flex items-center justify-center text-lg font-bold ${props.slug == project.slug ? 'bg-black text-white' : 'bg-gray-100 text-gray-400 hover:text-black'} rounded-xl`}>{project.name[0]}</Link>)}
            <Link href="/deploy/blank" class="h-12 aspect-square flex items-center justify-center text-gray-400 hover:text-black">
                <IconCirclePlus size={28}/>
            </Link>
        </div>
        <button onClick={logout} class="h-12 aspect-square flex items-center justify-center">
            <div class="w-9 h-9 flex items-center justify-center rounded-full overflow-hidden bg-red-100 text-red-800">
                {user.avatar != "" ? <img src={client.getFileUrl(user, user.avatar, { thumb: "64x64" })} class="object-cover"/> : <span class="font-bold uppercase">{(user.name != "" ? user.name : user.email)[0]}</span>}
            </div>
        </button>
    </div>
}