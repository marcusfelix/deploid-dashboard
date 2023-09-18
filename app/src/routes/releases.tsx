import { Record } from "pocketbase";
import { IconFileZip, IconPackageExport } from "@tabler/icons-preact";
import { useQuery } from "react-query";
import { client } from "../client";
import ReactMarkdown from "react-markdown";

type Props = {
    project: Record
}

export default function Releases(props: Props){

    const releases = useQuery(`deploy-releases-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/releases`, {
        method: "GET"
    }))

    const download = async (id: string) => {
        const url = await client.send(`/api/projects/${props.project.slug}/releases/${id}`, {
            method: "GET"
        })
        
        // Open download link in new tab
        window.open(url,'_blank');
    }

    // TODO: Test assets file downlooad
    return <div class="w-full md:flex-1 flex flex-col gap-6 shrink-0 snap-center mx-auto max-w-4xl">
        <div class="py-4">
            <h1 class="text-xl font-semibold">Releases</h1>
        </div>
        <div class="flex flex-col gap-4">
            {(releases.data ?? []).map((release: any) => <div class="flex flex-col gap-4 border-2 rounded-xl py-4">
                <div class="flex gap-4 px-4 items-center">
                    <IconPackageExport size={26}/>
                    <h3 class="text-lg flex-1 font-bold">{release.name}</h3>
                    <span class="text-sm text-gray-400">{new Date(release.created).toLocaleDateString()}</span>
                </div>
                <hr class="border-t-2 border-gray-300 border-dashed mx-4"/>
                <div class="prose px-4">
                    <ReactMarkdown>{release.body}</ReactMarkdown>
                </div>
                <div class="flex items-center gap-2 px-4">
                    <div>
                        <img src={release.author.avatar} class="w-6 h-6 rounded-full object-cover"/>
                    </div>
                    <span class="font-bold">{release.author.name}</span>
                </div>
                <div class="flex md:flex-row flex-col gap-4 px-4 mt-4">
                    {release.assets.map((e: any) => <button onClick={() => download(e.id)} target="_blank" class="flex items-center gap-2 px-3 py-2 bg-gray-100 hover:text-white hover:bg-black rounded-lg">
                        <IconFileZip size={20}/> {e.name}
                    </button>)}
                </div>
            </div>)}
            
        </div>
    </div>
}